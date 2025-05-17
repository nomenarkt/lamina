package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nomenarkt/lamina/internal/admin"
	"github.com/nomenarkt/lamina/internal/middleware"
	"github.com/stretchr/testify/assert"
)

func setupRouterWithDeleteSelf(_ *admin.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	v1 := r.Group("/api/v1")

	adminGroup := v1.Group("/admin", middleware.JWTMiddleware(), middleware.RequireRoles("admin"))
	adminGroup.DELETE("/me", func(c *gin.Context) {
		claimsVal, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth context"})
			return
		}
		claims := claimsVal.(*middleware.CustomClaims)

		// Prevent self-deletion
		if claims.Email == "admin@madagascarairlines.com" {
			c.JSON(http.StatusForbidden, gin.H{"error": "cannot delete yourself"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
	})

	return r
}

func TestAdminSelfDelete_ShouldFail(t *testing.T) {
	_ = os.Setenv("JWT_SECRET", "selfsecret")

	service := admin.NewAdminService(nil, nil)
	router := setupRouterWithDeleteSelf(service)

	token, _ := middleware.GenerateJWT("selfsecret", 1, "admin@madagascarairlines.com", "admin")

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "cannot delete yourself")
}
