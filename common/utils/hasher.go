// Package utils provides utility interfaces and implementations for security utilities.
package utils

// PasswordHasher defines the interface for hashing passwords securely.
type PasswordHasher interface {
	HashPassword(password string) (string, error)
}
