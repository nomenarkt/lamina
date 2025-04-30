package admin

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/nomenarkt/lamina/common/utils"
	"golang.org/x/crypto/bcrypt"
)

func RegisterRoutes(router *gin.RouterGroup, db *sqlx.DB) {
	adminService := NewAdminService(NewAdminRepository(db))

	adminGroup := router.Group("/admin")
	{

		adminGroup.POST("/create-user", func(c *gin.Context) {
			var req CreateUserRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// ✅ Extract UserID from JWT
			_, exists := utils.GetUserIDFromContext(c)
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				return
			}

			// ✅ Enforce admin-only access (🔐 NEW step)
			userRole := c.GetString("role")
			if userRole != "admin" {
				c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can create users"})
				return
			}

			// ✅ Domain check (already in your file)
			if !strings.HasSuffix(req.Email, "@madagascarairlines.com") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Email must belong to Madagascar Airlines"})
				return
			}

			// ✅ Hash password
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
				return
			}

			// ✅ Call service to save user
			if err := adminService.CreateUser(c, req, string(hashedPassword)); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
		})

	}
}
