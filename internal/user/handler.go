// Package user provides HTTP handlers for user-related operations such as profile retrieval and updates.
package user

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nomenarkt/lamina/common/utils"
	"github.com/nomenarkt/lamina/internal/middleware"
)

// ServiceInterface defines the contract for user service logic.
type ServiceInterface interface {
	GetMe(ctx context.Context, id int64) (*User, error)
	ListUsers(ctx context.Context) ([]User, error)
	UpdateUserProfile(ctx context.Context, userID int64, req UpdateProfileRequest) error
}

// Handler handles HTTP requests related to user operations.
type Handler struct {
	service ServiceInterface
}

// NewUserHandler creates a new Handler with the provided service.
func NewUserHandler(svc ServiceInterface) *Handler {
	return &Handler{service: svc}
}

// GetMe handles GET /user/me - returns the authenticated user's profile.
func (h *Handler) GetMe(c *gin.Context) {
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

// ListAll handles GET /user/ - lists all users.
func (h *Handler) ListAll(c *gin.Context) {
	users, err := h.service.ListUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// RegisterRoutes binds user-related endpoints to the router group.
func RegisterRoutes(router *gin.RouterGroup, h *Handler) {
	group := router.Group("/user")
	group.GET("/me", h.GetMe)
	group.GET("/", h.ListAll)
	group.PUT("/profile", h.UpdateProfile)
}

// ListUsers handles GET /user (alias for ListAll) - returns all users.
func (h *Handler) ListUsers(c *gin.Context) {
	users, err := h.service.ListUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// UpdateProfile handles PUT /user/profile - updates the authenticated user's profile.
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format", "details": err.Error()})
		return
	}

	if err := h.service.UpdateUserProfile(c, userID, req); err != nil {
		log.Printf("❌ UpdateUserProfile error for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated"})
}
