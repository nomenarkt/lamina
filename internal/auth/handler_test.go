package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nomenarkt/lamina/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Signup(c *gin.Context) {
	m.Called(c)
}

func (m *MockAuthService) SignupUser(ctx context.Context, req auth.SignupRequest) (auth.AuthResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(auth.AuthResponse), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, req auth.LoginRequest) (auth.AuthResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(auth.AuthResponse), args.Error(1)
}

func setupRouterWithMock(service auth.AuthServiceInterface) *gin.Engine {
	os.Setenv("JWT_SECRET", "testsecret123")
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1 := router.Group("/api/v1")
	auth.RegisterRoutes(v1, nil, service)
	return router
}

func TestSignup_CallsService(t *testing.T) {
	mockService := new(MockAuthService)
	router := setupRouterWithMock(mockService)

	mockService.On("Signup", mock.Anything).Return()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewBuffer([]byte(`{
		"email": "test@madagascarairlines.com",
		"password": "pass1234"
	}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	mockService.AssertCalled(t, "Signup", mock.Anything)
}

func TestLogin_Success(t *testing.T) {
	mockService := new(MockAuthService)
	router := setupRouterWithMock(mockService)

	loginReq := auth.LoginRequest{
		Email:    "test@madagascarairlines.com",
		Password: "pass1234",
	}
	loginRes := auth.AuthResponse{
		AccessToken:  "access123",
		RefreshToken: "refresh123",
	}
	mockService.On("Login", mock.Anything, loginReq).Return(loginRes, nil)

	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res auth.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, "access123", res.AccessToken)
	assert.Equal(t, "refresh123", res.RefreshToken)
}

func TestLogin_Failure(t *testing.T) {
	mockService := new(MockAuthService)
	router := setupRouterWithMock(mockService)

	reqData := auth.LoginRequest{
		Email:    "wrong@madagascarairlines.com",
		Password: "invalid",
	}
	mockService.On("Login", mock.Anything, reqData).Return(auth.AuthResponse{}, errors.New("invalid email or password"))

	body, _ := json.Marshal(reqData)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid email or password")
}

func TestLogin_BadJSON(t *testing.T) {
	mockService := new(MockAuthService)
	router := setupRouterWithMock(mockService)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(`{notvalid}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid character")
}
