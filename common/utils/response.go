package utils

import (
	"github.com/gin-gonic/gin"
)

const ContextUserIDKey = "userID"

func GetUserIDFromContext(c *gin.Context) (int64, bool) {
	id, exists := c.Get(ContextUserIDKey)
	if !exists {
		return 0, false
	}
	userID, ok := id.(int64)
	return userID, ok
}
