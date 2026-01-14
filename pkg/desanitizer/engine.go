// Package desanitizer provides the desanitization engine for restoring original values.
package desanitizer

import (
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/enterprise/opencode-enterprise-shield/pkg/types"
)

// Engine is the desanitization engine that restores original values from aliases.
type Engine struct{}

// NewEngine creates a new desanitization engine.
func NewEngine() *Engine {
	return &Engine{}
}

// Desanitize replaces aliases in content with original values from the session.
func (e *Engine) Desanitize(content string, session *types.Session) types.DesanitizationResult {
	startTime := time.Now()

	result := types.DesanitizationResult{
		DesanitizedContent: content,
		ReplacementsCount:  0,
		UnmatchedAliases:   make([]string, 0),
	}

	if len(session.ReverseMappings) == 0 {
		result.ProcessingTimeMs = time.Since(startTime).Milliseconds()
		return result
	}

	workingContent := content

	// Build list of aliases to replace, sorted by length (longest first)
	// This prevents partial replacements (e.g., SERVER_10 before SERVER_1)
	aliases := make([]string, 0, len(session.ReverseMappings))
	for alias := range session.ReverseMappings {
		aliases = append(aliases, alias)
	}
	sort.Slice(aliases, func(i, j int) bool {
		return len(aliases[i]) > len(aliases[j])
	})

	// Build a regex pattern that matches any alias
	// Using word boundaries to avoid partial matches
	if len(aliases) > 0 {
		// Escape special regex characters in aliases
		escapedAliases := make([]string, len(aliases))
		for i, alias := range aliases {
			escapedAliases[i] = regexp.QuoteMeta(alias)
		}

		pattern := `\b(` + strings.Join(escapedAliases, "|") + `)\b`
		re, err := regexp.Compile(pattern)
		if err != nil {
			// Fallback to simple string replacement
			for _, alias := range aliases {
				original, ok := session.GetOriginal(alias)
				if !ok {
					result.UnmatchedAliases = append(result.UnmatchedAliases, alias)
					continue
				}
				if strings.Contains(workingContent, alias) {
					workingContent = strings.ReplaceAll(workingContent, alias, original)
					result.ReplacementsCount++
				}
			}
		} else {
			// Use regex replacement
			workingContent = re.ReplaceAllStringFunc(workingContent, func(match string) string {
				original, ok := session.GetOriginal(match)
				if !ok {
					result.UnmatchedAliases = append(result.UnmatchedAliases, match)
					return match
				}
				result.ReplacementsCount++
				return original
			})
		}
	}

	result.DesanitizedContent = workingContent
	result.ProcessingTimeMs = time.Since(startTime).Milliseconds()

	return result
}

// DesanitizeWithContext performs desanitization while preserving JSON structure.
func (e *Engine) DesanitizeWithContext(content string, session *types.Session, preserveJSON bool) types.DesanitizationResult {
	// For now, use the standard desanitization
	// Future enhancement: parse JSON and selectively replace within string values
	return e.Desanitize(content, session)
}

