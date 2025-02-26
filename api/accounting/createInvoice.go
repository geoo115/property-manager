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

func CreateInvoice(c *gin.Context) {
	var input struct {
		TenantID       uint    `json:"tenant_id" binding:"required"`
		PropertyID     uint    `json:"property_id" binding:"required"`
		Amount         float64 `json:"amount" binding:"required,gt=0"`
		PaidAmount     float64 `json:"paid_amount" binding:"required,gte=0"`
		InvoiceDateStr string  `json:"invoice_date" binding:"required"`
		Category       string  `json:"category" binding:"required,oneof=rent deposit"`
		DueDateStr     string  `json:"due_date" binding:"required"`
		PaymentStatus  string  `json:"payment_status" binding:"required"`
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

	// Validate tenant and property
	var tenant models.User
	if err := db.DB.First(&tenant, input.TenantID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		return
	}

	var property models.Property
	if err := db.DB.First(&property, input.PropertyID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}

	invoice := models.Invoice{
		TenantID:      input.TenantID,
		PropertyID:    input.PropertyID,
		Amount:        input.Amount,
		PaidAmount:    input.PaidAmount,
		InvoiceDate:   invoiceDate,
		Category:      input.Category,
		DueDate:       dueDate,
		PaymentStatus: input.PaymentStatus,
	}

	// Create invoice
	if err := db.DB.Create(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating invoice", "details": err.Error()})
		return
	}

	// Preload associations
	if err := db.DB.Preload("Tenant").Preload("Property.Owner").First(&invoice, invoice.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching created invoice", "details": err.Error()})
		return
	}

	// Invalidate Redis caches
	ctx := context.Background()
	cacheKeys := []string{
		"invoices",
		fmt.Sprintf("invoice:%d", invoice.ID),
		fmt.Sprintf("tenant_invoices:%d", invoice.TenantID),
		fmt.Sprintf("landlord_invoices:%d", property.OwnerID),
	}
	for _, key := range cacheKeys {
		if err := db.RedisClient.Del(ctx, key).Err(); err != nil {
			fmt.Printf("Failed to delete Redis key %s: %v\n", key, err)
		}
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Invoice created successfully",
		"invoice": invoice,
	})
}
