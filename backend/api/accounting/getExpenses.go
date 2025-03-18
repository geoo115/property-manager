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

// GetExpenses fetches all expenses with Redis caching.
func GetExpenses(c *gin.Context) {
	ctx := context.Background()
	cacheKey := "expenses"

	// Check Redis cache
	cachedData, err := db.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var expenses []models.Expense
		if err := json.Unmarshal([]byte(cachedData), &expenses); err == nil {
			c.JSON(http.StatusOK, gin.H{"expenses": expenses, "cache": "hit"})
			return
		}
	}

	// Fetch from database on cache miss
	var expenses []models.Expense
	if err := db.DB.Preload("Property.Owner").Find(&expenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching expenses"})
		return
	}

	// Cache the result
	jsonData, _ := json.Marshal(expenses)
	db.RedisClient.Set(ctx, cacheKey, jsonData, 5*time.Minute)

	c.JSON(http.StatusOK, gin.H{"expenses": expenses, "cache": "miss"})
}
