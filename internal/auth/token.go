package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret     = os.Getenv("JWT_SECRET")         // access token secret
	refreshSecret = os.Getenv("JWT_REFRESH_SECRET") // optional: use a separate one
)

// Claims represents JWT claims used by the application.
type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateTokens(userID int64, email string, role string) (accessToken string, refreshToken string, err error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err = token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", err
	}

	refreshClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refresh.SignedString([]byte(refreshSecret)) // or reuse jwtSecret
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
