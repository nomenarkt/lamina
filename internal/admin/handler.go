package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/nomenarkt/lamina/common/utils"
	"github.com/nomenarkt/lamina/internal/middleware"
	"github.com/nomenarkt/lamina/internal/user"
)

func RegisterRoutes(r *gin.RouterGroup, db *sqlx.DB) {
	adminRepo := NewAdminRepository(db)
	hasher := &utils.BcryptHasher{}
	adminService := NewAdminService(adminRepo, hasher)

	adminGroup := r.Group("/admin", middleware.JWTMiddleware(), middleware.RequireRoles("admin"))

	adminGroup.POST("/create-user", func(c *gin.Context) {
		claimsVal, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authentication context"})
			return
		}

		claims, ok := claimsVal.(*user.UserClaims)
		if !ok || claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user claims"})
			return
		}

		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input format"})
			return
		}

		if err := adminService.CreateUser(c.Request.Context(), req, claims.Email); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "user successfully created"})
	})
}
