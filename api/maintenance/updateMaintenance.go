package maintenance

import (
	"net/http"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

// UpdateMaintenance updates an existing maintenance request.
func UpdateMaintenance(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		TenantID    uint   `json:"tenant_id" binding:"required"`
		PropertyID  uint   `json:"property_id" binding:"required"`
		Description string `json:"description" binding:"required"`
		RequestedAt string `json:"requested_at" binding:"required"`
		Status      string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid maintenance data", "details": err.Error()})
		return
	}

	// Parse the RequestedAt string into a time.Time
	parsedTime, err := time.Parse("2006-01-02 15:04:05-07:00", input.RequestedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid requested_at format",
			"details": "Expected format: 'YYYY-MM-DD HH:MM:SS+ZZ:ZZ' (e.g., '2025-02-18 00:00:00+00:00')",
		})
		return
	}

	var maintenance models.Maintenance
	if err := db.DB.First(&maintenance, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Maintenance not found"})
		return
	}

	// Update fields
	maintenance.TenantID = input.TenantID
	maintenance.PropertyID = input.PropertyID
	maintenance.Description = input.Description
	maintenance.RequestedAt = parsedTime
	maintenance.Status = input.Status

	if err := db.DB.Save(&maintenance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating maintenance request"})
		return
	}

	// Reload with preloaded associations.
	if err := db.DB.Preload("Tenant").Preload("Property.Owner").
		First(&maintenance, maintenance.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching updated maintenance details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Maintenance request updated successfully",
		"maintenance": maintenance,
	})
}
