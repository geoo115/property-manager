package lease

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

// GetLeaseByID fetches a lease by its ID with Redis caching.
func GetLeaseByID(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()
	cacheKey := "lease:" + id

	// Try to get cached data from Redis
	cachedData, err := db.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var lease models.Lease
		if json.Unmarshal([]byte(cachedData), &lease) == nil {
			c.JSON(http.StatusOK, gin.H{"lease": lease, "cache": "hit"})
			return
		}
	}

	// Fetch from database if cache miss
	var lease models.Lease
	if err := db.DB.Preload("Tenant").Preload("Property.Owner").
		First(&lease, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lease not found"})
		return
	}

	// Store in Redis
	jsonData, _ := json.Marshal(lease)
	db.RedisClient.Set(ctx, cacheKey, jsonData, 10*time.Minute)

	c.JSON(http.StatusOK, gin.H{"lease": lease, "cache": "miss"})
}
