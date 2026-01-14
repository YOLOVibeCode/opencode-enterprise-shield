package desanitizer

import (
	"testing"
	"time"

	"github.com/enterprise/opencode-enterprise-shield/pkg/types"
)

func TestDesanitize_Basic(t *testing.T) {
	engine := NewEngine()
	session := types.NewSession("test-session", "user@test.com", "engineering", 8*time.Hour)

	// Add mappings
	session.AddMapping("ServerDB01", "SERVER_0")
	session.AddMapping("users_prod", "TABLE_0")
	session.AddMapping("192.168.1.100", "IP_0")

	content := "To optimize SERVER_0.TABLE_0, add an index. The server IP_0 is currently overloaded."
	result := engine.Desanitize(content, session)

	expected := "To optimize ServerDB01.users_prod, add an index. The server 192.168.1.100 is currently overloaded."

	if result.DesanitizedContent != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result.DesanitizedContent)
	}

	if result.ReplacementsCount != 3 {
		t.Errorf("Expected 3 replacements, got %d", result.ReplacementsCount)
	}
}

func TestDesanitize_NoMappings(t *testing.T) {
	engine := NewEngine()
	session := types.NewSession("test-session", "user@test.com", "engineering", 8*time.Hour)

	content := "This is regular content with no aliases."
	result := engine.Desanitize(content, session)

	if result.DesanitizedContent != content {
		t.Error("Expected content to remain unchanged")
	}

	if result.ReplacementsCount != 0 {
		t.Error("Expected 0 replacements")
	}
}

func TestDesanitize_UnmatchedAliases(t *testing.T) {
	engine := NewEngine()
	session := types.NewSession("test-session", "user@test.com", "engineering", 8*time.Hour)

	// Only add one mapping
	session.AddMapping("ServerDB01", "SERVER_0")

	// Content has an unmatched alias
	content := "Check SERVER_0 and SERVER_1 for issues."
	result := engine.Desanitize(content, session)

	// SERVER_0 should be replaced, SERVER_1 should remain unchanged
	if result.ReplacementsCount != 1 {
		t.Errorf("Expected 1 replacement, got %d", result.ReplacementsCount)
	}

	// Verify SERVER_0 was replaced and SERVER_1 remains
	expectedContent := "Check ServerDB01 and SERVER_1 for issues."
	if result.DesanitizedContent != expectedContent {
		t.Errorf("Expected:\n%s\nGot:\n%s", expectedContent, result.DesanitizedContent)
	}
}

func TestDesanitize_MultipleOccurrences(t *testing.T) {
	engine := NewEngine()
	session := types.NewSession("test-session", "user@test.com", "engineering", 8*time.Hour)

	session.AddMapping("ServerDB01", "SERVER_0")

	content := "First mention of SERVER_0, second mention of SERVER_0, third SERVER_0."
	result := engine.Desanitize(content, session)

	expected := "First mention of ServerDB01, second mention of ServerDB01, third ServerDB01."
	if result.DesanitizedContent != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result.DesanitizedContent)
	}

	// Should count all 3 replacements
	if result.ReplacementsCount != 3 {
		t.Errorf("Expected 3 replacements, got %d", result.ReplacementsCount)
	}
}

func TestDesanitize_LongAliasFirst(t *testing.T) {
	engine := NewEngine()
	session := types.NewSession("test-session", "user@test.com", "engineering", 8*time.Hour)

	// Add aliases where one is a prefix of another
	session.AddMapping("ServerDB01", "SERVER_1")
	session.AddMapping("ServerDB10", "SERVER_10")

	// SERVER_10 should be replaced before SERVER_1 to avoid partial matches
	content := "Check SERVER_10 and SERVER_1 status."
	result := engine.Desanitize(content, session)

	expected := "Check ServerDB10 and ServerDB01 status."
	if result.DesanitizedContent != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result.DesanitizedContent)
	}
}

func TestDesanitize_EmptySession(t *testing.T) {
	engine := NewEngine()
	session := types.NewSession("test-session", "user@test.com", "engineering", 8*time.Hour)

	content := "Some content with SERVER_0"
	result := engine.Desanitize(content, session)

	// With no mappings, content should remain unchanged
	if result.DesanitizedContent != content {
		t.Error("Expected content to remain unchanged with empty session")
	}
}

