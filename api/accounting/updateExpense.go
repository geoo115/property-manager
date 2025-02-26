package accounting

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UpdateExpense updates an expense and invalidates Redis cache.
func UpdateExpense(c *gin.Context) {
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

	var input struct {
		PropertyID  uint    `json:"property_id" binding:"required"`
		Description string  `json:"description" binding:"required"`
		Category    string  `json:"category" binding:"required"`
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		ExpenseDate string  `json:"expense_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense data", "details": err.Error()})
		return
	}

	// Parse date
	expenseDate, err := time.Parse("2006-01-02", input.ExpenseDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense date format, use YYYY-MM-DD"})
		return
	}

	// Update fields
	expense.PropertyID = input.PropertyID
	expense.Description = input.Description
	expense.Category = input.Category
	expense.Amount = input.Amount
	expense.ExpenseDate = expenseDate

	if err := db.DB.Save(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating expense"})
		return
	}

	// Preload associations
	if err := db.DB.Preload("Property").First(&expense, expense.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching updated expense"})
		return
	}

	// Invalidate Redis caches
	ctx := context.Background()
	cacheKeys := []string{"expenses", fmt.Sprintf("expense:%s", id)}
	for _, key := range cacheKeys {
		db.RedisClient.Del(ctx, key)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense updated successfully", "expense": expense})
}
