package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nomenarkt/lamina/common/utils"
)

func TestMiddleware_ValidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "testsecret123") // ✅ Safe test isolation

	// Step 1: Generate a valid token
	accessToken, _, err := utils.GenerateTokens(123, "admin", "admin@madagascarairlines.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Step 2: Setup a Gin router with the Middleware
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Middleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Access granted"})
	})

	// Step 3: Create a request with the token in Authorization header
	req, err := http.NewRequest(http.MethodGet, "/protected", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Step 4: Perform the request
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Step 5: Validate
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", recorder.Code)
	}
}
