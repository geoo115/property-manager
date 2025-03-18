package accounting

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

// GetInvoices fetches all invoices with Redis caching.
func GetInvoices(c *gin.Context) {
	ctx := context.Background()
	cacheKey := "invoices"

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
	if err := db.DB.Preload("Tenant").Preload("Property").Find(&invoices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching invoices"})
		return
	}

	// Cache the result
	jsonData, _ := json.Marshal(invoices)
	db.RedisClient.Set(ctx, cacheKey, jsonData, 5*time.Minute)

	c.JSON(http.StatusOK, gin.H{"invoices": invoices, "cache": "miss"})
}
