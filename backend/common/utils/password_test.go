package utils

import (
	"testing"
)

func TestPasswordHashingAndValidation(t *testing.T) {
	// Step 1: Choose a test password
	password := "supersecure123"

	// Step 2: Hash the password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Step 3: Check that the original password matches the hashed password
	err = CheckPasswordHash(password, hashedPassword)
	if err != nil {
		t.Errorf("Password validation failed: %v", err)
	}

	// Step 4: Check that a wrong password does NOT match
	wrongPassword := "wrongpassword"
	err = CheckPasswordHash(wrongPassword, hashedPassword)
	if err == nil {
		t.Errorf("Validation should have failed for wrong password")
	}
}
