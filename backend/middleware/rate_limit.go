package middleware

import (
	"fmt"
	"strconv"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/logger"
	"github.com/geoo115/property-manager/response"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// RateLimitConfig holds configuration for rate limiting
type RateLimitConfig struct {
	Requests int                       // Number of requests allowed
	Window   time.Duration             // Time window for the rate limit
	KeyFunc  func(*gin.Context) string // Function to generate the key for rate limiting
}

// DefaultRateLimitConfig returns a default rate limit configuration
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Requests: 100,
		Window:   time.Minute,
		KeyFunc:  func(c *gin.Context) string { return c.ClientIP() },
	}
}

// RateLimit creates a rate limiting middleware
func RateLimit(config RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if Redis is available
		if db.RedisClient == nil {
			// Continue without rate limiting if Redis is not available
			c.Next()
			return
		}

		key := fmt.Sprintf("rate_limit:%s", config.KeyFunc(c))

		// Get current count
		currentStr, err := db.RedisClient.Get(db.Ctx, key).Result()
		if err != nil && err.Error() != "redis: nil" {
			logger.LogError(err, "Failed to get rate limit count from Redis", logrus.Fields{
				"key": key,
				"ip":  c.ClientIP(),
			})
			// Continue without rate limiting if Redis is unavailable
			c.Next()
			return
		}

		current := 0
		if currentStr != "" {
			current, _ = strconv.Atoi(currentStr)
		}

		// Check if limit exceeded
		if current >= config.Requests {
			logger.LogWarning("Rate limit exceeded", logrus.Fields{
				"ip":       c.ClientIP(),
				"current":  current,
				"limit":    config.Requests,
				"window":   config.Window,
				"endpoint": c.Request.URL.Path,
			})

			// Get TTL for rate limit reset time
			ttl, err := db.RedisClient.TTL(db.Ctx, key).Result()
			if err == nil && ttl > 0 {
				c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(ttl).Unix()))
			}

			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", config.Requests))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("Retry-After", fmt.Sprintf("%d", int(config.Window.Seconds())))

			response.TooManyRequests(c, "Rate limit exceeded")
			c.Abort()
			return
		}

		// Increment counter
		pipe := db.RedisClient.Pipeline()
		pipe.Incr(db.Ctx, key)
		if current == 0 {
			// Set expiration only on first request
			pipe.Expire(db.Ctx, key, config.Window)
		}
		_, err = pipe.Exec(db.Ctx)
		if err != nil {
			logger.LogError(err, "Failed to update rate limit count in Redis", logrus.Fields{
				"key": key,
				"ip":  c.ClientIP(),
			})
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", config.Requests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", config.Requests-current-1))

		c.Next()
	}
}

// IPRateLimit creates an IP-based rate limiting middleware
func IPRateLimit(requests int, window time.Duration) gin.HandlerFunc {
	config := RateLimitConfig{
		Requests: requests,
		Window:   window,
		KeyFunc:  func(c *gin.Context) string { return c.ClientIP() },
	}
	return RateLimit(config)
}

// UserRateLimit creates a user-based rate limiting middleware
func UserRateLimit(requests int, window time.Duration) gin.HandlerFunc {
	config := RateLimitConfig{
		Requests: requests,
		Window:   window,
		KeyFunc: func(c *gin.Context) string {
			if userID, exists := c.Get("user_id"); exists {
				return fmt.Sprintf("user:%v", userID)
			}
			return c.ClientIP() // Fallback to IP if user not authenticated
		},
	}
	return RateLimit(config)
}

// EndpointRateLimit creates an endpoint-specific rate limiting middleware
func EndpointRateLimit(requests int, window time.Duration) gin.HandlerFunc {
	config := RateLimitConfig{
		Requests: requests,
		Window:   window,
		KeyFunc: func(c *gin.Context) string {
			return fmt.Sprintf("%s:%s", c.ClientIP(), c.Request.URL.Path)
		},
	}
	return RateLimit(config)
}
