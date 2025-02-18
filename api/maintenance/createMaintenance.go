package maintenance

import (
	"net/http"
	"strconv"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

func CreateMaintenance(c *gin.Context) {
	// Extract lease ID from URL
	leaseIDStr := c.Param("id")
	leaseID, err := strconv.Atoi(leaseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lease ID"})
		return
	}

	// For tenant requests, we don't require tenant_id or property_id from the payload.
	var input struct {
		Description string    `json:"description" binding:"required"`
		RequestedAt time.Time `json:"requested_at"`
		Status      string    `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid maintenance data", "details": err.Error()})
		return
	}

	// Set default values
	if input.RequestedAt.IsZero() {
		input.RequestedAt = time.Now()
	}
	if input.Status == "" {
		input.Status = "pending"
	}

	// Retrieve the lease using the lease ID.
	var lease models.Lease
	if err := db.DB.First(&lease, leaseID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lease not found"})
		return
	}

	// Check if the lease is active (i.e., its EndDate is in the future).
	if lease.EndDate.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No active lease found for tenant on this property"})
		return
	}

	// Use the lease's TenantID and PropertyID.
	tenantID := lease.TenantID
	// Fallback: if lease.TenantID is 0, try using the user_id from the token.
	if tenantID == 0 {
		if uid, exists := c.Get("user_id"); exists {
			tenantID = uid.(uint)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant not found in token"})
			return
		}
	}

	maintenance := models.Maintenance{
		TenantID:    tenantID,
		PropertyID:  lease.PropertyID,
		Description: input.Description,
		RequestedAt: input.RequestedAt,
		Status:      input.Status,
	}

	if err := db.DB.Create(&maintenance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating maintenance request"})
		return
	}

	// Preload associations and return the created record.
	if err := db.DB.Preload("Tenant").Preload("Property.Owner").
		First(&maintenance, maintenance.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching maintenance details"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Maintenance request created successfully",
		"maintenance": maintenance,
	})
}

// CreateLandlordMaintenance creates a maintenance request for a landlord's property.
func CreateLandlordMaintenance(c *gin.Context) {
	propertyIDStr := c.Param("id")
	propertyID, err := strconv.ParseUint(propertyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property ID"})
		return
	}

	userID, _ := c.Get("user_id")

	// Verify that the landlord owns the property.
	var property models.Property
	if err := db.DB.Where("id = ? AND owner_id = ?", propertyID, userID).First(&property).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found or you do not own this property"})
		return
	}

	var input struct {
		TenantID    uint      `json:"tenant_id"`
		Description string    `json:"description" binding:"required"`
		RequestedAt time.Time `json:"requested_at"`
		Status      string    `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid maintenance data", "details": err.Error()})
		return
	}

	if input.RequestedAt.IsZero() {
		input.RequestedAt = time.Now()
	}

	if input.Status == "" {
		input.Status = "pending"
	}

	maintenance := models.Maintenance{
		TenantID:    input.TenantID,
		PropertyID:  uint(propertyID),
		Description: input.Description,
		RequestedAt: input.RequestedAt,
		Status:      input.Status,
	}

	if err := db.DB.Create(&maintenance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating maintenance request"})
		return
	}

	if err := db.DB.Preload("Tenant").Preload("Property.Owner").First(&maintenance, maintenance.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching maintenance details"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Maintenance request created successfully",
		"maintenance": maintenance,
	})
}
