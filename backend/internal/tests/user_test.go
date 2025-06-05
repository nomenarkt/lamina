package tests

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
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
	CreateUserFunc        func(ctx context.Context, u *user.User) error
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

func (m *MockUserService) CreateUser(ctx context.Context, u *user.User) error {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, u)
	}
	return nil
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

func TestCreateUser_Handler(t *testing.T) {
	mockService := &MockUserService{
		CreateUserFunc: func(_ context.Context, u *user.User) error {
			if u.Email == "duplicate@madagascarairlines.com" {
				return errors.New("email already in use")
			}
			u.ID = 42 // simulate DB assigning ID
			return nil
		},
	}

	handler := user.NewUserHandler(mockService)
	r := setupRouter(handler, false) // No auth required for user creation
	r.POST("/api/v1/user", handler.CreateUser)

	t.Run("success", func(t *testing.T) {
		body := `{
			"email": "new@madagascarairlines.com",
			"password_hash": "hashedpass",
			"role": "user",
			"user_type": "external",
			"status": "pending"
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/user", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), `"id":42`)
	})

	t.Run("duplicate email", func(t *testing.T) {
		body := `{
			"email": "duplicate@madagascarairlines.com",
			"password_hash": "hashedpass",
			"role": "user",
			"user_type": "external",
			"status": "pending"
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/user", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "email already in use")
	})

	t.Run("malformed JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/user", strings.NewReader(`bad json`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid JSON")
	})
}
