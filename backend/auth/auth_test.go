package auth_test

import (
	"backend/auth"
	"testing"
)

func TestAuthenticateUser(t *testing.T) {
	t.Run("Valid Credentials", func(t *testing.T) {
		sessionID, err := auth.Authenticate("admin", "1234")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if sessionID == "" {
			t.Errorf("Expected session ID, got empty string")
		}

		if !auth.IsValidSession(sessionID) {
			t.Errorf("Expected session to be valid")
		}
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		sessionID, err := auth.Authenticate("admin", "wrong")
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if sessionID != "" {
			t.Errorf("Expected empty session ID, got %s", sessionID)
		}
	})
	
	t.Run("Invalid Session ID", func(t *testing.T) {
		if auth.IsValidSession("invalid-session") {
			t.Errorf("Expected session to be invalid")
		}
	})
}
