// Package main provides the Enterprise Shield plugin entry point for OpenCode.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/enterprise/opencode-enterprise-shield/pkg/config"
	"github.com/enterprise/opencode-enterprise-shield/pkg/hooks"
	"github.com/enterprise/opencode-enterprise-shield/pkg/types"
)

var (
	// Set via ldflags at build time
	version     = "dev"
	buildNumber = "0"
	buildTime   = "unknown"
)

const (
	configPath = "~/.opencode/config/enterprise-shield.yaml"
)

// Plugin represents the Enterprise Shield plugin.
type Plugin struct {
	hook   *hooks.Hook
	config *config.FullConfig
}

// NewPlugin creates a new plugin instance.
func NewPlugin() (*Plugin, error) {
	cfg := config.LoadOrDefault(configPath)
	
	hook, err := hooks.NewHook(cfg.ToHooksConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize plugin: %w", err)
	}

	return &Plugin{
		hook:   hook,
		config: cfg,
	}, nil
}

// ProcessRequest handles outgoing requests to LLM.
func (p *Plugin) ProcessRequest(userID, content, provider string) types.Response {
	return p.hook.OnRequest(userID, content, provider)
}

// ProcessResponse handles incoming responses from LLM.
func (p *Plugin) ProcessResponse(content, sessionID string) types.DesanitizationResult {
	return p.hook.OnResponse(content, sessionID)
}

// ScanContent performs a compliance scan.
func (p *Plugin) ScanContent(content string) types.ComplianceResult {
	return p.hook.OnScan(content)
}

// Close cleans up plugin resources.
func (p *Plugin) Close() error {
	return p.hook.Close()
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "version":
		fmt.Printf("Enterprise Shield Plugin v%s\n", version)
		if buildNumber != "0" {
			fmt.Printf("Build: #%s\n", buildNumber)
		}
		if buildTime != "unknown" {
			fmt.Printf("Built: %s\n", buildTime)
		}

	case "init":
		// Initialize configuration
		cfg := config.DefaultFullConfig()
		if err := config.Save(cfg, configPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Configuration initialized at", configPath)

	case "scan":
		// Scan content from stdin or argument
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: enterprise-shield scan <content>")
			os.Exit(1)
		}
		content := os.Args[2]
		
		plugin, err := NewPlugin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		defer plugin.Close()

		result := plugin.ScanContent(content)
		printJSON(result)

	case "process":
		// Process a request (for testing)
		if len(os.Args) < 5 {
			fmt.Fprintln(os.Stderr, "Usage: enterprise-shield process <userID> <content> <provider>")
			os.Exit(1)
		}
		userID := os.Args[2]
		content := os.Args[3]
		provider := os.Args[4]

		plugin, err := NewPlugin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		defer plugin.Close()

		result := plugin.ProcessRequest(userID, content, provider)
		printJSON(result)

	case "serve":
		// Run as a service (for OpenCode integration)
		fmt.Println("Enterprise Shield Plugin v" + version)
		fmt.Println("Running in server mode...")
		
		plugin, err := NewPlugin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		defer plugin.Close()

		// In server mode, we would listen for JSON-RPC or stdio commands
		// This is a placeholder for the actual OpenCode integration
		fmt.Println("Ready to process requests.")
		fmt.Println("Press Ctrl+C to exit.")
		
		// Block forever (in real implementation, this would be the event loop)
		select {}

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Print(`Enterprise Shield Plugin for OpenCode

Usage:
  enterprise-shield <command> [arguments]

Commands:
  version              Show version information
  init                 Initialize default configuration
  scan <content>       Scan content for compliance violations
  process <user> <content> <provider>
                       Process a request (sanitize and check policy)
  serve                Run in server mode for OpenCode integration

Examples:
  enterprise-shield version
  enterprise-shield init
  enterprise-shield scan "My SSN is 123-45-6789"
  enterprise-shield process user@example.com "Query ServerDB01" openai

Configuration:
  Default config path: ~/.opencode/config/enterprise-shield.yaml
`)
}

func printJSON(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

