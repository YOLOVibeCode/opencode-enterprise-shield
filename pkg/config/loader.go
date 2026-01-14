// Package config provides configuration loading for Enterprise Shield.
package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/enterprise/opencode-enterprise-shield/pkg/hooks"
	"github.com/enterprise/opencode-enterprise-shield/pkg/types"
	"gopkg.in/yaml.v3"
)

// FullConfig represents the complete configuration file structure.
type FullConfig struct {
	Enabled    bool            `yaml:"enabled"`
	Session    SessionConfig   `yaml:"session"`
	Rules      []types.SanitizationRule `yaml:"rules"`
	Compliance ComplianceConfig `yaml:"compliance"`
	Policy     PolicyConfig    `yaml:"policy"`
	Audit      AuditConfig     `yaml:"audit"`
}

// SessionConfig holds session-related configuration.
type SessionConfig struct {
	TTL         string `yaml:"ttl"`
	MaxMappings int    `yaml:"maxMappings"`
	Encryption  bool   `yaml:"encryption"`
}

// ComplianceConfig holds compliance detection configuration.
type ComplianceConfig struct {
	BlockOnCritical bool             `yaml:"blockOnCritical"`
	Detectors       []DetectorConfig `yaml:"detectors"`
}

// DetectorConfig holds individual detector configuration.
type DetectorConfig struct {
	Type        string `yaml:"type"`
	Enabled     bool   `yaml:"enabled"`
	Severity    string `yaml:"severity"`
	ValidateLuhn bool  `yaml:"validateLuhn,omitempty"`
}

// PolicyConfig holds policy configuration.
type PolicyConfig struct {
	DefaultAccessLevel string `yaml:"defaultAccessLevel"`
	RequireAuth        bool   `yaml:"requireAuth"`
}

// AuditConfig holds audit logging configuration.
type AuditConfig struct {
	Enabled       bool   `yaml:"enabled"`
	LogPath       string `yaml:"logPath"`
	SignEntries   bool   `yaml:"signEntries"`
	RetentionDays int    `yaml:"retentionDays"`
}

// Load loads configuration from a YAML file.
func Load(path string) (*FullConfig, error) {
	// Expand home directory
	if path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(home, path[1:])
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config FullConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// LoadOrDefault loads configuration or returns defaults.
func LoadOrDefault(path string) *FullConfig {
	config, err := Load(path)
	if err != nil {
		return DefaultFullConfig()
	}
	return config
}

// DefaultFullConfig returns the default full configuration.
func DefaultFullConfig() *FullConfig {
	return &FullConfig{
		Enabled: true,
		Session: SessionConfig{
			TTL:         "8h",
			MaxMappings: 10000,
			Encryption:  true,
		},
		Compliance: ComplianceConfig{
			BlockOnCritical: true,
			Detectors: []DetectorConfig{
				{Type: "ssn", Enabled: true, Severity: "critical"},
				{Type: "credit_card", Enabled: true, Severity: "critical", ValidateLuhn: true},
				{Type: "api_key", Enabled: true, Severity: "critical"},
			},
		},
		Policy: PolicyConfig{
			DefaultAccessLevel: "sanitized_only",
			RequireAuth:        true,
		},
		Audit: AuditConfig{
			Enabled:       true,
			LogPath:       "~/.opencode/logs/enterprise-shield",
			SignEntries:   true,
			RetentionDays: 365,
		},
	}
}

// ToHooksConfig converts the full config to a hooks.Config.
func (c *FullConfig) ToHooksConfig() *hooks.Config {
	ttl, _ := time.ParseDuration(c.Session.TTL)
	if ttl == 0 {
		ttl = 8 * time.Hour
	}

	return &hooks.Config{
		Enabled:         c.Enabled,
		SessionTTL:      ttl,
		MaxMappings:     c.Session.MaxMappings,
		BlockOnCritical: c.Compliance.BlockOnCritical,
		AuditLogPath:    c.Audit.LogPath,
		SignAuditLogs:   c.Audit.SignEntries,
		RetentionDays:   c.Audit.RetentionDays,
	}
}

// Save saves configuration to a YAML file.
func Save(config *FullConfig, path string) error {
	// Expand home directory
	if path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		path = filepath.Join(home, path[1:])
	}

	// Create directory if needed
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0640)
}

