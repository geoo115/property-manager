package maintenance

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

// GetMaintenance fetches a maintenance request by ID with Redis caching.
func GetMaintenance(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()
	cacheKey := "maintenance:" + id

	// Try to get cached data from Redis
	cachedData, err := db.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var maintenance models.Maintenance
		if json.Unmarshal([]byte(cachedData), &maintenance) == nil {
			c.JSON(http.StatusOK, gin.H{"maintenance": maintenance, "cache": "hit"})
			return
		}
	}

	// Fetch from database
	var maintenance models.Maintenance
	if err := db.DB.Preload("RequestedBy").Preload("Property.Owner").
		First(&maintenance, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Maintenance not found"})
		return
	}

	// Store in Redis
	jsonData, _ := json.Marshal(maintenance)
	db.RedisClient.Set(ctx, cacheKey, jsonData, 10*time.Minute)

	c.JSON(http.StatusOK, gin.H{"maintenance": maintenance, "cache": "miss"})
}
