// Package session provides session management for Enterprise Shield.
package session

import (
	"sync"
	"time"

	"github.com/enterprise/opencode-enterprise-shield/pkg/types"
	"github.com/google/uuid"
)

// Manager handles session lifecycle and storage.
type Manager struct {
	store      *Store
	defaultTTL time.Duration
	maxMappings int
	mu         sync.RWMutex
}

// NewManager creates a new session manager.
func NewManager(defaultTTL time.Duration, maxMappings int) *Manager {
	return &Manager{
		store:       NewStore(),
		defaultTTL:  defaultTTL,
		maxMappings: maxMappings,
	}
}

// GetOrCreate retrieves an existing session or creates a new one.
func (m *Manager) GetOrCreate(userID, department string, sessionID string) (*types.Session, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Try to get existing session
	if sessionID != "" {
		session, ok := m.store.Get(sessionID)
		if ok && !session.IsExpired() && session.UserID == userID {
			session.Touch()
			return session, false
		}
	}

	// Try to find by user ID
	existingID, ok := m.store.GetByUserID(userID)
	if ok {
		session, ok := m.store.Get(existingID)
		if ok && !session.IsExpired() {
			session.Touch()
			return session, false
		}
	}

	// Create new session
	newSessionID := "sess_" + uuid.New().String()[:12]
	session := types.NewSession(newSessionID, userID, department, m.defaultTTL)
	m.store.Set(session)

	return session, true
}

// Get retrieves a session by ID.
func (m *Manager) Get(sessionID string) (*types.Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, ok := m.store.Get(sessionID)
	if !ok || session.IsExpired() {
		return nil, false
	}
	return session, true
}

// Delete removes a session.
func (m *Manager) Delete(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store.Delete(sessionID)
}

// Clear removes all sessions for a user.
func (m *Manager) Clear(userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	sessionID, ok := m.store.GetByUserID(userID)
	if ok {
		m.store.Delete(sessionID)
	}
}

// CleanupExpired removes all expired sessions.
func (m *Manager) CleanupExpired() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.store.CleanupExpired()
}

// GetStats returns session statistics.
func (m *Manager) GetStats() SessionStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.store.GetStats()
}

// Export exports session data for backup.
func (m *Manager) Export(sessionID string) (*types.Session, bool) {
	return m.Get(sessionID)
}

// SessionStats contains session statistics.
type SessionStats struct {
	TotalSessions   int `json:"totalSessions"`
	ActiveSessions  int `json:"activeSessions"`
	ExpiredSessions int `json:"expiredSessions"`
	TotalMappings   int `json:"totalMappings"`
}

