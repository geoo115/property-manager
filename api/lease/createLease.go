package lease

import (
	"net/http"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

func CreateLease(c *gin.Context) {
	var input struct {
		TenantID        uint      `json:"tenant_id" binding:"required"`
		PropertyID      uint      `json:"property_id" binding:"required"`
		StartDate       time.Time `json:"start_date" binding:"required"`
		EndDate         time.Time `json:"end_date" binding:"required"`
		MonthlyRent     float64   `json:"monthly_rent" binding:"required,gte=0"`
		SecurityDeposit float64   `json:"security_deposit" binding:"required,gte=0"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lease data", "details": err.Error()})
		return
	}

	if !input.EndDate.After(input.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "EndDate must be after StartDate"})
		return
	}

	var tenant models.User
	if err := db.DB.First(&tenant, input.TenantID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant does not exist"})
		return
	}

	var property models.Property
	if err := db.DB.First(&property, input.PropertyID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Property not found"})
		return
	}

	// Create the lease
	lease := models.Lease{
		TenantID:        input.TenantID,
		PropertyID:      input.PropertyID,
		StartDate:       input.StartDate,
		EndDate:         input.EndDate,
		MonthlyRent:     input.MonthlyRent,
		SecurityDeposit: input.SecurityDeposit,
	}

	if err := db.DB.Create(&lease).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating lease"})
		return
	}

	// Reload the lease with preloaded Tenant & Property data
	if err := db.DB.Preload("Tenant").Preload("Property.Owner").First(&lease, lease.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching lease details"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Lease created successfully",
		"lease":   lease,
	})
}
