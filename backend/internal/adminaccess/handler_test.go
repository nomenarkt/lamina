package adminaccess

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nomenarkt/lamina/internal/access"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRouterWithPolicy(t *testing.T) *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Setup in-memory Casbin enforcer
	enforcer := access.InitTestEnforcer(t)
	access.SetEnforcer(enforcer)

	// Setup permissions
	_, err := enforcer.AddPolicy("planner", "orgunit:1", "/api/crew", "assign")
	require.NoError(t, err)
	_, err = enforcer.AddGroupingPolicy("user:201", "planner", "orgunit:1")
	require.NoError(t, err)

	// Setup test router
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", "201")
		c.Set("role", "planner")
		c.Next()
	})

	// Simulated protected route
	r.POST("/api/crew/assign", func(c *gin.Context) {
		e := access.GetEnforcer()
		ok, _ := e.Enforce("user:201", "orgunit:1", "/api/crew", "assign")
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"message": "forbidden"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "crew assigned"})
	})

	return r
}

func TestPlanner_AssignCrew_Success(t *testing.T) {
	router := setupTestRouterWithPolicy(t)

	body := `{}`
	req, _ := http.NewRequest(http.MethodPost, "/api/crew/assign", bytes.NewBuffer([]byte(body)))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var data map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&data)
	assert.Equal(t, "crew assigned", data["message"])
}

func TestGetUserEffectivePoliciesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	enforcer := access.InitTestEnforcer(t)
	access.SetEnforcer(enforcer)

	// Setup test policy and user
	_, err := enforcer.AddPolicy("planner", "orgunit:1", "/api/crew", "assign")
	require.NoError(t, err)
	_, err = enforcer.AddGroupingPolicy("user:201", "planner", "orgunit:1")
	require.NoError(t, err)

	router := gin.New()
	router.GET("/user/:id/policies", func(c *gin.Context) {
		c.Set("userID", "201")
		c.Set("role", "planner")
		c.Next()
	}, GetUserEffectivePoliciesHandler)

	req := httptest.NewRequest("GET", "/user/201/policies?org_unit_id=1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "/api/crew")
}

func TestViewer_CannotAccessPolicies(t *testing.T) {
	gin.SetMode(gin.TestMode)

	enforcer := access.InitTestEnforcer(t)
	access.SetEnforcer(enforcer)

	// Setup viewer user
	_, _ = enforcer.AddGroupingPolicy("user:999", "viewer", "orgunit:1")

	router := gin.New()
	router.GET("/user/:id/policies", func(c *gin.Context) {
		c.Set("userID", "999")
		c.Set("role", "viewer")
		c.Next()
	}, GetUserEffectivePoliciesHandler)

	req := httptest.NewRequest("GET", "/user/999/policies?org_unit_id=1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
