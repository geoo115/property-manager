package accounting

import (
	"context"
	"fmt"
	"net/http"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DeleteExpense deletes an expense and invalidates Redis cache.
func DeleteExpense(c *gin.Context) {
	id := c.Param("id")
	var expense models.Expense

	if err := db.DB.First(&expense, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching expense"})
		}
		return
	}

	if err := db.DB.Delete(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting expense"})
		return
	}

	// Invalidate Redis caches
	ctx := context.Background()
	cacheKeys := []string{"expenses", fmt.Sprintf("expense:%s", id)}
	for _, key := range cacheKeys {
		db.RedisClient.Del(ctx, key)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}
