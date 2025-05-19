package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nomenarkt/lamina/internal/user"
)

// Claims = JWT payload returned to frontend.
type Claims struct {
	UserID      int64                   `json:"userID"`
	Email       string                  `json:"email"`
	Departments []user.DepartmentAccess `json:"departments"`
	Modules     []string                `json:"modules"`
	Role        string                  `json:"role"`
	jwt.RegisteredClaims
}

// GenerateTokens returns access + refresh tokens for a user.
func GenerateTokens(secret, refreshSecret string, u user.User) (string, string, error) {
	claims := Claims{
		UserID:      u.ID,
		Email:       u.Email,
		Departments: u.Departments,
		Modules:     u.Modules,
		Role:        u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), //Add(24 * time.Hour))
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
	})
	refreshToken, err := refresh.SignedString([]byte(refreshSecret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// GenerateTokensFromEnv uses env vars for signing keys.
func GenerateTokensFromEnv(u user.User) (string, string, error) {
	return GenerateTokens(os.Getenv("JWT_SECRET"), os.Getenv("JWT_REFRESH_SECRET"), u)
}
