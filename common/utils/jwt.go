package utils

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ✅ Struct defining what your JWT will contain
type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// ✅ Function to generate Access and Refresh Tokens
func GenerateTokens(userID int64, email string, role string) (string, string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", "", fmt.Errorf("JWT_SECRET is not set in the environment")
	}

	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	refreshToken := "SIMPLE-REFRESH-TOKEN-STUB"

	return accessToken, refreshToken, nil
}

// ✅ NEW: ParseToken to validate and extract Claims from JWT
// JWT Token Parsing Logic
// -----------------------
// The "Missing HMAC secret" error occurs when the server attempts to parse a JWT
// but cannot find the JWT_SECRET environment variable.
//
// This is *not* related to client-side HMAC signatures or missing request headers.
//
// ✅ To fix this error:
// 1. Set JWT_SECRET in your .env file or server environment.
// 2. Make sure Docker or your deployment environment passes it at runtime.
// 3. Avoid Docker layer caching issues that skip source recompilation.
//
// This error comes from jwt.ParseWithClaims when given an empty signing key.

func ParseToken(tokenString string) (*Claims, error) {
	fmt.Println("JWT_SECRET (runtime):", os.Getenv("JWT_SECRET"))

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is not set in the environment")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

// ✅ NEW: Middleware to check auth and inject user info into Gin Context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"error": "Missing or invalid Authorization header"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		claims, err := ParseToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Inject into context
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
