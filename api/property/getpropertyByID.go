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

	// Find the property with the given ID
	var property models.Property
	if err := db.DB.Preload("Units").Preload("Owner").Preload("Tenant").First(&property, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}

	// Respond with the property data
	c.JSON(http.StatusOK, gin.H{
		"property": property,
	})
}
