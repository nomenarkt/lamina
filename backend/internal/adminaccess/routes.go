// Package adminaccess handles admin-level access management like role and policy control.
package adminaccess

import (
	"github.com/gin-gonic/gin"
	"github.com/nomenarkt/lamina/internal/middleware"
)

// RegisterRoutes sets up admin access routes under the given router group.
func RegisterRoutes(rg *gin.RouterGroup) {
	admin := rg.Group("/admin")
	admin.Use(middleware.RequireRoles("admin"))

	admin.POST("/roles", AssignRole)
	admin.DELETE("/roles", RemoveRole)

	admin.POST("/policies", AddPolicy)
	admin.DELETE("/policies", DeletePolicy)
	admin.GET("/policies", ListPolicies)
}
