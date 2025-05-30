package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nomenarkt/lamina/internal/admin"
	"github.com/nomenarkt/lamina/internal/middleware"
	testutils "github.com/nomenarkt/lamina/internal/tests/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouterWithService(service *admin.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	v1 := r.Group("/api/v1")

	adminGroup := v1.Group("/admin", middleware.JWTMiddleware(), middleware.RequireRoles("admin"))
	adminGroup.POST("/create-user", func(c *gin.Context) {
		claimsVal, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authentication context"})
			return
		}

		claims, ok := claimsVal.(*middleware.CustomClaims)
		if !ok || claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user claims"})
			return
		}

		var req admin.CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input format"})
			return
		}

		err := service.InviteUser(c.Request.Context(), req, claims.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to invite user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "user successfully invited"})
	})

	return r
}

func TestCreateUser_Unauthorized(t *testing.T) {
	_ = os.Setenv("JWT_SECRET", "mytestsecret")

	service := admin.NewAdminService(new(testutils.MockAdminRepo), new(testutils.MockHasher))
	router := setupRouterWithService(service)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/create-user", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateUser_ForbiddenForViewer(t *testing.T) {
	_ = os.Setenv("JWT_SECRET", "mytestsecret")

	service := admin.NewAdminService(new(testutils.MockAdminRepo), new(testutils.MockHasher))
	router := setupRouterWithService(service)

	token, _ := middleware.GenerateJWT("mytestsecret", 1234, "viewer@madagascarairlines.com", "viewer")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/create-user", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCreateUser_Success_WithAdminRole(t *testing.T) {
	_ = os.Setenv("JWT_SECRET", "mytestsecret")

	mockRepo := new(testutils.MockAdminRepo)
	mockHasher := new(testutils.MockHasher)

	reqEmail := "successcase@madagascarairlines.com"
	tokenUserID := int64(99)

	mockRepo.On("IsEmailExists", reqEmail).Return(false, nil)
	mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(nil)
	mockRepo.On("FindUserIDByEmail", mock.Anything, reqEmail).Return(tokenUserID, nil)
	mockRepo.On("SetConfirmationToken", mock.Anything, tokenUserID, mock.AnythingOfType("string")).Return(nil)

	service := admin.NewAdminService(mockRepo, mockHasher)
	router := setupRouterWithService(service)

	token, _ := middleware.GenerateJWT("mytestsecret", 1, "admin@madagascarairlines.com", "admin")

	payload := `{
		"email": "successcase@madagascarairlines.com",
		"role": "editor"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/create-user", strings.NewReader(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	t.Logf("RESPONSE BODY: %s", w.Body.String())
	assert.Equal(t, http.StatusCreated, w.Code)
	mockRepo.AssertExpectations(t)
}
