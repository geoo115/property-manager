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

// GetMaintenances fetches all maintenance requests with Redis caching.
func GetMaintenances(c *gin.Context) {
	ctx := context.Background()
	cacheKey := "maintenances"

	var maintenances []models.Maintenance
	userRole, _ := c.Get("user_role")
	userID, _ := c.Get("user_id")

	// Different cache keys per role
	if userRole == "tenant" {
		leaseIDStr := c.Param("id")
		if leaseIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Lease ID is required"})
			return
		}
		cacheKey = fmt.Sprintf("maintenances:tenant:%d:lease:%s", userID.(uint), leaseIDStr)
	} else if userRole == "maintenanceTeam" {
		cacheKey = "maintenances:team"
	} else if userRole == "admin" {
		cacheKey = "maintenances:all"
	}

	// Try to get cached data from Redis
	cachedData, err := db.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		if json.Unmarshal([]byte(cachedData), &maintenances) == nil {
			c.JSON(http.StatusOK, gin.H{"maintenances": maintenances, "cache": "hit"})
			return
		}
	}

	query := db.DB.Model(&models.Maintenance{})

	switch userRole {
	case "admin":
		query = query.Preload("Reporter").Preload("Property")
	case "tenant":
		leaseIDStr := c.Param("id")
		leaseID, err := strconv.ParseUint(leaseIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lease ID"})
			return
		}

		var lease models.Lease
		if err := db.DB.Where("id = ? AND tenant_id = ?", leaseID, userID).First(&lease).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Lease not found or access denied"})
			return
		}

		// Fetch all maintenance requests for the property tied to the lease
		query = query.Where("property_id = ?", lease.PropertyID).
			Preload("Property").Preload("Reporter")
	case "maintenanceTeam":
		query = query.Preload("Property").Preload("Reporter")
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	if err := query.Find(&maintenances).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching maintenance requests"})
		return
	}

	// Store in Redis
	jsonData, _ := json.Marshal(maintenances)
	db.RedisClient.Set(ctx, cacheKey, jsonData, 10*time.Minute)

	c.JSON(http.StatusOK, gin.H{"maintenances": maintenances, "cache": "miss"})
}
