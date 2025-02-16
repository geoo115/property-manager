package lease

import (
	"net/http"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

func UpdateLease(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		TenantID        uint      `json:"tenant_id" binding:"required"`
		PropertyID      uint      `json:"property_id" binding:"required"`
		StartDate       time.Time `json:"start_date" binding:"required"`
		EndDate         time.Time `json:"end_date" binding:"required"`
		MonthlyRent     float64   `json:"monthly_rent" binding:"required,gte=0"`
		SecurityDeposit float64   `json:"security_deposit" binding:"required,gte=0"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lease input", "details": err.Error()})
		return
	}

	if !input.EndDate.After(input.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date must be after start date"})
		return
	}

	var lease models.Lease
	if err := db.DB.Preload("Tenant").Preload("Property.Owner").
		First(&lease, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lease not found"})
		return
	}

	// Enforce role-based access
	userRole, _ := c.Get("user_role")
	userID, _ := c.Get("user_id")

	if userRole == "tenant" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Tenants cannot update leases"})
		return
	}
	if userRole == "landlord" && lease.Property.OwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Update lease
	lease.TenantID = input.TenantID
	lease.PropertyID = input.PropertyID
	lease.StartDate = input.StartDate
	lease.EndDate = input.EndDate
	lease.MonthlyRent = input.MonthlyRent
	lease.SecurityDeposit = input.SecurityDeposit

	if err := db.DB.Save(&lease).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating lease"})
		return
	}

	// Reload with preloaded data
	if err := db.DB.Preload("Tenant").Preload("Property.Owner").
		First(&lease, lease.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching updated lease"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lease updated successfully",
		"lease":   lease,
	})
}
