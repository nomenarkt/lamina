package tests

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

// mockMiddleware injects a fixed userID into the Gin context to simulate an authenticated request
func mockMiddleware(userID int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	}
}

func setupRouter(handler *user.Handler, withAuth bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	group := r.Group("/api/v1/user")
	if withAuth {
		group.Use(mockMiddleware(1))
	}

	group.GET("/me", handler.GetMe)
	group.GET("/", handler.ListUsers)

	return r
}

func TestGetMe_Unauthorized(t *testing.T) {
	svc := &MockUserService{}
	handler := user.NewUserHandler(svc)
	r := setupRouter(handler, false)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized")
}

func TestGetMe_Success(t *testing.T) {
	svc := &MockUserService{
		GetMeFunc: func(_ context.Context, id int64) (*user.User, error) {
			return &user.User{ID: id, Email: "john@example.com"}, nil
		},
	}
	handler := user.NewUserHandler(svc)
	r := setupRouter(handler, true)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "john@example.com")
}

func TestListUsers_Success(t *testing.T) {
	svc := &MockUserService{
		ListUsersFunc: func(_ context.Context) ([]user.User, error) {
			return []user.User{
				{ID: 1, Email: "user1@madagascarairlines.com"},
				{ID: 2, Email: "user2@madagascarairlines.com"},
			}, nil
		},
	}
	handler := user.NewUserHandler(svc)
	r := setupRouter(handler, true)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "user1@madagascarairlines.com")
}

func TestListUsers_Failure(t *testing.T) {
	svc := &MockUserService{
		ListUsersFunc: func(_ context.Context) ([]user.User, error) {
			return nil, errors.New("DB failure")
		},
	}
	handler := user.NewUserHandler(svc)
	r := setupRouter(handler, true)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "DB failure")
}
