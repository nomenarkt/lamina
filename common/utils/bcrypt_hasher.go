package utils

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct{}

func (b *BcryptHasher) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
