// Package admin handles the admin routes and request processing.
package admin

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nomenarkt/lamina/internal/middleware"
)

// RegisterRoutes sets up the admin routes and applies JWT and role-based middleware.
func RegisterRoutes(r *gin.RouterGroup, adminService *Service) {
	adminGroup := r.Group("/admin", middleware.JWTMiddleware(), middleware.RequireRoles("admin"))

	adminGroup.POST("/create-user", func(c *gin.Context) {
		claimsVal, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authentication context"})
			return
		}

		claims, ok := claimsVal.(*middleware.CustomClaims)
		if !ok || claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user claims"})
			return
		}

		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input format"})
			return
		}

		log.Printf("Received request to create user: %+v", req)

		err := adminService.InviteUser(c.Request.Context(), req, claims.Email)
		if err != nil {
			log.Printf("CreateUser failed: %v", err)
			if strings.Contains(err.Error(), "duplicate key value") && strings.Contains(err.Error(), "users_email_key") {
				c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "user successfully created"})
	})
}
