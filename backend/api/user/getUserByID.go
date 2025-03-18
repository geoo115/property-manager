package user

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

// GetUserByID fetches a user by ID with Redis caching.
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()
	cacheKey := "user:" + id

	// Try to get cached data from Redis
	cachedData, err := db.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var user models.User
		if json.Unmarshal([]byte(cachedData), &user) == nil {
			c.JSON(http.StatusOK, gin.H{"user": user, "cache": "hit"})
			return
		}
	}

	// Fetch from database if cache miss
	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Store in Redis
	jsonData, _ := json.Marshal(user)
	db.RedisClient.Set(ctx, cacheKey, jsonData, 10*time.Minute)

	c.JSON(http.StatusOK, gin.H{"user": user, "cache": "miss"})
}
