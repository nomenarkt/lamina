package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/nomenarkt/lamina/common/utils"
)

func RegisterRoutes(router *gin.RouterGroup, db *sqlx.DB) {
	userService := NewUserService(NewUserRepository(db))

	users := router.Group("/user")
	{
		users.GET("/me", func(c *gin.Context) {
			userID, exists := utils.GetUserIDFromContext(c)
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				return
			}

			u, err := userService.GetProfile(c, userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, u)
		})

		users.GET("/", func(c *gin.Context) {
			allUsers, err := userService.ListUsers(c)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, allUsers)
		})
	}
}
