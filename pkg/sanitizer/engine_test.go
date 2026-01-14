package sanitizer

import (
	"testing"
	"time"

	"github.com/enterprise/opencode-enterprise-shield/pkg/types"
)

func TestSanitize_ServerNames(t *testing.T) {
	engine := NewEngine(DefaultRules())
	session := types.NewSession("test-session", "user@test.com", "engineering", 8*time.Hour)

	content := "Query ServerDB01 and ProductionDB for user data"
	result := engine.Sanitize(content, session)

	if !result.WasSanitized {
		t.Error("Expected content to be sanitized")
	}

	// Check that server names were replaced
	if result.SanitizedContent == content {
		t.Error("Expected sanitized content to be different")
	}

	// Verify mappings were created
	if len(result.MappingsCreated) == 0 {
		t.Error("Expected mappings to be created")
	}

	// Verify aliases follow pattern
	for original, alias := range result.MappingsCreated {
		t.Logf("Mapping: %s -> %s", original, alias)
		if alias[:7] != "SERVER_" {
			t.Errorf("Expected SERVER_ prefix, got %s", alias)
		}
	}
}

func TestSanitize_IPAddresses(t *testing.T) {
	engine := NewEngine(DefaultRules())
	session := types.NewSession("test-session", "user@test.com", "engineering", 8*time.Hour)

	content := "Connect to 192.168.1.100 or 10.0.0.50"
	result := engine.Sanitize(content, session)

	if !result.WasSanitized {
		t.Error("Expected content to be sanitized")
	}

	// Check that IPs were replaced
	if result.SanitizedContent == content {
		t.Error("Expected sanitized content to be different")
	}

	// Verify mappings
	for original, alias := range result.MappingsCreated {
		t.Logf("Mapping: %s -> %s", original, alias)
		if alias[:3] != "IP_" {
			t.Errorf("Expected IP_ prefix, got %s", alias)
		}
	}
}

func TestSanitize_MixedContent(t *testing.T) {
	engine := NewEngine(DefaultRules())
	session := types.NewSession("test-session", "user@test.com", "engineering", 8*time.Hour)

	content := "SELECT * FROM ServerDB01.users_prod WHERE ip='192.168.1.100'"
	result := engine.Sanitize(content, session)

	if !result.WasSanitized {
		t.Error("Expected content to be sanitized")
	}

	t.Logf("Original: %s", content)
	t.Logf("Sanitized: %s", result.SanitizedContent)
	t.Logf("Mappings: %v", result.MappingsCreated)
}

func TestSanitize_SessionPersistence(t *testing.T) {
	engine := NewEngine(DefaultRules())
	session := types.NewSession("test-session", "user@test.com", "engineering", 8*time.Hour)

	// First request
	content1 := "Query ServerDB01"
	result1 := engine.Sanitize(content1, session)

	// Second request with same server name
	content2 := "Also check ServerDB01 and ServerDB02"
	result2 := engine.Sanitize(content2, session)

	// ServerDB01 should have the same alias in both
	alias1 := result1.MappingsCreated["ServerDB01"]
	
	// In result2, ServerDB01 should NOT be in MappingsCreated (it already exists)
	if _, ok := result2.MappingsCreated["ServerDB01"]; ok {
		t.Error("ServerDB01 should not be in new mappings (already existed)")
	}

	// But it should still be replaced with the same alias
	if session.Mappings["ServerDB01"] != alias1 {
		t.Error("ServerDB01 alias should persist across requests")
	}

	// ServerDB02 should be new
	if _, ok := result2.MappingsCreated["ServerDB02"]; !ok {
		t.Error("ServerDB02 should be in new mappings")
	}
}

func TestSanitize_NoSensitiveData(t *testing.T) {
	engine := NewEngine(DefaultRules())
	session := types.NewSession("test-session", "user@test.com", "engineering", 8*time.Hour)

	content := "How do I optimize a SQL query?"
	result := engine.Sanitize(content, session)

	if result.WasSanitized {
		t.Error("Expected no sanitization for generic content")
	}

	if result.SanitizedContent != content {
		t.Error("Expected content to remain unchanged")
	}
}

func TestAliasGenerator(t *testing.T) {
	gen := NewAliasGenerator()
	session := types.NewSession("test-session", "user@test.com", "engineering", 8*time.Hour)

	// Generate several aliases
	alias1 := gen.Generate(session, "SERVER")
	alias2 := gen.Generate(session, "SERVER")
	alias3 := gen.Generate(session, "IP")

	if alias1 != "SERVER_0" {
		t.Errorf("Expected SERVER_0, got %s", alias1)
	}
	if alias2 != "SERVER_1" {
		t.Errorf("Expected SERVER_1, got %s", alias2)
	}
	if alias3 != "IP_0" {
		t.Errorf("Expected IP_0, got %s", alias3)
	}
}

