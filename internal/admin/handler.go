package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/nomenarkt/lamina/common/utils"
	"github.com/nomenarkt/lamina/internal/user"
)

func RegisterRoutes(r *gin.RouterGroup, db *sqlx.DB) {
	adminRepo := NewAdminRepository(db)
	hasher := &utils.BcryptHasher{}
	adminService := NewAdminService(adminRepo, hasher)

	r.POST("/admin/create-user", func(c *gin.Context) {
		// Must be an admin
		claims, ok := c.Get("user")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userClaims := claims.(*user.UserClaims)
		if userClaims.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin only"})
			return
		}

		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		err := adminService.CreateUser(c.Request.Context(), req, userClaims.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "user created"})
	})
}
