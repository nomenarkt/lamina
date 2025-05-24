// Package admin defines data models for the admin functionality.
package admin

// CreateUserRequest is used by admin handlers to invite new users.
// It supports optional access expiration for external users.
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role"`               // optional; defaults to "user"
	Duration string `json:"duration,omitempty"` // optional; parsed like "1w", "2h", etc.
}
