// Package auth handles HTTP routing for authentication.
package auth

import (
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
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
}
