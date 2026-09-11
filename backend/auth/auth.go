package auth

import (
	"errors"
	"sync"
	"time"
)

var (
	sessions = make(map[string]time.Time)
	mu       sync.RWMutex
)

// Authenticate verifies credentials and returns a new session ID if successful.
func Authenticate(username, password string) (string, error) {
	if username == "admin" && password == "1234" {
		sessionID := generateSessionID()
		mu.Lock()
		sessions[sessionID] = time.Now().Add(1 * time.Hour)
		mu.Unlock()
		return sessionID, nil
	}
	return "", errors.New("invalid credentials")
}

// IsValidSession checks if a session ID exists and is active.
func IsValidSession(sessionID string) bool {
	mu.RLock()
	expiry, exists := sessions[sessionID]
	mu.RUnlock()

	if !exists {
		return false
	}
	if time.Now().After(expiry) {
		mu.Lock()
		delete(sessions, sessionID)
		mu.Unlock()
		return false
	}
	return true
}

func generateSessionID() string {
	// Simple UUID mock for this exercise
	return "session-" + time.Now().Format("20060102150405.000")
}
