package property

import (
	"net/http"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

func GetProperties(c *gin.Context) {
	var properties []models.Property

	// Apply filtering if query parameters are provided
	query := db.DB.Model(&models.Property{})

	// Filter by availability (optional)
	if available := c.Query("available"); available != "" {
		query = query.Where("available = ?", available == "true")
	}

	// Filter by city (optional)
	if city := c.Query("city"); city != "" {
		query = query.Where("city = ?", city)
	}

	// Filter by owner ID (optional)
	if ownerID := c.Query("owner_id"); ownerID != "" {
		query = query.Where("owner_id = ?", ownerID)
	}

	// Fetch properties along with related units, owner, and tenant (if any)
	if err := query.
		Preload("Units").
		Preload("Owner").
		Preload("Tenant").
		Find(&properties).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching properties"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"properties": properties})
}
