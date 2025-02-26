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

// GetUsers fetches all users with Redis caching.
func GetUsers(c *gin.Context) {
	ctx := context.Background()
	cacheKey := "users"

	// Try to get cached data from Redis
	cachedData, err := db.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var users []models.User
		if json.Unmarshal([]byte(cachedData), &users) == nil {
			c.JSON(http.StatusOK, gin.H{"users": users, "cache": "hit"})
			return
		}
	}

	var users []models.User
	if err := db.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users"})
		return
	}

	// Store in Redis
	jsonData, _ := json.Marshal(users)
	db.RedisClient.Set(ctx, cacheKey, jsonData, 10*time.Minute)

	c.JSON(http.StatusOK, gin.H{"users": users, "cache": "miss"})
}

func GetActiveUsers(c *gin.Context) {
	var users []models.User
	err := db.DB.Where("role = ?", "tenant").Find(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}
