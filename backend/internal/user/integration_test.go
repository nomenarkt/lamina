package user

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nomenarkt/lamina/internal/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	// ✅ required for PostgreSQL driver registration
	_ "github.com/lib/pq"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetMe(ctx context.Context, id int64) (*User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	u := args.Get(0).(User)
	return &u, args.Error(1)
}

func (m *MockService) ListUsers(ctx context.Context) ([]User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]User), args.Error(1)
}

func (m *MockService) UpdateUserProfile(ctx context.Context, userID int64, req UpdateProfileRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *MockService) CreateUser(_ context.Context, _ *User) error {
	return nil // No-op for integration test unless you're testing CreateUser
}

func TestIntegration_GetUserMe_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = os.Setenv("JWT_SECRET", "test-secret")

	// ✅ Generate a valid token
	token, err := middleware.GenerateJWT("test-secret", 3190, "john@madagascarairlines.com", "admin")
	assert.NoError(t, err)

	// ✅ Mock user and service
	mockSvc := new(MockService)
	expectedUser := User{ID: 3190, Email: "john@madagascarairlines.com"}
	mockSvc.On("GetMe", mock.Anything, int64(3190)).Return(expectedUser, nil)

	// ✅ Setup Gin + handler + middleware
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.JWTMiddleware())
	handler := NewUserHandler(mockSvc)
	RegisterRoutes(router.Group(""), handler)

	// ✅ Perform HTTP request
	req, _ := http.NewRequest("GET", "/user/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// ✅ Assert response
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "john@madagascarairlines.com")
}
