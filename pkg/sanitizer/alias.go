// Package sanitizer provides alias generation for sanitization.
package sanitizer

import (
	"fmt"

	"github.com/enterprise/opencode-enterprise-shield/pkg/types"
)

// AliasGenerator generates unique aliases for sensitive data.
type AliasGenerator struct{}

// NewAliasGenerator creates a new alias generator.
func NewAliasGenerator() *AliasGenerator {
	return &AliasGenerator{}
}

// Generate creates a new alias using the format {PREFIX}_{COUNTER}.
// Examples: SERVER_0, TABLE_0, IP_1
func (g *AliasGenerator) Generate(session *types.Session, prefix string) string {
	counter := session.GetNextCounter(prefix)
	return fmt.Sprintf("%s_%d", prefix, counter)
}

