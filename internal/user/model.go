// Package user contains data models related to user identity and profile management.
package user

import "time"

// User represents a registered user in the system.
type User struct {
	ID                int64     `db:"id"`
	CompanyID         *int      `db:"company_id"` // Nullable: link to company
	Email             string    `db:"email"`
	PasswordHash      string    `db:"password_hash"`
	Role              string    `db:"role"`
	Status            string    `db:"status"`
	ConfirmationToken *string   `db:"confirmation_token"` // Nullable: for email confirmation flows
	CreatedAt         time.Time `db:"created_at"`         // Timestamp of registration
	FullName          *string   `db:"full_name"`          // Nullable: display name
}

// Claims holds the minimal identity claims extracted from a JWT token.
type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

// UpdateProfileRequest defines the payload for updating user profile data.
type UpdateProfileRequest struct {
	FullName  string `json:"full_name" binding:"required"`
	CompanyID *int   `json:"company_id,omitempty"` // Optional: for system-internal users only
}
