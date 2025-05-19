// Package auth defines data models related to authentication operations.
package auth

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
