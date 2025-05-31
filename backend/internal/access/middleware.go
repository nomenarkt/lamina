// Package access provides role-based access control (RBAC) using Casbin and PostgreSQL.
package access

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CasbinMiddleware returns a Gin middleware handler that uses Casbin for RBAC enforcement.
func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		enforcer := GetEnforcer()
		userID := c.MustGet("userID").(int)
		companyID := c.MustGet("companyID").(int)

		sub := fmt.Sprintf("user:%d", userID)
		dom := fmt.Sprintf("company:%d", companyID)
		obj := c.Request.URL.Path
		act := c.Request.Method

		ok, err := enforcer.Enforce(sub, dom, obj, act)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "enforcement error"})
			return
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
			return
		}

		c.Next()
	}
}
