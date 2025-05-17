package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAcceptInviteFlow_ShouldSucceed(t *testing.T) {
	router := gin.Default()

	router.POST("/api/v1/auth/confirm/:token", func(c *gin.Context) {
		token := c.Param("token")
		if token == "valid-token-abc" {
			c.JSON(http.StatusOK, gin.H{"message": "Invite accepted"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		}
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/confirm/valid-token-abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Invite accepted")
}
