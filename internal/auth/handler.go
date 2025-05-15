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
		tokens, err := service.Login(c.Request.Context(), req) // match interface from service.go
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, tokens)
	})

	router.GET("/auth/confirm/:token", func(c *gin.Context) {
		token := c.Param("token")
		err := service.(*Service).ConfirmRegistration(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:5173"
		}
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/login")
	})

}
