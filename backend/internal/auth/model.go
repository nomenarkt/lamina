// Package auth defines data models related to authentication operations.
package auth

import (
	"errors"
	"strings"
)

// SignupRequest represents the expected payload for a user signup request.
type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest represents the expected payload for a user login request.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Response contains JWT tokens returned upon successful authentication.
type Response struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// PasswordPayload is used for password confirmation during signup or invite completion.
type PasswordPayload struct {
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}

// Validate ensures that both passwords match and are not empty.
func (p PasswordPayload) Validate() error {
	if strings.TrimSpace(p.Password) == "" || strings.TrimSpace(p.ConfirmPassword) == "" {
		return errors.New("password and confirm password must not be empty")
	}
	if p.Password != p.ConfirmPassword {
		return errors.New("passwords do not match")
	}
	return nil
}

// User represents the structure of a user in the system's database.
type User struct {
	ID           int64  `db:"id"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
	Role         string `db:"role"`
	Status       string `db:"status"`
	FirstName    string `db:"first_name"`
	LastName     string `db:"last_name"`
	PhoneNumber  string `db:"phone_number"`
	Department   string `db:"department"`
	CreatedAt    string `db:"created_at"`
}
