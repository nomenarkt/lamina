// Package auth handles HTTP routing for authentication.
package auth

import (
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
		if err := service.ConfirmRegistration(c.Request.Context(), token); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:5173"
		}
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/login")
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
