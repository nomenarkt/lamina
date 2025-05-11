package auth

type SignupRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

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
