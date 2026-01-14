// Package sanitizer provides default sanitization rules.
package sanitizer

import "github.com/enterprise/opencode-enterprise-shield/pkg/types"

// DefaultRules returns the default set of sanitization rules.
func DefaultRules() []types.SanitizationRule {
	return []types.SanitizationRule{
		// Server/Database names
		{
			RuleID:      "server_names",
			Name:        "Server Names",
			Description: "Detects server and database names like ServerDB01, ProductionDB",
			Pattern:     `\b[A-Z][a-zA-Z]*DB\d*\b`,
			Prefix:      "SERVER",
			Severity:    types.SeverityMedium,
			Enabled:     true,
			Order:       10,
		},
		{
			RuleID:      "prod_databases",
			Name:        "Production Databases",
			Description: "Detects production database names",
			Pattern:     `\b[a-zA-Z]+[_-]?[Pp]rod(uction)?\b`,
			Prefix:      "SERVER",
			Severity:    types.SeverityMedium,
			Enabled:     true,
			Order:       11,
		},

		// Table names
		{
			RuleID:      "table_names_prod",
			Name:        "Production Table Names",
			Description: "Detects table names with prod suffix",
			Pattern:     `\b[a-z_]+_prod\b`,
			Prefix:      "TABLE",
			Severity:    types.SeverityMedium,
			Enabled:     true,
			Order:       20,
		},
		{
			RuleID:      "table_names_users",
			Name:        "User Table Names",
			Description: "Detects user-related table names",
			Pattern:     `\b(users?|accounts?|customers?|employees?)(_\w+)?\b`,
			Prefix:      "TABLE",
			Severity:    types.SeverityMedium,
			Enabled:     true,
			Exceptions:  []string{`^user$`, `^account$`}, // Don't match generic words
			Order:       21,
		},

		// IP Addresses
		{
			RuleID:      "private_ip_10",
			Name:        "Private IP (10.x.x.x)",
			Description: "Detects RFC 1918 private IPs in 10.0.0.0/8 range",
			Pattern:     `\b10\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`,
			Prefix:      "IP",
			Severity:    types.SeverityHigh,
			Enabled:     true,
			Order:       30,
		},
		{
			RuleID:      "private_ip_172",
			Name:        "Private IP (172.16-31.x.x)",
			Description: "Detects RFC 1918 private IPs in 172.16.0.0/12 range",
			Pattern:     `\b172\.(1[6-9]|2\d|3[01])\.\d{1,3}\.\d{1,3}\b`,
			Prefix:      "IP",
			Severity:    types.SeverityHigh,
			Enabled:     true,
			Order:       31,
		},
		{
			RuleID:      "private_ip_192",
			Name:        "Private IP (192.168.x.x)",
			Description: "Detects RFC 1918 private IPs in 192.168.0.0/16 range",
			Pattern:     `\b192\.168\.\d{1,3}\.\d{1,3}\b`,
			Prefix:      "IP",
			Severity:    types.SeverityHigh,
			Enabled:     true,
			Order:       32,
		},

		// Connection strings
		{
			RuleID:      "connection_string",
			Name:        "Connection Strings",
			Description: "Detects database connection strings",
			Pattern:     `(?i)(server|data source|host)=[^;]+;`,
			Prefix:      "CONNSTR",
			Severity:    types.SeverityHigh,
			Enabled:     true,
			Order:       40,
		},

		// File paths
		{
			RuleID:      "windows_path",
			Name:        "Windows File Paths",
			Description: "Detects Windows file paths",
			Pattern:     `[A-Za-z]:\\[^\s*?"<>|:]+`,
			Prefix:      "PATH",
			Severity:    types.SeverityMedium,
			Enabled:     true,
			Order:       50,
		},
		{
			RuleID:      "unc_path",
			Name:        "UNC Paths",
			Description: "Detects UNC network paths",
			Pattern:     `\\\\[a-zA-Z0-9._-]+\\[^\s]+`,
			Prefix:      "PATH",
			Severity:    types.SeverityMedium,
			Enabled:     true,
			Order:       51,
		},

		// Hostnames
		{
			RuleID:      "internal_hostname",
			Name:        "Internal Hostnames",
			Description: "Detects internal hostnames with common suffixes",
			Pattern:     `\b[a-z][a-z0-9-]*\.(internal|local|corp|lan)\b`,
			Prefix:      "HOST",
			Severity:    types.SeverityMedium,
			Enabled:     true,
			Order:       60,
		},

		// Email addresses (internal)
		{
			RuleID:      "internal_email",
			Name:        "Internal Email Addresses",
			Description: "Detects email addresses (for internal domain detection)",
			Pattern:     `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`,
			Prefix:      "EMAIL",
			Severity:    types.SeverityLow,
			Enabled:     false, // Disabled by default, enable for stricter environments
			Order:       70,
		},
	}
}

