// Package types contains shared type definitions for Enterprise Shield.
package types

import "time"

// Severity represents the severity level of a violation or rule.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
)

// Action represents the action to take for a request.
type Action string

const (
	ActionAllow                 Action = "allow"
	ActionAllowWithSanitization Action = "allow_with_sanitization"
	ActionAllowWithWarning      Action = "allow_with_warning"
	ActionBlock                 Action = "block"
	ActionRateLimited           Action = "rate_limited"
)

// AccessLevel represents user access level for RBAC.
type AccessLevel string

const (
	AccessUnrestricted  AccessLevel = "unrestricted"
	AccessSanitizedOnly AccessLevel = "sanitized_only"
	AccessBlocked       AccessLevel = "blocked"
)

// SessionStatus represents the status of a session.
type SessionStatus string

const (
	SessionActive     SessionStatus = "active"
	SessionExpired    SessionStatus = "expired"
	SessionTerminated SessionStatus = "terminated"
)

// Session represents a user session with mappings.
type Session struct {
	SessionID       string            `json:"sessionId"`
	UserID          string            `json:"userId"`
	Department      string            `json:"department,omitempty"`
	CreatedAt       time.Time         `json:"createdAt"`
	ExpiresAt       time.Time         `json:"expiresAt"`
	LastAccessedAt  time.Time         `json:"lastAccessedAt,omitempty"`
	Status          SessionStatus     `json:"status"`
	Mappings        map[string]string `json:"mappings"`        // Original → Alias
	ReverseMappings map[string]string `json:"reverseMappings"` // Alias → Original
	RequestCount    int               `json:"requestCount"`
	Counters        map[string]int    `json:"counters"` // Per-prefix counters
}

// NewSession creates a new session for a user.
func NewSession(sessionID, userID, department string, ttl time.Duration) *Session {
	now := time.Now()
	return &Session{
		SessionID:       sessionID,
		UserID:          userID,
		Department:      department,
		CreatedAt:       now,
		ExpiresAt:       now.Add(ttl),
		LastAccessedAt:  now,
		Status:          SessionActive,
		Mappings:        make(map[string]string),
		ReverseMappings: make(map[string]string),
		RequestCount:    0,
		Counters:        make(map[string]int),
	}
}

// IsExpired checks if the session has expired.
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt) || s.Status != SessionActive
}

// AddMapping adds a mapping and its reverse to the session.
func (s *Session) AddMapping(original, alias string) {
	s.Mappings[original] = alias
	s.ReverseMappings[alias] = original
}

// GetAlias returns the alias for an original value.
func (s *Session) GetAlias(original string) (string, bool) {
	alias, ok := s.Mappings[original]
	return alias, ok
}

// GetOriginal returns the original value for an alias.
func (s *Session) GetOriginal(alias string) (string, bool) {
	original, ok := s.ReverseMappings[alias]
	return original, ok
}

// GetNextCounter returns the next counter value for a prefix.
func (s *Session) GetNextCounter(prefix string) int {
	count := s.Counters[prefix]
	s.Counters[prefix] = count + 1
	return count
}

// Touch updates the last accessed time.
func (s *Session) Touch() {
	s.LastAccessedAt = time.Now()
	s.RequestCount++
}

// SanitizationRule defines a pattern for detecting and masking sensitive data.
type SanitizationRule struct {
	RuleID      string   `json:"ruleId" yaml:"ruleId"`
	Name        string   `json:"name" yaml:"name"`
	Description string   `json:"description,omitempty" yaml:"description"`
	Pattern     string   `json:"pattern" yaml:"pattern"`
	Prefix      string   `json:"prefix" yaml:"prefix"`
	Severity    Severity `json:"severity" yaml:"severity"`
	Enabled     bool     `json:"enabled" yaml:"enabled"`
	Exceptions  []string `json:"exceptions,omitempty" yaml:"exceptions"`
	Order       int      `json:"order,omitempty" yaml:"order"`
}

