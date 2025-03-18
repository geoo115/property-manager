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

func GetExpensesForLandlord(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	ctx := context.Background()
	cacheKey := "landlord_expenses:" + strconv.Itoa(int(userID.(uint))) // Fix here too

	cachedData, err := db.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var expenses []models.Expense
		if json.Unmarshal([]byte(cachedData), &expenses) == nil {
			c.JSON(http.StatusOK, gin.H{"expenses": expenses, "cache": "hit"})
			return
		}
	}

	var expenses []models.Expense
	if err := db.DB.Joins("JOIN properties ON properties.id = expenses.property_id").
		Where("properties.owner_id = ?", userID).
		Preload("Property.Owner").
		Find(&expenses).Error; err != nil {
		log.Printf("Error fetching expenses for landlord: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching expenses for landlord"})
		return
	}

	jsonData, _ := json.Marshal(expenses)
	db.RedisClient.Set(ctx, cacheKey, jsonData, 10*time.Minute)

	c.JSON(http.StatusOK, gin.H{"expenses": expenses, "cache": "miss"})
}
