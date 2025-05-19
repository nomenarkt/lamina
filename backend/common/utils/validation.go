package utils

import "regexp"

var emailRegex = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

// IsValidEmail validates the format of an email address.
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}
