package lease

import (
	"net/http"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

func GetActiveLeaseForTenant(c *gin.Context) {
	// Get the user ID from the token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
		return
	}

	// Query the database for the tenant's active lease and preload tenant and property owner
	var lease models.Lease
	if err := db.DB.Preload("Tenant").Preload("Property.Owner").
		Where("tenant_id = ? AND end_date > ?", userID, time.Now()).First(&lease).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active lease found for tenant"})
		return
	}

	c.JSON(http.StatusOK, lease)
}

func GetLeasesForTenant(c *gin.Context) {
	// Use "user_id" if that’s the key set in the context.
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, user_id not found"})
		return
	}

	// Ensure userID is of type uint, as it's expected in your database query
	tenantID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	// Query to get leases for the tenant
	var leases []models.Lease
	if err := db.DB.Where("tenant_id = ?", tenantID).Preload("Property").Find(&leases).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching leases"})
		return
	}

	// Return the leases
	c.JSON(http.StatusOK, gin.H{"leases": leases})
}
