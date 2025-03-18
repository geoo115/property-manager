package accounting

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetInvoiceByID fetches an invoice by ID with Redis caching.
func GetInvoiceByID(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()
	cacheKey := "invoice:" + id

	// Check Redis cache
	cachedData, err := db.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var invoice models.Invoice
		if err := json.Unmarshal([]byte(cachedData), &invoice); err == nil {
			c.JSON(http.StatusOK, gin.H{"invoice": invoice, "cache": "hit"})
			return
		}
	}

	// Fetch from database on cache miss
	var invoice models.Invoice
	if err := db.DB.Preload("Tenant").Preload("Property").First(&invoice, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching invoice"})
		}
		return
	}

	// Cache the result
	jsonData, _ := json.Marshal(invoice)
	db.RedisClient.Set(ctx, cacheKey, jsonData, 10*time.Minute)

	c.JSON(http.StatusOK, gin.H{"invoice": invoice, "cache": "miss"})
}
