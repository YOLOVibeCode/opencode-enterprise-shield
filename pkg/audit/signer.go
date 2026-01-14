// Package audit provides entry signing for audit integrity.
package audit

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"

	"github.com/enterprise/opencode-enterprise-shield/pkg/types"
)

// Signer provides Ed25519 signing for audit entries.
type Signer struct {
	publicKey  ed25519.PublicKey
	privateKey ed25519.PrivateKey
}

// NewSigner creates a new signer with a generated key pair.
func NewSigner() *Signer {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	return &Signer{
		publicKey:  pub,
		privateKey: priv,
	}
}

// NewSignerWithKey creates a signer with an existing private key.
func NewSignerWithKey(privateKey ed25519.PrivateKey) *Signer {
	return &Signer{
		publicKey:  privateKey.Public().(ed25519.PublicKey),
		privateKey: privateKey,
	}
}

// Sign signs an audit entry and returns the base64-encoded signature.
func (s *Signer) Sign(entry types.AuditEntry) (string, error) {
	// Create canonical representation (deterministic JSON)
	canonical := struct {
		Timestamp        string `json:"timestamp"`
		UserID           string `json:"userId"`
		SessionID        string `json:"sessionId"`
		RequestHash      string `json:"requestHash"`
		Action           string `json:"action"`
		PreviousEntryHash string `json:"previousEntryHash"`
	}{
		Timestamp:         entry.Timestamp.UTC().Format("2006-01-02T15:04:05Z07:00"),
		UserID:            entry.UserID,
		SessionID:         entry.SessionID,
		RequestHash:       entry.RequestHash,
		Action:            string(entry.Action),
		PreviousEntryHash: entry.PreviousEntryHash,
	}

	data, err := json.Marshal(canonical)
	if err != nil {
		return "", err
	}

	// Sign the canonical data
	signature := ed25519.Sign(s.privateKey, data)

	return "ed25519:" + base64.StdEncoding.EncodeToString(signature), nil
}

// Verify verifies an audit entry signature.
func (s *Signer) Verify(entry types.AuditEntry, signatureStr string) bool {
	// Extract the base64 signature
	if len(signatureStr) <= 8 || signatureStr[:8] != "ed25519:" {
		return false
	}

	signature, err := base64.StdEncoding.DecodeString(signatureStr[8:])
	if err != nil {
		return false
	}

	// Create canonical representation
	canonical := struct {
		Timestamp         string `json:"timestamp"`
		UserID            string `json:"userId"`
		SessionID         string `json:"sessionId"`
		RequestHash       string `json:"requestHash"`
		Action            string `json:"action"`
		PreviousEntryHash string `json:"previousEntryHash"`
	}{
		Timestamp:         entry.Timestamp.UTC().Format("2006-01-02T15:04:05Z07:00"),
		UserID:            entry.UserID,
		SessionID:         entry.SessionID,
		RequestHash:       entry.RequestHash,
		Action:            string(entry.Action),
		PreviousEntryHash: entry.PreviousEntryHash,
	}

	data, err := json.Marshal(canonical)
	if err != nil {
		return false
	}

	return ed25519.Verify(s.publicKey, data, signature)
}

// GetPublicKey returns the public key.
func (s *Signer) GetPublicKey() ed25519.PublicKey {
	return s.publicKey
}

// GetPublicKeyBase64 returns the public key as base64.
func (s *Signer) GetPublicKeyBase64() string {
	return base64.StdEncoding.EncodeToString(s.publicKey)
}

