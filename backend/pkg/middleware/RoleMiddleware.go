package middleware

import (
	"DeliFood/backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RoleMiddleware ensures only users with allowed roles can access a route
func RoleMiddleware(requiredRole models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists || models.Role(userRole.(string)) != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}
		c.Next()
	}
}
