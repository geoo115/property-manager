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

// GetExpenseByID fetches an expense by ID with Redis caching.
func GetExpenseByID(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()
	cacheKey := "expense:" + id

	// Check Redis cache
	cachedData, err := db.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var expense models.Expense
		if err := json.Unmarshal([]byte(cachedData), &expense); err == nil {
			c.JSON(http.StatusOK, gin.H{"expense": expense, "cache": "hit"})
			return
		}
	}

	// Fetch from database on cache miss
	var expense models.Expense
	if err := db.DB.Preload("Property.Owner").First(&expense, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching expense"})
		}
		return
	}

	// Cache the result
	jsonData, _ := json.Marshal(expense)
	db.RedisClient.Set(ctx, cacheKey, jsonData, 10*time.Minute)

	c.JSON(http.StatusOK, gin.H{"expense": expense, "cache": "miss"})
}
