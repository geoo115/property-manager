package accounting

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

func GetInvoicesForLandlord(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	ctx := context.Background()
	// Fix type assertion from int to uint
	cacheKey := "landlord_invoices:" + strconv.Itoa(int(userID.(uint))) // Changed userID.(int) to userID.(uint)

	// Try to get cached data from Redis
	cachedData, err := db.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var invoices []models.Invoice
		if json.Unmarshal([]byte(cachedData), &invoices) == nil {
			c.JSON(http.StatusOK, gin.H{"invoices": invoices, "cache": "hit"})
			return
		}
	}

	// Fetch from database if cache miss
	var invoices []models.Invoice
	if err := db.DB.Joins("JOIN properties ON properties.id = invoices.property_id").
		Where("properties.owner_id = ?", userID).
		Preload("Tenant").
		Preload("Property.Owner").
		Find(&invoices).Error; err != nil {
		log.Printf("Error fetching invoices for landlord: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching invoices for landlord"})
		return
	}

	// Store in Redis
	jsonData, _ := json.Marshal(invoices)
	db.RedisClient.Set(ctx, cacheKey, jsonData, 10*time.Minute)

	c.JSON(http.StatusOK, gin.H{"invoices": invoices, "cache": "miss"})
}
