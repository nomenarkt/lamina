package user

import "time"

type User struct {
	ID                int64     `db:"id"`
	CompanyID         *int      `db:"company_id"` // nullable
	Email             string    `db:"email"`
	PasswordHash      string    `db:"password_hash"`
	Role              string    `db:"role"`
	Status            string    `db:"status"`
	ConfirmationToken *string   `db:"confirmation_token"` // nullable
	CreatedAt         time.Time `db:"created_at"`         // timestamp
	FullName          *string   `db:"full_name"`          // nullable
}

type UserClaims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

type UpdateProfileRequest struct {
	FullName  string `json:"full_name" binding:"required"`
	CompanyID *int   `json:"company_id,omitempty"` // only allowed for internal users
}
