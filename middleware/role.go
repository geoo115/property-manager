package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware enhanced with better admin handling
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User role not found in context"})
			return
		}

		// Allow admins to access any route
		if userRole == "admin" {
			c.Next()
			return
		}

		// Validate required role
		if userRole != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": fmt.Sprintf("Access denied. Required role: %s", requiredRole),
			})
			return
		}

		c.Next()
	}
}
