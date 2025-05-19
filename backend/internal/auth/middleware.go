package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Context keys
const (
	ContextUserIDKey      = "userID"
	ContextModulesKey     = "modules"
	ContextDepartmentsKey = "departments"
	ContextUserRoleKey    = "userRole"
)

// Middleware validates JWTs and injects user identity into Gin context.
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
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextModulesKey, claims.Modules)
		c.Set(ContextDepartmentsKey, claims.Departments)
		c.Set(ContextUserRoleKey, claims.Role)
		c.Next()
	}
}

// RequireRoles returns a Gin middleware that allows access only to users with the specified roles.
func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get(ContextUserRoleKey) // use your const here

		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: missing role"})
			return
		}

		role, ok := roleValue.(string)
		if !ok || role == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: invalid role format"})
			return
		}

		for _, allowed := range allowedRoles {
			if role == allowed {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient role"})
	}
}
