package utils

type PasswordHasher interface {
	HashPassword(password string) (string, error)
}
