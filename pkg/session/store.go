// Package session provides the session store for Enterprise Shield.
package session

import (
	"sync"
	"time"

	"github.com/enterprise/opencode-enterprise-shield/pkg/types"
)

// Store is an in-memory session store.
type Store struct {
	sessions map[string]*types.Session
	userIndex map[string]string // userID -> sessionID
	mu       sync.RWMutex
}

// NewStore creates a new session store.
func NewStore() *Store {
	return &Store{
		sessions:  make(map[string]*types.Session),
		userIndex: make(map[string]string),
	}
}

// Get retrieves a session by ID.
func (s *Store) Get(sessionID string) (*types.Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[sessionID]
	return session, ok
}

// GetByUserID retrieves a session ID by user ID.
func (s *Store) GetByUserID(userID string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sessionID, ok := s.userIndex[userID]
	return sessionID, ok
}

// Set stores a session.
func (s *Store) Set(session *types.Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove old session for this user if exists
	if oldSessionID, ok := s.userIndex[session.UserID]; ok {
		delete(s.sessions, oldSessionID)
	}

	s.sessions[session.SessionID] = session
	s.userIndex[session.UserID] = session.SessionID
}

// Delete removes a session.
func (s *Store) Delete(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if session, ok := s.sessions[sessionID]; ok {
		delete(s.userIndex, session.UserID)
		delete(s.sessions, sessionID)
	}
}

// CleanupExpired removes all expired sessions.
func (s *Store) CleanupExpired() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	count := 0
	now := time.Now()

	for sessionID, session := range s.sessions {
		if now.After(session.ExpiresAt) {
			delete(s.userIndex, session.UserID)
			delete(s.sessions, sessionID)
			count++
		}
	}

	return count
}

// GetStats returns store statistics.
func (s *Store) GetStats() SessionStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := SessionStats{
		TotalSessions: len(s.sessions),
	}

	now := time.Now()
	for _, session := range s.sessions {
		if now.Before(session.ExpiresAt) && session.Status == types.SessionActive {
			stats.ActiveSessions++
		} else {
			stats.ExpiredSessions++
		}
		stats.TotalMappings += len(session.Mappings)
	}

	return stats
}

// ListAll returns all session IDs.
func (s *Store) ListAll() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids := make([]string, 0, len(s.sessions))
	for id := range s.sessions {
		ids = append(ids, id)
	}
	return ids
}

