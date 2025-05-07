package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterRoutes(router *gin.RouterGroup, db *sqlx.DB, service AuthServiceInterface) {
	if service == nil {
		service = NewAuthService(NewAuthRepository(db))
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
}
