package compliance

import (
	"testing"

	"github.com/enterprise/opencode-enterprise-shield/pkg/types"
)

func TestDetector_SSN(t *testing.T) {
	detector := NewDetector(true)

	result := detector.Scan("My SSN is 123-45-6789")

	if !result.HasViolations {
		t.Error("Expected SSN violation to be detected")
	}

	if !result.ShouldBlock {
		t.Error("Expected SSN to trigger blocking")
	}

	if len(result.Violations) != 1 {
		t.Errorf("Expected 1 violation, got %d", len(result.Violations))
	}

	if result.Violations[0].Type != "SSN" {
		t.Errorf("Expected SSN type, got %s", result.Violations[0].Type)
	}

	if result.Violations[0].Severity != types.SeverityCritical {
		t.Errorf("Expected critical severity, got %s", result.Violations[0].Severity)
	}
}

func TestDetector_CreditCard(t *testing.T) {
	detector := NewDetector(true)

	// Valid Visa card (passes Luhn)
	result := detector.Scan("Card: 4111111111111111")

	if !result.HasViolations {
		t.Error("Expected credit card violation to be detected")
	}

	if !result.ShouldBlock {
		t.Error("Expected credit card to trigger blocking")
	}
}

func TestDetector_InvalidCreditCard(t *testing.T) {
	detector := NewDetector(true)

	// Invalid card (fails Luhn)
	result := detector.Scan("Card: 4111111111111112")

	if result.HasViolations {
		t.Error("Expected invalid credit card to not trigger violation (fails Luhn)")
	}
}

func TestDetector_AWSKey(t *testing.T) {
	detector := NewDetector(true)

	result := detector.Scan("AWS key: AKIAIOSFODNN7EXAMPLE")

	if !result.HasViolations {
		t.Error("Expected AWS key violation to be detected")
	}
}

func TestDetector_GitHubToken(t *testing.T) {
	detector := NewDetector(true)

	result := detector.Scan("Token: ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")

	if !result.HasViolations {
		t.Error("Expected GitHub token violation to be detected")
	}
}

func TestDetector_PrivateKey(t *testing.T) {
	detector := NewDetector(true)

	content := `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA2Z3qX2B...
-----END RSA PRIVATE KEY-----`

	result := detector.Scan(content)

	if !result.HasViolations {
		t.Error("Expected private key violation to be detected")
	}
}

func TestDetector_Password(t *testing.T) {
	detector := NewDetector(true)

	result := detector.Scan("password=mysecretpassword123")

	if !result.HasViolations {
		t.Error("Expected password violation to be detected")
	}
}

func TestDetector_OpenAIKey(t *testing.T) {
	detector := NewDetector(true)

	result := detector.Scan("sk-1234567890123456789012345678901234567890123456ab")

	if !result.HasViolations {
		t.Error("Expected OpenAI key violation to be detected")
	}
}

func TestDetector_NoViolations(t *testing.T) {
	detector := NewDetector(true)

	result := detector.Scan("How do I optimize a database query?")

	if result.HasViolations {
		t.Error("Expected no violations for generic content")
	}

	if result.ShouldBlock {
		t.Error("Expected no blocking for generic content")
	}
}

func TestDetector_MultipleViolations(t *testing.T) {
	detector := NewDetector(true)

	content := "SSN: 123-45-6789, Card: 4111111111111111, Key: AKIAIOSFODNN7EXAMPLE"
	result := detector.Scan(content)

	if !result.HasViolations {
		t.Error("Expected violations to be detected")
	}

	if len(result.Violations) < 3 {
		t.Errorf("Expected at least 3 violations, got %d", len(result.Violations))
	}
}

func TestLuhnValidation(t *testing.T) {
	tests := []struct {
		number string
		valid  bool
	}{
		{"4111111111111111", true},  // Visa test card
		{"5500000000000004", true},  // Mastercard test card
		{"340000000000009", true},   // Amex test card
		{"4111111111111112", false}, // Invalid
		{"1234567890123456", false}, // Invalid
	}

	for _, test := range tests {
		result := ValidateLuhn(test.number)
		if result != test.valid {
			t.Errorf("Luhn(%s) = %v, want %v", test.number, result, test.valid)
		}
	}
}

func TestGetCardType(t *testing.T) {
	tests := []struct {
		number   string
		expected string
	}{
		{"4111111111111111", "visa"},
		{"5500000000000004", "mastercard"},
		{"340000000000009", "amex"},
		{"6011000000000004", "discover"},
		{"1234567890123456", "unknown"},
	}

	for _, test := range tests {
		result := GetCardType(test.number)
		if result != test.expected {
			t.Errorf("GetCardType(%s) = %s, want %s", test.number, result, test.expected)
		}
	}
}

