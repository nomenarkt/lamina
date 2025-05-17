package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"github.com/nomenarkt/lamina/internal/middleware"
)

func generateToken(secret string, userID int64, role string) string {
	claims := middleware.CustomClaims{
		UserID: userID,
		Email:  "test@madagascarairlines.com",
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func TestJWTMiddleware_Success(t *testing.T) {
	t.Setenv("JWT_SECRET", "mytestsecret")

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.JWTMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	token := generateToken("mytestsecret", 3190, "admin")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJWTMiddleware_MissingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.JWTMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTMiddleware_InvalidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "mytestsecret")

	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.JWTMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRequireRoles_Forbidden(t *testing.T) {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("role", "viewer")
	}, middleware.RequireRoles("admin"))

	router.GET("/admin", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/admin", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
