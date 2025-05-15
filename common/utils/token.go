package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateSecureToken returns a secure random hex string of the given length.
func GenerateSecureToken(n int) (string, error) {

	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
