package user

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nomenarkt/lamina/common/utils"
	"github.com/nomenarkt/lamina/internal/middleware"
)

type UserServiceInterface interface {
	GetMe(ctx context.Context, id int64) (*User, error)
	ListUsers(ctx context.Context) ([]User, error)
	UpdateUserProfile(ctx context.Context, userID int64, req UpdateProfileRequest) error
}

type UserHandler struct {
	service UserServiceInterface
}

func NewUserHandler(svc UserServiceInterface) *UserHandler {
	return &UserHandler{service: svc}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID, exists := utils.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.service.GetMe(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) ListAll(c *gin.Context) {
	users, err := h.service.ListUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func RegisterRoutes(router *gin.RouterGroup, h *UserHandler) {
	group := router.Group("/user")
	group.GET("/me", h.GetMe)
	group.GET("/", h.ListAll)
	group.PUT("/profile", h.UpdateProfile)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.service.ListUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format", "details": err.Error()})
		return
	}

	if err := h.service.UpdateUserProfile(c, userID, req); err != nil {
		// ✅ Log the actual error and return it in the response
		log.Printf("❌ UpdateUserProfile error for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated"})
}
