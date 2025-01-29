package middleware

import (
	"DeliFood/backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware ensures only users with allowed roles can access a route
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No token provided"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid or expired token"})
			c.Abort()
			return
		}

		// Check if user's role is in the allowed roles
		for _, role := range allowedRoles {
			if claims.Role == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Insufficient permissions"})
		c.Abort()
	}
}
