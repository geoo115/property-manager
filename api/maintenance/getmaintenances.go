package maintenance

import (
	"net/http"
	"strconv"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

func GetMaintenances(c *gin.Context) {
	var maintenances []models.Maintenance
	userRole, _ := c.Get("user_role")
	userID, _ := c.Get("user_id")

	query := db.DB.Model(&models.Maintenance{})

	// ✅ Admins can see all maintenance requests
	if userRole == "admin" {
		query = query.Preload("Tenant").Preload("Property")
	} else if userRole == "tenant" {
		// ✅ Tenants should only see maintenance for their **active leases**
		query = query.Joins("JOIN leases ON leases.property_id = maintenances.property_id").
			Where("leases.tenant_id = ? AND leases.end_date > ?", userID, time.Now()).
			Preload("Property")
	} else if userRole == "maintenanceTeam" {
		// ✅ Maintenance team sees all maintenance requests
		query = query.Preload("Property").Preload("Tenant")
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	if err := query.Find(&maintenances).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching maintenance requests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"maintenances": maintenances})
}

func GetLandlordMaintenances(c *gin.Context) {
	propertyIDStr := c.Param("id")
	propertyID, err := strconv.ParseUint(propertyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property ID"})
		return
	}

	userID, _ := c.Get("user_id")

	var maintenances []models.Maintenance

	// Verify that the landlord owns the property.
	var property models.Property
	if err := db.DB.Where("id = ? AND owner_id = ?", propertyID, userID).First(&property).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found or you do not own this property"})
		return
	}

	if err := db.DB.Where("property_id = ?", propertyID).Preload("Tenant").Preload("Property.Owner").Find(&maintenances).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch maintenances"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"maintenances": maintenances})
}
