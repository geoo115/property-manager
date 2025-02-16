package property

import (
	"net/http"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

func GetPropertyByID(c *gin.Context) {
	// Get property ID from URL parameter
	id := c.Param("id")

	// Retrieve user role and ID from context
	userRole, _ := c.Get("user_role")
	userID, _ := c.Get("user_id")

	// Find the property
	var property models.Property
	if err := db.DB.Preload("Units").Preload("Owner").
		First(&property, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}

	// Enforce role-based access
	if userRole == "landlord" && property.OwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Respond with property data
	c.JSON(http.StatusOK, gin.H{
		"property": property,
	})
}
