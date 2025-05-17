package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nomenarkt/lamina/internal/middleware"
	"github.com/stretchr/testify/assert"
)

func generateExpiredToken(secret string) string {
	claims := middleware.CustomClaims{
		UserID: 3190,
		Email:  "expired@madagascarairlines.com",
		Role:   "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tok, _ := token.SignedString([]byte(secret))
	return tok
}

func TestMiddleware_ExpiredToken_ShouldReject(t *testing.T) {
	t.Setenv("JWT_SECRET", "expired_secret")

	r := gin.New()
	r.Use(middleware.JWTMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+generateExpiredToken("expired_secret"))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
