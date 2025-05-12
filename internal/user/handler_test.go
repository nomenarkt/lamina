package user_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nomenarkt/lamina/internal/user"
	"github.com/stretchr/testify/assert"
)

// MockUserService uses function fields to simulate service methods
type MockUserService struct {
	GetMeFunc             func(ctx context.Context, id int64) (*user.User, error)
	GetProfileFunc        func(ctx context.Context, id int64) (*user.User, error)
	ListUsersFunc         func(ctx context.Context) ([]user.User, error)
	UpdateUserProfileFunc func(ctx context.Context, userID int64, req user.UpdateProfileRequest) error
}

func (m *MockUserService) GetMe(ctx context.Context, id int64) (*user.User, error) {
	return m.GetMeFunc(ctx, id)
}
func (m *MockUserService) GetProfile(ctx context.Context, id int64) (*user.User, error) {
	return m.GetProfileFunc(ctx, id)
}
func (m *MockUserService) ListUsers(ctx context.Context) ([]user.User, error) {
	return m.ListUsersFunc(ctx)
}

func (m *MockUserService) UpdateUserProfile(ctx context.Context, userID int64, req user.UpdateProfileRequest) error {
	return m.UpdateUserProfileFunc(ctx, userID, req)
}

// Injects userID into Gin context properly
func injectUserID(c *gin.Context, id string) {
	c.Set("userID", id)
}

func setupRouter(handler *user.UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	group := r.Group("/api/v1/user")
	group.GET("/me", handler.GetMe)
	group.GET("/", handler.ListUsers)

	return r
}

func TestGetMe_Unauthorized(t *testing.T) {
	svc := &MockUserService{}
	handler := user.NewUserHandler(svc)
	r := setupRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized")
}

func TestGetMe_Success(t *testing.T) {
	svc := &MockUserService{
		GetMeFunc: func(ctx context.Context, id int64) (*user.User, error) {
			return &user.User{ID: id, Email: "john@example.com"}, nil
		},
	}
	handler := user.NewUserHandler(svc)
	r := gin.New()
	r.GET("/api/v1/user/me", func(c *gin.Context) {
		c.Set("userID", int64(1)) // ✅ This makes GetUserIDFromContext return correctly
		handler.GetMe(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "john@example.com")
}

func TestListUsers_Success(t *testing.T) {
	svc := &MockUserService{
		ListUsersFunc: func(ctx context.Context) ([]user.User, error) {
			return []user.User{
				{ID: 1, Email: "user1@madagascarairlines.com"},
				{ID: 2, Email: "user2@madagascarairlines.com"},
			}, nil
		},
	}
	handler := user.NewUserHandler(svc)

	r := gin.New()
	r.GET("/api/v1/user/", func(c *gin.Context) {
		c.Set("userID", int64(1))
		handler.ListUsers(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "user1@madagascarairlines.com")
}

func TestListUsers_Failure(t *testing.T) {
	svc := &MockUserService{
		ListUsersFunc: func(ctx context.Context) ([]user.User, error) {
			return nil, errors.New("DB failure")
		},
	}
	handler := user.NewUserHandler(svc)

	r := gin.New()
	r.GET("/api/v1/user/", func(c *gin.Context) {
		c.Set("userID", int64(1))
		handler.ListUsers(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "DB failure")
}
