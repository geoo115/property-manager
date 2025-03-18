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

// GetLeases fetches all leases with Redis caching.
func GetLeases(c *gin.Context) {
	ctx := context.Background()
	cacheKey := "leases"

	// Try to get cached data from Redis
	cachedData, err := db.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var leases []models.Lease
		if json.Unmarshal([]byte(cachedData), &leases) == nil {
			c.JSON(http.StatusOK, gin.H{"leases": leases, "cache": "hit"})
			return
		}
	}

	var leases []models.Lease
	if err := db.DB.Preload("Tenant").Preload("Property.Owner").Find(&leases).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching leases"})
		return
	}

	// Store in Redis
	jsonData, _ := json.Marshal(leases)
	db.RedisClient.Set(ctx, cacheKey, jsonData, 10*time.Minute)

	c.JSON(http.StatusOK, gin.H{"leases": leases, "cache": "miss"})
}
