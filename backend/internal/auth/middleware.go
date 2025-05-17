package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

// Context keys
const (
	ContextUserIDKey = "userID"
	ContextRoleKey   = "role"
)

// isDebugEnabled returns true if APP_DEBUG env is explicitly set to "true".
func isDebugEnabled() bool {
	return os.Getenv("APP_DEBUG") == "true"
}

// Middleware validates JWT tokens and injects user information into the context.
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization Header"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization Header"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			if isDebugEnabled() {
				log.Printf("[JWT DEBUG] Token Header alg: %v\n", t.Header["alg"])
			}

			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			if isDebugEnabled() {
				log.Printf("[JWT DEBUG] Token parse error: %v\n", err)
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if isDebugEnabled() {
			log.Printf("[JWT DEBUG] Token is valid: %v\n", token.Valid)
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			if isDebugEnabled() {
				log.Println("[JWT DEBUG] Token claims type mismatch")
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		if isDebugEnabled() {
			log.Printf("[JWT DEBUG] Token Claims: userID=%d, email=%s, role=%s\n", claims.UserID, claims.Email, claims.Role)
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextRoleKey, claims.Role)
		c.Next()
	}
}
