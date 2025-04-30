package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/nomenarkt/lamina/common/utils"
)

type UserHandler struct {
	service *UserService
}

func NewUserHandler(s *UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	// Extract userID from context (injected by auth middleware)
	idStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userID, err := strconv.ParseInt(idStr.(string), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.service.GetMe(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

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