// UserPolicy defines access rules for a user.
type UserPolicy struct {
	PolicyID          string      `json:"policyId"`
	UserID            string      `json:"userId"`
	Department        string      `json:"department,omitempty"`
	AccessLevel       AccessLevel `json:"accessLevel"`
	AllowedProviders  []string    `json:"allowedProviders,omitempty"`
	DailyRequestLimit int         `json:"dailyRequestLimit"`
	HourlyRequestLimit int        `json:"hourlyRequestLimit"`
	RequiredRules     []string    `json:"requiredRules,omitempty"`
	Enabled           bool        `json:"enabled"`
}

// Violation represents a detected violation.
type Violation struct {
	RuleID         string   `json:"ruleId"`
	RuleName       string   `json:"ruleName"`
	Type           string   `json:"type"`
	Severity       Severity `json:"severity"`
	MatchedPattern string   `json:"matchedPattern,omitempty"`
	RedactedValue  string   `json:"redactedValue"`
	Position       int      `json:"position"`
	Length         int      `json:"length"`
}

// SanitizationResult contains the result of sanitization.
type SanitizationResult struct {
	SanitizedContent string            `json:"sanitizedContent"`
	WasSanitized     bool              `json:"wasSanitized"`
	MappingsCreated  map[string]string `json:"mappingsCreated"`
	Violations       []Violation       `json:"violations"`
	ProcessingTimeMs int64             `json:"processingTimeMs"`
	ShouldBlock      bool              `json:"shouldBlock"`
	BlockReason      string            `json:"blockReason,omitempty"`
}

// DesanitizationResult contains the result of desanitization.
type DesanitizationResult struct {
	DesanitizedContent string   `json:"desanitizedContent"`
	ReplacementsCount  int      `json:"replacementsCount"`
	UnmatchedAliases   []string `json:"unmatchedAliases,omitempty"`
	ProcessingTimeMs   int64    `json:"processingTimeMs"`
}

// ComplianceResult contains the result of compliance scanning.
type ComplianceResult struct {
	HasViolations bool        `json:"hasViolations"`
	ShouldBlock   bool        `json:"shouldBlock"`
	Violations    []Violation `json:"violations"`
}

// PolicyDecision contains the policy evaluation result.
type PolicyDecision struct {
	Action             Action   `json:"action"`
	Reason             string   `json:"reason,omitempty"`
	PolicyApplied      string   `json:"policyApplied,omitempty"`
	RequiredSanitization []string `json:"requiredSanitization,omitempty"`
}

// AuditEntry represents an audit log entry.
type AuditEntry struct {
	EntryID           string      `json:"entryId"`
	Timestamp         time.Time   `json:"timestamp"`
	UserID            string      `json:"userId"`
	SessionID         string      `json:"sessionId,omitempty"`
	Department        string      `json:"department,omitempty"`
	Provider          string      `json:"provider,omitempty"`
	RequestHash       string      `json:"requestHash"`
	ResponseHash      string      `json:"responseHash,omitempty"`
	WasSanitized      bool        `json:"wasSanitized"`
	Violations        []Violation `json:"violations,omitempty"`
	Action            Action      `json:"action"`
	ProcessingTimeMs  int64       `json:"processingTimeMs"`
	Signature         string      `json:"signature,omitempty"`
	PreviousEntryHash string      `json:"previousEntryHash,omitempty"`
}

// Request represents an incoming request to process.
type Request struct {
	UserID     string            `json:"userId"`
	SessionID  string            `json:"sessionId,omitempty"`
	Department string            `json:"department,omitempty"`
	Provider   string            `json:"provider,omitempty"`
	Content    string            `json:"content"`
	Headers    map[string]string `json:"headers,omitempty"`
}

// Response represents the processed response.
type Response struct {
	Content         string            `json:"content"`
	SessionID       string            `json:"sessionId"`
	WasSanitized    bool              `json:"wasSanitized"`
	WasDesanitized  bool              `json:"wasDesanitized"`
	MappingsCreated map[string]string `json:"mappingsCreated,omitempty"`
	Blocked         bool              `json:"blocked"`
	BlockReason     string            `json:"blockReason,omitempty"`
	Violations      []Violation       `json:"violations,omitempty"`
}

