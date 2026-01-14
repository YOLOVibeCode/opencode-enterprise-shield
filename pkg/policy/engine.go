// Package policy provides the policy engine for access control.
package policy

import (
	"sync"

	"github.com/enterprise/opencode-enterprise-shield/pkg/types"
)

// Engine evaluates access policies for users.
type Engine struct {
	policies       map[string]*types.UserPolicy // userID -> policy
	deptPolicies   map[string]*types.UserPolicy // department -> default policy
	defaultPolicy  *types.UserPolicy
	mu             sync.RWMutex
}

// NewEngine creates a new policy engine.
func NewEngine() *Engine {
	return &Engine{
		policies:      make(map[string]*types.UserPolicy),
		deptPolicies:  make(map[string]*types.UserPolicy),
		defaultPolicy: DefaultPolicy(),
	}
}

// DefaultPolicy returns the default policy for users without specific policies.
func DefaultPolicy() *types.UserPolicy {
	return &types.UserPolicy{
		PolicyID:           "default",
		AccessLevel:        types.AccessSanitizedOnly,
		AllowedProviders:   []string{"openai", "anthropic", "azure_openai", "google"},
		DailyRequestLimit:  500,
		HourlyRequestLimit: 50,
		Enabled:            true,
	}
}

// Evaluate evaluates the policy for a request context.
func (e *Engine) Evaluate(ctx PolicyContext) types.PolicyDecision {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Get effective policy
	policy := e.getEffectivePolicy(ctx.UserID, ctx.Department)

	// Check if policy is enabled
	if !policy.Enabled {
		return types.PolicyDecision{
			Action:        types.ActionBlock,
			Reason:        "User policy is disabled",
			PolicyApplied: policy.PolicyID,
		}
	}

	// Check access level
	switch policy.AccessLevel {
	case types.AccessBlocked:
		return types.PolicyDecision{
			Action:        types.ActionBlock,
			Reason:        "User access is blocked",
			PolicyApplied: policy.PolicyID,
		}
	case types.AccessUnrestricted:
		return types.PolicyDecision{
			Action:        types.ActionAllow,
			Reason:        "User has unrestricted access",
			PolicyApplied: policy.PolicyID,
		}
	case types.AccessSanitizedOnly:
		// Check provider allowlist
		if len(policy.AllowedProviders) > 0 && ctx.Provider != "" {
			allowed := false
			for _, p := range policy.AllowedProviders {
				if p == ctx.Provider {
					allowed = true
					break
				}
			}
			if !allowed {
				return types.PolicyDecision{
					Action:        types.ActionBlock,
					Reason:        "Provider not in allowed list",
					PolicyApplied: policy.PolicyID,
				}
			}
		}

		return types.PolicyDecision{
			Action:               types.ActionAllowWithSanitization,
			Reason:               "Request requires sanitization",
			PolicyApplied:        policy.PolicyID,
			RequiredSanitization: policy.RequiredRules,
		}
	}

	// Default: allow with sanitization
	return types.PolicyDecision{
		Action:        types.ActionAllowWithSanitization,
		PolicyApplied: policy.PolicyID,
	}
}

// getEffectivePolicy returns the effective policy for a user.
func (e *Engine) getEffectivePolicy(userID, department string) *types.UserPolicy {
	// Check user-specific policy first
	if policy, ok := e.policies[userID]; ok {
		return policy
	}

	// Check department policy
	if department != "" {
		if policy, ok := e.deptPolicies[department]; ok {
			return policy
		}
	}

	// Return default policy
	return e.defaultPolicy
}

// SetUserPolicy sets or updates a user's policy.
func (e *Engine) SetUserPolicy(userID string, policy *types.UserPolicy) {
	e.mu.Lock()
	defer e.mu.Unlock()

	policy.UserID = userID
	e.policies[userID] = policy
}

// SetDepartmentPolicy sets the default policy for a department.
func (e *Engine) SetDepartmentPolicy(department string, policy *types.UserPolicy) {
	e.mu.Lock()
	defer e.mu.Unlock()

	policy.Department = department
	e.deptPolicies[department] = policy
}

// GetUserPolicy retrieves a user's policy.
func (e *Engine) GetUserPolicy(userID string) (*types.UserPolicy, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	policy, ok := e.policies[userID]
	return policy, ok
}

// DeleteUserPolicy removes a user's policy.
func (e *Engine) DeleteUserPolicy(userID string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	delete(e.policies, userID)
}

// PolicyContext contains context for policy evaluation.
type PolicyContext struct {
	UserID     string
	Department string
	Provider   string
	Content    string
	SourceIP   string
}

