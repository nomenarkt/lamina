// Package utils provides utility helpers for handling Gin context and application constants.
package utils

import (
	"github.com/gin-gonic/gin"
)

// ContextUserIDKey is the key used to store the user ID in the Gin context.
const ContextUserIDKey = "userID"

// GetUserIDFromContext retrieves the user ID from the Gin context.
// It returns the ID and a boolean indicating whether it was successfully retrieved.
func GetUserIDFromContext(c *gin.Context) (int64, bool) {
	id, exists := c.Get(ContextUserIDKey)
	if !exists {
		return 0, false
	}
	userID, ok := id.(int64)
	return userID, ok
}
