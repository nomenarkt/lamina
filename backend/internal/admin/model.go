// Package admin defines data models for the admin functionality.
package admin

// CreateUserRequest represents the JSON payload to create a new user account.
type CreateUserRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}
