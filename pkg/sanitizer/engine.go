// Package sanitizer provides the sanitization engine for detecting and masking sensitive data.
package sanitizer

import (
	"fmt"
	"regexp"
	"sort"
	"sync"
	"time"

	"github.com/enterprise/opencode-enterprise-shield/pkg/types"
)

// Engine is the sanitization engine that processes content.
type Engine struct {
	rules         []types.SanitizationRule
	compiledRules map[string]*regexp.Regexp
	mu            sync.RWMutex
	aliasGen      *AliasGenerator
	regexTimeout  time.Duration
}

// NewEngine creates a new sanitization engine with the given rules.
func NewEngine(rules []types.SanitizationRule) *Engine {
	e := &Engine{
		rules:         make([]types.SanitizationRule, 0),
		compiledRules: make(map[string]*regexp.Regexp),
		aliasGen:      NewAliasGenerator(),
		regexTimeout:  50 * time.Millisecond,
	}
	e.LoadRules(rules)
	return e
}

// LoadRules loads and compiles sanitization rules.
func (e *Engine) LoadRules(rules []types.SanitizationRule) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.rules = make([]types.SanitizationRule, 0, len(rules))
	e.compiledRules = make(map[string]*regexp.Regexp)

	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}

		compiled, err := regexp.Compile(rule.Pattern)
		if err != nil {
			return fmt.Errorf("failed to compile pattern for rule %s: %w", rule.Name, err)
		}

		e.rules = append(e.rules, rule)
		e.compiledRules[rule.RuleID] = compiled
	}

	// Sort rules by order
	sort.Slice(e.rules, func(i, j int) bool {
		return e.rules[i].Order < e.rules[j].Order
	})

	return nil
}

// Sanitize processes content and replaces sensitive data with aliases.
func (e *Engine) Sanitize(content string, session *types.Session) types.SanitizationResult {
	startTime := time.Now()

	result := types.SanitizationResult{
		SanitizedContent: content,
		WasSanitized:     false,
		MappingsCreated:  make(map[string]string),
		Violations:       make([]types.Violation, 0),
	}

	e.mu.RLock()
	rules := e.rules
	compiledRules := e.compiledRules
	e.mu.RUnlock()

	workingContent := content

	for _, rule := range rules {
		compiled, ok := compiledRules[rule.RuleID]
		if !ok {
			continue
		}

		// Find all matches
		matches := compiled.FindAllStringIndex(workingContent, -1)
		if len(matches) == 0 {
			continue
		}

		// Process matches in reverse order to preserve positions
		for i := len(matches) - 1; i >= 0; i-- {
			match := matches[i]
			matchedValue := workingContent[match[0]:match[1]]

			// Check exceptions
			if e.isException(matchedValue, rule.Exceptions) {
				continue
			}

			// Check if critical severity should block
			if rule.Severity == types.SeverityCritical {
				result.ShouldBlock = true
				result.BlockReason = fmt.Sprintf("Critical violation detected: %s", rule.Name)
			}

			// Get or create alias
			alias, isNew := e.getOrCreateAlias(session, matchedValue, rule.Prefix)

			// Replace in content
			workingContent = workingContent[:match[0]] + alias + workingContent[match[1]:]

			// Record violation
			violation := types.Violation{
				RuleID:        rule.RuleID,
				RuleName:      rule.Name,
				Type:          rule.Prefix,
				Severity:      rule.Severity,
				RedactedValue: redactValue(matchedValue),
				Position:      match[0],
				Length:        len(matchedValue),
			}
			result.Violations = append(result.Violations, violation)

			// Track new mappings
			if isNew {
				result.MappingsCreated[matchedValue] = alias
			}

			result.WasSanitized = true
		}
	}

	result.SanitizedContent = workingContent
	result.ProcessingTimeMs = time.Since(startTime).Milliseconds()

	return result
}

// getOrCreateAlias retrieves existing alias or creates a new one.
func (e *Engine) getOrCreateAlias(session *types.Session, original, prefix string) (string, bool) {
	// Check if alias already exists
	if alias, ok := session.GetAlias(original); ok {
		return alias, false
	}

	// Generate new alias
	alias := e.aliasGen.Generate(session, prefix)
	session.AddMapping(original, alias)

	return alias, true
}

// isException checks if a value matches any exception pattern.
func (e *Engine) isException(value string, exceptions []string) bool {
	for _, exception := range exceptions {
		matched, err := regexp.MatchString(exception, value)
		if err == nil && matched {
			return true
		}
	}
	return false
}

// redactValue creates a redacted version of a value for logging.
func redactValue(value string) string {
	if len(value) <= 3 {
		return "***"
	}
	if len(value) <= 6 {
		return value[:1] + "***" + value[len(value)-1:]
	}
	return value[:3] + "***" + value[len(value)-3:]
}

// GetRules returns a copy of the loaded rules.
func (e *Engine) GetRules() []types.SanitizationRule {
	e.mu.RLock()
	defer e.mu.RUnlock()

	rules := make([]types.SanitizationRule, len(e.rules))
	copy(rules, e.rules)
	return rules
}

