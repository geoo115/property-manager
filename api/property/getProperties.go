package property

import (
	"net/http"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

func GetProperties(c *gin.Context) {
	var properties []models.Property

	// Get user role and ID from context (set by JWT middleware)
	userRole, _ := c.Get("user_role")
	userID, _ := c.Get("user_id") // Assuming stored as uint in context

	// Start building query
	query := db.DB.Model(&models.Property{})

	// Role-based filtering
	if userRole == "landlord" {
		// Landlords can only access their own properties
		query = query.Where("owner_id = ?", userID)
	}

	// Apply optional filters
	if available := c.Query("available"); available != "" {
		query = query.Where("available = ?", available == "true")
	}

	if city := c.Query("city"); city != "" {
		query = query.Where("city = ?", city)
	}

	if ownerID := c.Query("owner_id"); ownerID != "" && userRole == "admin" {
		// Only admins can filter by owner_id
		query = query.Where("owner_id = ?", ownerID)
	}

	// Fetch properties along with related units, owner
	if err := query.Preload("Units").Preload("Owner").
		Find(&properties).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching properties"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"properties": properties})
}
