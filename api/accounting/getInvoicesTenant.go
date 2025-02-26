package accounting

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

// GetInvoicesForTenant fetches invoices for a tenant with Redis caching.
func GetInvoicesForTenant(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	ctx := context.Background()
	cacheKey := "tenant_invoices:" + strconv.Itoa(userID.(int))

	// Check Redis cache
	cachedData, err := db.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var invoices []models.Invoice
		if err := json.Unmarshal([]byte(cachedData), &invoices); err == nil {
			c.JSON(http.StatusOK, gin.H{"invoices": invoices, "cache": "hit"})
			return
		}
	}

	// Fetch from database on cache miss
	var invoices []models.Invoice
	if err := db.DB.Preload("Property").
		Where("tenant_id = ?", userID).
		Find(&invoices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching invoices for tenant"})
		return
	}

	// Cache the result
	jsonData, _ := json.Marshal(invoices)
	db.RedisClient.Set(ctx, cacheKey, jsonData, 10*time.Minute)

	c.JSON(http.StatusOK, gin.H{"invoices": invoices, "cache": "miss"})
}
