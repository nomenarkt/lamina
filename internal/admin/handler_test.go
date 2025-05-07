package admin_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/nomenarkt/lamina/internal/admin"
	"github.com/nomenarkt/lamina/internal/middleware"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	os.Setenv("JWT_SECRET", "mytestsecret")
	gin.SetMode(gin.TestMode)

	r := gin.New()
	v1 := r.Group("/api/v1")
	db := &sqlx.DB{} // dummy DB for signature compliance
	admin.RegisterRoutes(v1, db)
	return r
}

func TestCreateUser_Unauthorized(t *testing.T) {
	router := setupRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/create-user", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateUser_ForbiddenForViewer(t *testing.T) {
	os.Setenv("JWT_SECRET", "mytestsecret")

	router := setupRouter()
	token, _ := middleware.GenerateJWT("mytestsecret", 1234, "viewer@madagascarairlines.com", "viewer")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/create-user", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}
