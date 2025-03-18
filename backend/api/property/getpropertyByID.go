package property

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

// GetPropertyByID fetches a property by ID with Redis caching.
func GetPropertyByID(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()
	cacheKey := "property:" + id

	// Try to get cached data from Redis
	cachedData, err := db.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var property models.Property
		if json.Unmarshal([]byte(cachedData), &property) == nil {
			c.JSON(http.StatusOK, gin.H{"property": property, "cache": "hit"})
			return
		}
	}

	// Fetch from database if cache miss
	var property models.Property
	if err := db.DB.Preload("Units").Preload("Owner").
		First(&property, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}

	// Store in Redis
	jsonData, _ := json.Marshal(property)
	db.RedisClient.Set(ctx, cacheKey, jsonData, 10*time.Minute)

	c.JSON(http.StatusOK, gin.H{"property": property, "cache": "miss"})
}
