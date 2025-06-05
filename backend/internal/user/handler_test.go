package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockService implements ServiceInterface for unit testing.
type mockService struct {
	mock.Mock
}

func (m *mockService) GetMe(ctx context.Context, id int64) (*User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	user := args.Get(0).(User)
	return &user, args.Error(1)
}

func (m *mockService) ListUsers(ctx context.Context) ([]User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]User), args.Error(1)
}

func (m *mockService) UpdateUserProfile(ctx context.Context, userID int64, req UpdateProfileRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *mockService) CreateUser(_ context.Context, _ *User) error {
	return nil
}

// mock middleware: simulate setting userID in context
func setUserIDMiddleware(userID int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	}
}

func TestHandler_GetMe(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := new(mockService)
	h := NewUserHandler(ms)

	router := gin.New()
	router.Use(setUserIDMiddleware(100))
	router.GET("/user/me", h.GetMe)

	expected := User{ID: 100, Email: "me@example.com"}
	ms.On("GetMe", mock.Anything, int64(100)).Return(expected, nil)

	req := httptest.NewRequest("GET", "/user/me", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "me@example.com")
}

func TestHandler_GetMe_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := new(mockService)
	h := NewUserHandler(ms)

	router := gin.New()
	router.GET("/user/me", h.GetMe)

	req := httptest.NewRequest("GET", "/user/me", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "Unauthorized")
}

func TestHandler_ListAll(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := new(mockService)
	h := NewUserHandler(ms)

	router := gin.New()
	router.GET("/user", h.ListAll)

	mockUsers := []User{
		{ID: 1, Email: "a@example.com"},
		{ID: 2, Email: "b@example.com"},
	}
	ms.On("ListUsers", mock.Anything).Return(mockUsers, nil)

	req := httptest.NewRequest("GET", "/user", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "a@example.com")
	assert.Contains(t, resp.Body.String(), "b@example.com")
}

func TestHandler_ListUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := new(mockService)
	h := NewUserHandler(ms)

	router := gin.New()
	router.GET("/user/list-alias", h.ListUsers)

	mockUsers := []User{
		{ID: 3, Email: "alias@example.com"},
	}
	ms.On("ListUsers", mock.Anything).Return(mockUsers, nil)

	req := httptest.NewRequest("GET", "/user/list-alias", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "alias@example.com")
}

func TestHandler_UpdateProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := new(mockService)
	h := NewUserHandler(ms)

	router := gin.New()
	router.Use(setUserIDMiddleware(101))
	router.PUT("/user/profile", h.UpdateProfile)

	body := UpdateProfileRequest{
		FullName: "Updated User",
	}
	payload, _ := json.Marshal(body)

	ms.On("UpdateUserProfile", mock.Anything, int64(101), body).Return(nil)

	req := httptest.NewRequest("PUT", "/user/profile", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "profile updated")
}

func TestHandler_UpdateProfile_BadJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := new(mockService)
	h := NewUserHandler(ms)

	router := gin.New()
	router.Use(setUserIDMiddleware(101))
	router.PUT("/user/profile", h.UpdateProfile)

	req := httptest.NewRequest("PUT", "/user/profile", bytes.NewBufferString("{not-json"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "invalid format")
}

func TestHandler_UpdateProfile_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := new(mockService)
	h := NewUserHandler(ms)

	router := gin.New()
	router.Use(setUserIDMiddleware(101))
	router.PUT("/user/profile", h.UpdateProfile)

	body := UpdateProfileRequest{FullName: "Broken"}
	payload, _ := json.Marshal(body)

	ms.On("UpdateUserProfile", mock.Anything, int64(101), body).Return(errors.New("fail"))

	req := httptest.NewRequest("PUT", "/user/profile", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Contains(t, resp.Body.String(), "fail")
}

func TestRegisterRoutes_BindsCorrectly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := new(mockService)
	h := NewUserHandler(ms)

	// 🛠️ Stub ListUsers so it doesn't panic
	ms.On("ListUsers", mock.Anything).Return([]User{}, nil)

	router := gin.New()
	api := router.Group("/api")
	RegisterRoutes(api, h)

	// Basic check that /user is bound
	req := httptest.NewRequest("GET", "/api/user/", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}
