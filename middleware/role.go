package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware ensures only authorized users can access certain routes
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")

		// Log debugging information
		fmt.Println("Middleware: Checking role for user ->", userRole)

		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User role not found in context"})
			return
		}

		// Allow admins to access any route
		if userRole == "admin" {
			fmt.Println("Admin Access Granted")
			c.Next()
			return
		}

		// If the required role does not match, deny access
		if userRole != requiredRole {
			fmt.Println("Access Denied. Required:", requiredRole, "User has:", userRole)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": fmt.Sprintf("Access denied. Required role: %s", requiredRole),
			})
			return
		}

		fmt.Println("Access Granted to:", userRole)
		c.Next()
	}
}
