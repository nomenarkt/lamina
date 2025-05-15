package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAndParseToken(t *testing.T) {
	secret := "testsecret"
	t.Setenv("JWT_SECRET", secret)

	accessToken, _, err := GenerateTokens("testsecret", "testrefresh", 123, "admin@madagascarairlines.com", "admin")
	assert.NoError(t, err)

	parsedToken, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return []byte(secret), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(*Claims)
	assert.True(t, ok)
	assert.Equal(t, int64(123), claims.UserID)
	assert.Equal(t, "admin@madagascarairlines.com", claims.Email)
	assert.Equal(t, "admin", claims.Role)
	assert.WithinDuration(t, time.Now(), claims.IssuedAt.Time, time.Minute)
}
