package utils

import (
	"os"
	"testing"
)

func TestJWTGenerateAndParse(t *testing.T) {
	// Setup: Fake secret for testing
	os.Setenv("JWT_SECRET", "testsecret123")

	userID := int64(3190)
	email := "m.rakotoarison@madagascarairlines.com"
	role := "admin"

	// Generate token
	accessToken, _, err := GenerateTokens(userID, email, role)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Parse token
	claims, err := ParseToken(accessToken)
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	// Validate extracted claims
	if claims.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("Expected Email %s, got %s", email, claims.Email)
	}
	if claims.Role != role {
		t.Errorf("Expected Role %s, got %s", role, claims.Role)
	}
}
