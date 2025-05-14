package auth

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAndParseToken(t *testing.T) {
	if err := os.Setenv("JWT_SECRET", "testsecret"); err != nil {
		t.Fatalf("failed to set JWT_SECRET: %v", err)
	}

	tokenStr, _, err := GenerateTokens(123, "admin", "admin@madagascarairlines.com")
	assert.NoError(t, err)

	parsedToken, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(_ *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(*Claims)
	assert.True(t, ok)
	assert.Equal(t, int64(123), claims.UserID)
	assert.Equal(t, "admin", claims.Role)
	assert.Equal(t, "admin@madagascarairlines.com", claims.Email)
	assert.WithinDuration(t, time.Now(), claims.IssuedAt.Time, time.Minute)
}
