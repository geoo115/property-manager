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

func UpdateInvoice(c *gin.Context) {
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

	var input struct {
		TenantID          uint    `json:"tenant_id" binding:"required"`
		PropertyID        uint    `json:"property_id" binding:"required"`
		Amount            float64 `json:"amount" binding:"required,gt=0"`
		PaidAmount        float64 `json:"paid_amount" binding:"required,gte=0"`
		InvoiceDateStr    string  `json:"invoice_date" binding:"required"`
		Category          string  `json:"category" binding:"required,oneof=rent deposit"`
		DueDateStr        string  `json:"due_date" binding:"required"`
		PaymentStatus     string  `json:"payment_status" binding:"required"`
		RefundedAmount    float64 `json:"refunded_amount" binding:"gte=0"` // Removed 'required'
		RecurringInterval string  `json:"recurring_interval"`
		Recurring         bool    `json:"recurring"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice data", "details": err.Error()})
		return
	}

	// Parse dates
	invoiceDate, err := time.Parse("2006-01-02", input.InvoiceDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice date format, use YYYY-MM-DD", "details": err.Error()})
		return
	}
	dueDate, err := time.Parse("2006-01-02", input.DueDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due date format, use YYYY-MM-DD", "details": err.Error()})
		return
	}

	invoice.TenantID = input.TenantID
	invoice.PropertyID = input.PropertyID
	invoice.Amount = input.Amount
	invoice.PaidAmount = input.PaidAmount
	invoice.InvoiceDate = invoiceDate
	invoice.Category = input.Category
	invoice.DueDate = dueDate
	invoice.PaymentStatus = input.PaymentStatus
	invoice.RefundedAmount = input.RefundedAmount
	invoice.RecurringInterval = input.RecurringInterval
	invoice.Recurring = input.Recurring

	if err := db.DB.Save(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating invoice"})
		return
	}

	if err := db.DB.Preload("Tenant").Preload("Property").First(&invoice, invoice.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching updated invoice"})
		return
	}

	// Invalidate Redis caches
	ctx := context.Background()
	cacheKeys := []string{
		"invoices",
		fmt.Sprintf("invoice:%s", id),
		fmt.Sprintf("tenant_invoices:%d", invoice.TenantID),
		fmt.Sprintf("landlord_invoices:%d", invoice.Property.OwnerID),
	}
	for _, key := range cacheKeys {
		if err := db.RedisClient.Del(ctx, key).Err(); err != nil {
			fmt.Printf("Failed to delete Redis key %s: %v\n", key, err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Invoice updated successfully",
		"invoice": invoice,
	})
}
