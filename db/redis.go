package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"), // Example: "localhost:6379"
		Password: "",                      // No password by default
		DB:       0,                       // Use default DB
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	fmt.Println("✅ Connected to Redis")
}

// Redis-based rate limiting middleware
func RateLimit(limit int, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%s", ip)

		current, err := RedisClient.Incr(Ctx, key).Result()
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Internal error"})
			return
		}

		if current == 1 {
			RedisClient.Expire(Ctx, key, duration)
		}

		if current > int64(limit) {
			c.AbortWithStatusJSON(429, gin.H{"error": "Too many requests"})
			return
		}

		c.Next()
	}
}
