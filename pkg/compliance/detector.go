// Package compliance provides PII and sensitive data detection.
package compliance

import (
	"regexp"
	"strings"

	"github.com/enterprise/opencode-enterprise-shield/pkg/types"
)

// Detector scans content for compliance violations (PII, secrets, etc.).
type Detector struct {
	patterns      map[string]*Pattern
	blockCritical bool
}

// Pattern represents a compliance detection pattern.
type Pattern struct {
	Name        string
	Type        string
	Regex       *regexp.Regexp
	Severity    types.Severity
	Enabled     bool
	Validator   func(string) bool // Optional validation function
	Description string
}

// NewDetector creates a new compliance detector.
func NewDetector(blockCritical bool) *Detector {
	d := &Detector{
		patterns:      make(map[string]*Pattern),
		blockCritical: blockCritical,
	}
	d.loadDefaultPatterns()
	return d
}

// loadDefaultPatterns loads the default compliance detection patterns.
func (d *Detector) loadDefaultPatterns() {
	// SSN (Social Security Number)
	d.patterns["ssn"] = &Pattern{
		Name:        "Social Security Number",
		Type:        "SSN",
		Regex:       regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`),
		Severity:    types.SeverityCritical,
		Enabled:     true,
		Description: "Detects US Social Security Numbers",
	}

	// Credit Card Numbers (with Luhn validation)
	d.patterns["credit_card"] = &Pattern{
		Name:        "Credit Card Number",
		Type:        "CREDIT_CARD",
		Regex:       regexp.MustCompile(`\b\d{4}[- ]?\d{4}[- ]?\d{4}[- ]?\d{4}\b`),
		Severity:    types.SeverityCritical,
		Enabled:     true,
		Validator:   ValidateLuhn,
		Description: "Detects credit card numbers with Luhn validation",
	}

	// AWS Access Key
	d.patterns["aws_key"] = &Pattern{
		Name:        "AWS Access Key",
		Type:        "API_KEY",
		Regex:       regexp.MustCompile(`\bAKIA[0-9A-Z]{16}\b`),
		Severity:    types.SeverityCritical,
		Enabled:     true,
		Description: "Detects AWS access key IDs",
	}

	// AWS Secret Key
	d.patterns["aws_secret"] = &Pattern{
		Name:        "AWS Secret Key",
		Type:        "API_KEY",
		Regex:       regexp.MustCompile(`(?i)aws.{0,20}secret.{0,20}['\"][0-9a-zA-Z/+=]{40}['\"]`),
		Severity:    types.SeverityCritical,
		Enabled:     true,
		Description: "Detects AWS secret access keys",
	}

	// GitHub Token
	d.patterns["github_token"] = &Pattern{
		Name:        "GitHub Token",
		Type:        "API_KEY",
		Regex:       regexp.MustCompile(`\bghp_[a-zA-Z0-9]{36}\b`),
		Severity:    types.SeverityCritical,
		Enabled:     true,
		Description: "Detects GitHub personal access tokens",
	}

	// GitHub OAuth
	d.patterns["github_oauth"] = &Pattern{
		Name:        "GitHub OAuth Token",
		Type:        "API_KEY",
		Regex:       regexp.MustCompile(`\bgho_[a-zA-Z0-9]{36}\b`),
		Severity:    types.SeverityCritical,
		Enabled:     true,
		Description: "Detects GitHub OAuth tokens",
	}

	// Generic API Key
	d.patterns["generic_api_key"] = &Pattern{
		Name:        "Generic API Key",
		Type:        "API_KEY",
		Regex:       regexp.MustCompile(`(?i)(api[_-]?key|apikey)\s*[:=]\s*['"]?[a-zA-Z0-9]{20,}['"]?`),
		Severity:    types.SeverityHigh,
		Enabled:     true,
		Description: "Detects generic API keys in code",
	}

	// Private Key
	d.patterns["private_key"] = &Pattern{
		Name:        "Private Key",
		Type:        "PRIVATE_KEY",
		Regex:       regexp.MustCompile(`-----BEGIN (RSA |EC |DSA |OPENSSH )?PRIVATE KEY-----`),
		Severity:    types.SeverityCritical,
		Enabled:     true,
		Description: "Detects private key headers",
	}

	// Password in string
	d.patterns["password_string"] = &Pattern{
		Name:        "Password in String",
		Type:        "PASSWORD",
		Regex:       regexp.MustCompile(`(?i)(password|passwd|pwd)\s*[:=]\s*['"]?[^\s'"]{8,}['"]?`),
		Severity:    types.SeverityHigh,
		Enabled:     true,
		Description: "Detects passwords in code strings",
	}

	// Bearer Token
	d.patterns["bearer_token"] = &Pattern{
		Name:        "Bearer Token",
		Type:        "API_KEY",
		Regex:       regexp.MustCompile(`(?i)bearer\s+[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+`),
		Severity:    types.SeverityHigh,
		Enabled:     true,
		Description: "Detects JWT bearer tokens",
	}

	// OpenAI API Key
	d.patterns["openai_key"] = &Pattern{
		Name:        "OpenAI API Key",
		Type:        "API_KEY",
		Regex:       regexp.MustCompile(`\bsk-[a-zA-Z0-9]{48}\b`),
		Severity:    types.SeverityCritical,
		Enabled:     true,
		Description: "Detects OpenAI API keys",
	}

	// Anthropic API Key
	d.patterns["anthropic_key"] = &Pattern{
		Name:        "Anthropic API Key",
		Type:        "API_KEY",
		Regex:       regexp.MustCompile(`\bsk-ant-[a-zA-Z0-9-]{40,}\b`),
		Severity:    types.SeverityCritical,
		Enabled:     true,
		Description: "Detects Anthropic API keys",
	}

	// Slack Token
	d.patterns["slack_token"] = &Pattern{
		Name:        "Slack Token",
		Type:        "API_KEY",
		Regex:       regexp.MustCompile(`xox[baprs]-[0-9]{10,13}-[0-9]{10,13}-[a-zA-Z0-9]{24}`),
		Severity:    types.SeverityCritical,
		Enabled:     true,
		Description: "Detects Slack API tokens",
	}

	// Azure Storage Key
	d.patterns["azure_storage"] = &Pattern{
		Name:        "Azure Storage Key",
		Type:        "API_KEY",
		Regex:       regexp.MustCompile(`(?i)AccountKey=[a-zA-Z0-9+/=]{88}`),
		Severity:    types.SeverityCritical,
		Enabled:     true,
		Description: "Detects Azure storage account keys",
	}
}

// Scan scans content for compliance violations.
func (d *Detector) Scan(content string) types.ComplianceResult {
	result := types.ComplianceResult{
		HasViolations: false,
		ShouldBlock:   false,
		Violations:    make([]types.Violation, 0),
	}

	for _, pattern := range d.patterns {
		if !pattern.Enabled {
			continue
		}

		matches := pattern.Regex.FindAllStringIndex(content, -1)
		for _, match := range matches {
			matchedValue := content[match[0]:match[1]]

			// Run validator if present
			if pattern.Validator != nil && !pattern.Validator(matchedValue) {
				continue // Failed validation, skip
			}

			violation := types.Violation{
				RuleID:        strings.ToLower(pattern.Type),
				RuleName:      pattern.Name,
				Type:          pattern.Type,
				Severity:      pattern.Severity,
				RedactedValue: redactSensitive(matchedValue),
				Position:      match[0],
				Length:        len(matchedValue),
			}

			result.Violations = append(result.Violations, violation)
			result.HasViolations = true

			if d.blockCritical && pattern.Severity == types.SeverityCritical {
				result.ShouldBlock = true
			}
		}
	}

	return result
}

// EnablePattern enables or disables a specific pattern.
func (d *Detector) EnablePattern(name string, enabled bool) {
	if pattern, ok := d.patterns[name]; ok {
		pattern.Enabled = enabled
	}
}

// AddPattern adds a custom pattern to the detector.
func (d *Detector) AddPattern(name string, pattern *Pattern) {
	d.patterns[name] = pattern
}

// redactSensitive creates a redacted version of sensitive data for logging.
func redactSensitive(value string) string {
	if len(value) <= 4 {
		return "****"
	}
	if len(value) <= 8 {
		return value[:2] + "****" + value[len(value)-2:]
	}
	return value[:4] + "****" + value[len(value)-4:]
}

