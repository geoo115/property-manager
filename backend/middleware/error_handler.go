package middleware

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/geoo115/property-manager/logger"
	"github.com/geoo115/property-manager/response"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ErrorHandler provides centralized error handling
func ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			logger.LogError(nil, "Panic recovered", logrus.Fields{
				"error":  err,
				"path":   c.Request.URL.Path,
				"method": c.Request.Method,
				"ip":     c.ClientIP(),
				"stack":  string(debug.Stack()),
			})
			response.InternalServerError(c, "Internal server error", nil)
		} else {
			logger.LogError(nil, "Panic recovered", logrus.Fields{
				"error":  recovered,
				"path":   c.Request.URL.Path,
				"method": c.Request.Method,
				"ip":     c.ClientIP(),
				"stack":  string(debug.Stack()),
			})
			response.InternalServerError(c, "Internal server error", nil)
		}
	})
}

// ValidationErrorHandler handles validation errors
func ValidationErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any validation errors
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				if err.Type == gin.ErrorTypeBind {
					logger.LogWarning("Validation error", logrus.Fields{
						"error":  err.Error(),
						"path":   c.Request.URL.Path,
						"method": c.Request.Method,
						"ip":     c.ClientIP(),
					})
					response.BadRequest(c, "Validation failed", err.Err)
					return
				}
			}
		}
	}
}

// NotFoundHandler handles 404 errors
func NotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.LogWarning("Route not found", logrus.Fields{
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
			"ip":     c.ClientIP(),
		})

		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Endpoint not found",
			"error":   "The requested resource does not exist",
			"path":    c.Request.URL.Path,
		})
	}
}

// MethodNotAllowedHandler handles 405 errors
func MethodNotAllowedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.LogWarning("Method not allowed", logrus.Fields{
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
			"ip":     c.ClientIP(),
		})

		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"success": false,
			"message": "Method not allowed",
			"error":   "The HTTP method is not supported for this endpoint",
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
		})
	}
}

// RateLimitExceededHandler handles rate limit errors
func RateLimitExceededHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.LogWarning("Rate limit exceeded", logrus.Fields{
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
			"ip":     c.ClientIP(),
		})

		c.JSON(http.StatusTooManyRequests, gin.H{
			"success":     false,
			"message":     "Rate limit exceeded",
			"error":       "Too many requests. Please try again later.",
			"retry_after": 60, // seconds
		})
	}
}

// DatabaseErrorHandler handles database errors
func DatabaseErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check for database errors in the response
		if c.Writer.Status() >= 500 {
			logger.LogError(nil, "Database error", logrus.Fields{
				"path":   c.Request.URL.Path,
				"method": c.Request.Method,
				"ip":     c.ClientIP(),
				"status": c.Writer.Status(),
			})
		}
	}
}

// SecurityErrorHandler handles security-related errors
func SecurityErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Log security-related errors
		if c.Writer.Status() == 401 || c.Writer.Status() == 403 {
			logger.LogWarning("Security error", logrus.Fields{
				"path":       c.Request.URL.Path,
				"method":     c.Request.Method,
				"ip":         c.ClientIP(),
				"status":     c.Writer.Status(),
				"user_agent": c.Request.UserAgent(),
			})
		}
	}
}

// RequestLogger logs all requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Log request details
		duration := time.Since(start)

		fields := logrus.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"duration":   duration,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}

		// Add user info if available
		if userID, exists := c.Get("user_id"); exists {
			fields["user_id"] = userID
		}
		if username, exists := c.Get("username"); exists {
			fields["username"] = username
		}

		if c.Writer.Status() >= 400 {
			logger.LogWarning("Request completed with error", fields)
		} else {
			logger.LogInfo("Request completed", fields)
		}
	}
}
