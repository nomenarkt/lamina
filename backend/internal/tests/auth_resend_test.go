package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Extend MockAuthService with ResendConfirmation
func (m *MockAuthService) ResendConfirmation(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func TestResendConfirmation_Success(t *testing.T) {
	mockService := new(MockAuthService)
	router := setupRouterWithMock(t, mockService)

	email := "internal@madagascarairlines.com"
	mockService.On("ResendConfirmation", mock.Anything, email).Return(nil)

	body := map[string]string{"email": email}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/resend-confirmation", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Confirmation email resent")
	mockService.AssertCalled(t, "ResendConfirmation", mock.Anything, email)
}

func TestResendConfirmation_AlreadyConfirmed(t *testing.T) {
	mockService := new(MockAuthService)
	router := setupRouterWithMock(t, mockService)

	email := "confirmed@madagascarairlines.com"
	mockService.On("ResendConfirmation", mock.Anything, email).
		Return(errors.New("user already confirmed or invalid status"))

	body := map[string]string{"email": email}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/resend-confirmation", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "User already confirmed")
}

func TestResendConfirmation_UserNotFound(t *testing.T) {
	mockService := new(MockAuthService)
	router := setupRouterWithMock(t, mockService)

	email := "notfound@madagascarairlines.com"
	mockService.On("ResendConfirmation", mock.Anything, email).
		Return(errors.New("user not found"))

	body := map[string]string{"email": email}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/resend-confirmation", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "User not found")
}

func TestResendConfirmation_ExternalUserBlocked(t *testing.T) {
	mockService := new(MockAuthService)
	router := setupRouterWithMock(t, mockService)

	email := "external@otherdomain.com"
	mockService.On("ResendConfirmation", mock.Anything, email).
		Return(errors.New("resend allowed only for internal users"))

	body := map[string]string{"email": email}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/resend-confirmation", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Resend allowed only for internal users")
}
