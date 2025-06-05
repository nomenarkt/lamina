package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// ✅ Set test JWT secret used in both middleware and token generation
	_ = os.Setenv("JWT_SECRET", "test-secret-123")
	os.Exit(m.Run())
}

func TestJWTMiddleware_ValidToken(t *testing.T) {
	router := gin.New()
	token, err := GenerateJWT("test-secret-123", 123, "user@example.com", "admin")
	assert.NoError(t, err)

	router.Use(JWTMiddleware())
	router.GET("/secure", func(c *gin.Context) {
		userID := GetUserID(c)
		role := GetUserRole(c)
		c.JSON(200, gin.H{"userID": userID, "role": role})
	})

	req, _ := http.NewRequest("GET", "/secure", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "123")
	assert.Contains(t, resp.Body.String(), "admin")
}

func TestJWTMiddleware_MissingToken(t *testing.T) {
	router := gin.New()
	router.Use(JWTMiddleware())
	router.GET("/secure", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req, _ := http.NewRequest("GET", "/secure", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestJWTMiddleware_InvalidToken(t *testing.T) {
	router := gin.New()
	router.Use(JWTMiddleware())
	router.GET("/secure", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req, _ := http.NewRequest("GET", "/secure", nil)
	req.Header.Set("Authorization", "Bearer not-a-valid-token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestJWTMiddleware_ExpiredToken(t *testing.T) {
	// Create an already expired token
	claims := CustomClaims{
		UserID: 123,
		Email:  "expired@example.com",
		Role:   "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte("test-secret-123"))
	assert.NoError(t, err)

	router := gin.New()
	router.Use(JWTMiddleware())
	router.GET("/secure", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req, _ := http.NewRequest("GET", "/secure", nil)
	req.Header.Set("Authorization", "Bearer "+signed)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestRequireRoles_AcceptsAllowedRole(t *testing.T) {
	router := gin.New()
	token, _ := GenerateJWT("test-secret-123", 321, "admin@example.com", "admin")

	router.Use(JWTMiddleware(), RequireRoles("admin", "planner"))
	router.GET("/admin", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "welcome"})
	})

	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "welcome")
}

func TestRequireRoles_RejectsForbiddenRole(t *testing.T) {
	router := gin.New()
	token, _ := GenerateJWT("test-secret-123", 321, "viewer@example.com", "viewer")

	router.Use(JWTMiddleware(), RequireRoles("admin"))
	router.GET("/admin", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "should not get here"})
	})

	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusForbidden, resp.Code)
	assert.Contains(t, resp.Body.String(), "Forbidden: insufficient role")
}
