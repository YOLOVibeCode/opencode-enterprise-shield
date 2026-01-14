// Package hooks provides OpenCode integration middleware for Enterprise Shield.
package hooks

import (
	"fmt"
	"time"

	"github.com/enterprise/opencode-enterprise-shield/pkg/audit"
	"github.com/enterprise/opencode-enterprise-shield/pkg/compliance"
	"github.com/enterprise/opencode-enterprise-shield/pkg/desanitizer"
	"github.com/enterprise/opencode-enterprise-shield/pkg/policy"
	"github.com/enterprise/opencode-enterprise-shield/pkg/sanitizer"
	"github.com/enterprise/opencode-enterprise-shield/pkg/session"
	"github.com/enterprise/opencode-enterprise-shield/pkg/types"
)

// Shield is the main Enterprise Shield middleware.
type Shield struct {
	sanitizer      *sanitizer.Engine
	desanitizer    *desanitizer.Engine
	compliance     *compliance.Detector
	sessionManager *session.Manager
	policyEngine   *policy.Engine
	auditLogger    *audit.Logger
	config         *Config
}

// Config holds the Shield configuration.
type Config struct {
	Enabled         bool          `yaml:"enabled"`
	SessionTTL      time.Duration `yaml:"sessionTTL"`
	MaxMappings     int           `yaml:"maxMappings"`
	BlockOnCritical bool          `yaml:"blockOnCritical"`
	AuditLogPath    string        `yaml:"auditLogPath"`
	SignAuditLogs   bool          `yaml:"signAuditLogs"`
	RetentionDays   int           `yaml:"retentionDays"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		Enabled:         true,
		SessionTTL:      8 * time.Hour,
		MaxMappings:     10000,
		BlockOnCritical: true,
		AuditLogPath:    "~/.opencode/logs/enterprise-shield",
		SignAuditLogs:   true,
		RetentionDays:   365,
	}
}

// NewShield creates a new Shield middleware instance.
func NewShield(config *Config) (*Shield, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Initialize components
	sanitizerEngine := sanitizer.NewEngine(sanitizer.DefaultRules())
	desanitizerEngine := desanitizer.NewEngine()
	complianceDetector := compliance.NewDetector(config.BlockOnCritical)
	sessionManager := session.NewManager(config.SessionTTL, config.MaxMappings)
	policyEngine := policy.NewEngine()

	// Initialize audit logger
	auditLogger, err := audit.NewLogger(config.AuditLogPath, config.SignAuditLogs, config.RetentionDays)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize audit logger: %w", err)
	}

	return &Shield{
		sanitizer:      sanitizerEngine,
		desanitizer:    desanitizerEngine,
		compliance:     complianceDetector,
		sessionManager: sessionManager,
		policyEngine:   policyEngine,
		auditLogger:    auditLogger,
		config:         config,
	}, nil
}

// ProcessRequest processes an outgoing request (before sending to LLM).
func (s *Shield) ProcessRequest(req types.Request) types.Response {
	startTime := time.Now()

	response := types.Response{
		SessionID: req.SessionID,
	}

	// Check if shield is enabled
	if !s.config.Enabled {
		response.Content = req.Content
		return response
	}

	// Step 1: Policy check
	policyCtx := policy.PolicyContext{
		UserID:     req.UserID,
		Department: req.Department,
		Provider:   req.Provider,
		Content:    req.Content,
	}
	policyDecision := s.policyEngine.Evaluate(policyCtx)

	if policyDecision.Action == types.ActionBlock {
		response.Blocked = true
		response.BlockReason = policyDecision.Reason
		s.logRequest(req, response, policyDecision.Action, nil, time.Since(startTime).Milliseconds())
		return response
	}

	// Step 2: Compliance scan
	complianceResult := s.compliance.Scan(req.Content)
	if complianceResult.ShouldBlock {
		response.Blocked = true
		response.BlockReason = "Critical compliance violation detected"
		response.Violations = complianceResult.Violations
		s.logRequest(req, response, types.ActionBlock, complianceResult.Violations, time.Since(startTime).Milliseconds())
		return response
	}

	// Step 3: Get or create session
	sess, _ := s.sessionManager.GetOrCreate(req.UserID, req.Department, req.SessionID)
	response.SessionID = sess.SessionID

	// Step 4: Sanitization (if required)
	if policyDecision.Action == types.ActionAllowWithSanitization {
		sanitizeResult := s.sanitizer.Sanitize(req.Content, sess)
		
		if sanitizeResult.ShouldBlock {
			response.Blocked = true
			response.BlockReason = sanitizeResult.BlockReason
			response.Violations = sanitizeResult.Violations
			s.logRequest(req, response, types.ActionBlock, sanitizeResult.Violations, time.Since(startTime).Milliseconds())
			return response
		}

		response.Content = sanitizeResult.SanitizedContent
		response.WasSanitized = sanitizeResult.WasSanitized
		response.MappingsCreated = sanitizeResult.MappingsCreated
		response.Violations = append(response.Violations, sanitizeResult.Violations...)
	} else {
		response.Content = req.Content
	}

	// Log the request
	allViolations := append(complianceResult.Violations, response.Violations...)
	s.logRequest(req, response, policyDecision.Action, allViolations, time.Since(startTime).Milliseconds())

	return response
}

// ProcessResponse processes an incoming response (from LLM back to user).
func (s *Shield) ProcessResponse(content string, sessionID string) types.DesanitizationResult {
	// Get session
	sess, ok := s.sessionManager.Get(sessionID)
	if !ok {
		return types.DesanitizationResult{
			DesanitizedContent: content,
			ReplacementsCount:  0,
		}
	}

	// Desanitize the response
	return s.desanitizer.Desanitize(content, sess)
}

// ScanContent performs a compliance scan without processing.
func (s *Shield) ScanContent(content string) types.ComplianceResult {
	return s.compliance.Scan(content)
}

// GetSession retrieves session information.
func (s *Shield) GetSession(sessionID string) (*types.Session, bool) {
	return s.sessionManager.Get(sessionID)
}

// ClearSession clears a user's session.
func (s *Shield) ClearSession(userID string) {
	s.sessionManager.Clear(userID)
}

// SetUserPolicy sets a user's policy.
func (s *Shield) SetUserPolicy(userID string, policy *types.UserPolicy) {
	s.policyEngine.SetUserPolicy(userID, policy)
}

// GetStats returns shield statistics.
func (s *Shield) GetStats() ShieldStats {
	sessionStats := s.sessionManager.GetStats()
	return ShieldStats{
		SessionStats: sessionStats,
		RulesLoaded:  len(s.sanitizer.GetRules()),
	}
}

// Close cleans up resources.
func (s *Shield) Close() error {
	return s.auditLogger.Close()
}

// logRequest logs a request to the audit log.
func (s *Shield) logRequest(req types.Request, resp types.Response, action types.Action, violations []types.Violation, processingMs int64) {
	entry := s.auditLogger.CreateEntry(
		req.UserID,
		resp.SessionID,
		req.Department,
		req.Provider,
		resp.WasSanitized,
		violations,
		action,
		processingMs,
	)
	s.auditLogger.Log(entry)
}

// ShieldStats contains statistics about the shield.
type ShieldStats struct {
	SessionStats session.SessionStats `json:"sessionStats"`
	RulesLoaded  int                  `json:"rulesLoaded"`
}

// --- OpenCode Hook Interface ---

// Hook represents an OpenCode plugin hook.
type Hook struct {
	shield *Shield
}

// NewHook creates a new OpenCode hook.
func NewHook(config *Config) (*Hook, error) {
	shield, err := NewShield(config)
	if err != nil {
		return nil, err
	}
	return &Hook{shield: shield}, nil
}

// OnRequest is called before a request is sent to the LLM.
func (h *Hook) OnRequest(userID, content, provider string) types.Response {
	req := types.Request{
		UserID:   userID,
		Content:  content,
		Provider: provider,
	}
	return h.shield.ProcessRequest(req)
}

// OnResponse is called when a response is received from the LLM.
func (h *Hook) OnResponse(content, sessionID string) types.DesanitizationResult {
	return h.shield.ProcessResponse(content, sessionID)
}

// OnScan performs a compliance scan.
func (h *Hook) OnScan(content string) types.ComplianceResult {
	return h.shield.ScanContent(content)
}

// Close cleans up hook resources.
func (h *Hook) Close() error {
	return h.shield.Close()
}

