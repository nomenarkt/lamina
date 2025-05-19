// Package admin defines data models for the admin functionality.
package admin

// CreateUserRequest represents the JSON payload to create a new user account.
type CreateUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required"`
}
