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

// DeleteInvoice deletes an invoice and invalidates Redis cache.
func DeleteInvoice(c *gin.Context) {
	id := c.Param("id")
	var invoice models.Invoice

	if err := db.DB.First(&invoice, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching invoice"})
		}
		return
	}

	if err := db.DB.Delete(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting invoice"})
		return
	}

	// Invalidate Redis caches
	ctx := context.Background()
	cacheKeys := []string{"invoices", fmt.Sprintf("invoice:%s", id)}
	for _, key := range cacheKeys {
		db.RedisClient.Del(ctx, key)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invoice deleted successfully"})
}
