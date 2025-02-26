package accounting

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

// CreateExpense creates a new expense with Redis cache invalidation.
func CreateExpense(c *gin.Context) {
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

	// Parse expense date with fallback to YYYY-MM-DD
	expenseDate, err := time.Parse(time.RFC3339, input.ExpenseDate)
	if err != nil {
		expenseDate, err = time.Parse("2006-01-02", input.ExpenseDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use RFC3339 or YYYY-MM-DD"})
			return
		}
	}

	expense := models.Expense{
		PropertyID:  input.PropertyID,
		Description: input.Description,
		Category:    input.Category,
		Amount:      input.Amount,
		ExpenseDate: expenseDate,
	}

	if err := db.DB.Create(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating expense"})
		return
	}

	// Preload associations
	if err := db.DB.Preload("Property.Owner").First(&expense, expense.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error loading property and owner"})
		return
	}

	// Invalidate relevant Redis caches
	ctx := context.Background()
	cacheKeys := []string{"expenses", fmt.Sprintf("expense:%d", expense.ID)}
	for _, key := range cacheKeys {
		db.RedisClient.Del(ctx, key)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Expense created successfully",
		"expense": expense,
	})
}
