// Package utils provides utility functions such as password hashing.
package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// BcryptHasher is a concrete implementation of the PasswordHasher interface using bcrypt.
type BcryptHasher struct{}

// HashPassword returns a bcrypt hashed password from the given plain-text password.
func (b *BcryptHasher) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
