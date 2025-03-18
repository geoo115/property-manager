package maintenance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

// GetLandlordMaintenances fetches maintenance requests for a landlord's property with Redis caching.
func GetLandlordMaintenances(c *gin.Context) {
	propertyIDStr := c.Param("id")
	propertyID, err := strconv.ParseUint(propertyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property ID"})
		return
	}

	userID, _ := c.Get("user_id")
	ctx := context.Background()
	cacheKey := fmt.Sprintf("maintenances:landlord:%d:property:%s", userID.(uint), propertyIDStr)

	// Try to get cached data from Redis
	cachedData, err := db.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var maintenances []models.Maintenance
		if json.Unmarshal([]byte(cachedData), &maintenances) == nil {
			c.JSON(http.StatusOK, gin.H{"maintenances": maintenances, "cache": "hit"})
			return
		}
	}

	// Verify landlord ownership
	var property models.Property
	if err := db.DB.Where("id = ? AND owner_id = ?", propertyID, userID).First(&property).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found or you do not own this property"})
		return
	}

	var maintenances []models.Maintenance
	if err := db.DB.Where("property_id = ?", propertyID).
		Preload("Reporter").Preload("Property.Owner").
		Find(&maintenances).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch maintenances"})
		return
	}

	// Store in Redis
	jsonData, _ := json.Marshal(maintenances)
	db.RedisClient.Set(ctx, cacheKey, jsonData, 10*time.Minute)

	c.JSON(http.StatusOK, gin.H{"maintenances": maintenances, "cache": "miss"})
}
