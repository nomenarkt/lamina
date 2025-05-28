// Package auth handles HTTP routing for authentication.
package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// RegisterRoutes registers the authentication endpoints for signup and login.
func RegisterRoutes(router *gin.RouterGroup, db *sqlx.DB, service ServiceInterface) {
	if service == nil {
		service = NewService(NewAuthRepository(db))
	}

	router.POST("/auth/signup", func(c *gin.Context) {
		service.Signup(c)
	})

	router.POST("/auth/login", func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tokens, err := service.Login(c.Request.Context(), req)
		if err != nil {
			if errors.Is(err, ErrUnconfirmedAccount) {
				c.JSON(http.StatusForbidden, gin.H{"error": "account not confirmed"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			}
			return
		}
		c.JSON(http.StatusOK, tokens)
	})

	router.GET("/auth/confirm/:token", func(c *gin.Context) {
		token := c.Param("token")
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:3000"
		}

		err := service.ConfirmRegistration(c.Request.Context(), token)
		if err != nil {
			log.Printf("📛 Email confirmation error: %v", err)
			reason := "invalid"
			switch err.Error() {
			case "token expired":
				reason = "expired"
			case "user already confirmed":
				reason = "already-confirmed"
			}
			c.Redirect(http.StatusFound, fmt.Sprintf("%s/confirm-error?reason=%s", frontendURL, reason))
			return
		}

		// Optional JSON fallback for fetch-based tools or Postman
		if c.GetHeader("Accept") == "application/json" {
			c.JSON(http.StatusOK, gin.H{"message": "account confirmed"})
			return
		}

		c.Redirect(http.StatusFound, fmt.Sprintf("%s/email-confirmed", frontendURL))
	})

	router.POST("/auth/complete-invite", func(c *gin.Context) {
		var req struct {
			Token string `json:"token" binding:"required"`
			PasswordPayload
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		if err := req.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := service.CompleteInvite(c.Request.Context(), req.Token, req.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  resp.AccessToken,
			"refresh_token": resp.RefreshToken,
		})
	})

	router.POST("/auth/resend-confirmation", func(c *gin.Context) {
		var req struct {
			Email string `json:"email" binding:"required,email"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("❌ Invalid email input: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
			return
		}

		log.Printf("🔁 Resend confirmation requested for email: %s", req.Email)

		if err := service.ResendConfirmation(c.Request.Context(), req.Email); err != nil {
			log.Printf("❌ Resend failed: %v", err)
			switch err.Error() {
			case "user already confirmed or invalid status":
				c.JSON(http.StatusBadRequest, gin.H{"error": "User already confirmed"})
			case "user not found":
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			case "resend allowed only for internal users":
				c.JSON(http.StatusForbidden, gin.H{"error": "Resend allowed only for internal users"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not resend confirmation"})
			}
			return
		}

		log.Printf("✅ Resend email dispatched successfully for %s", req.Email)
		c.JSON(http.StatusOK, gin.H{"message": "Confirmation email resent"})
	})
}
