// Package adminaccess handles admin-level access management like role and policy control.
package adminaccess

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nomenarkt/lamina/internal/access"
)

// AssignRoleRequest contains fields to assign a role to a user.
type AssignRoleRequest struct {
	UserID    int    `json:"user_id" binding:"required"`
	Function  string `json:"function" binding:"required"`    // e.g., "planner", "auditor"
	OrgUnitID int    `json:"org_unit_id" binding:"required"` // domain scope
}

// PolicyRequest defines the request payload for modifying access policies.
type PolicyRequest struct {
	Role      string `json:"role" binding:"required"`        // function/role name
	OrgUnitID int    `json:"org_unit_id" binding:"required"` // org domain
	Object    string `json:"object" binding:"required"`      // e.g., /api/flights
	Action    string `json:"action" binding:"required"`      // e.g., read/write
}

// AssignRole assigns a role to a user for a given domain.
func AssignRole(c *gin.Context) {
	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	e := access.GetEnforcer()
	sub := fmt.Sprintf("user:%d", req.UserID)
	dom := fmt.Sprintf("orgunit:%d", req.OrgUnitID)

	if _, err := e.AddGroupingPolicy(sub, req.Function, dom); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assigned"})
}

// RemoveRole removes a role from a user in a domain.
func RemoveRole(c *gin.Context) {
	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	e := access.GetEnforcer()
	sub := fmt.Sprintf("user:%d", req.UserID)
	dom := fmt.Sprintf("orgunit:%d", req.OrgUnitID)

	if _, err := e.RemoveGroupingPolicy(sub, req.Function, dom); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role removed"})
}

// AddPolicy adds a new access policy rule.
func AddPolicy(c *gin.Context) {
	var req PolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	e := access.GetEnforcer()
	dom := fmt.Sprintf("orgunit:%d", req.OrgUnitID)

	if _, err := e.AddPolicy(req.Role, dom, req.Object, req.Action); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Policy added"})
}

// DeletePolicy deletes an existing access policy rule.
func DeletePolicy(c *gin.Context) {
	var req PolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	e := access.GetEnforcer()
	dom := fmt.Sprintf("orgunit:%d", req.OrgUnitID)

	if _, err := e.RemovePolicy(req.Role, dom, req.Object, req.Action); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Policy removed"})
}

// ListPolicies lists all defined access policy rules.
func ListPolicies(c *gin.Context) {
	e := access.GetEnforcer()
	policies, err := e.GetPolicy()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"policies": policies})
}
