package property

import (
	"net/http"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

func DeleteProperty(c *gin.Context) {
	id := c.Param("id")

	var property models.Property
	if err := db.DB.First(&property, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}

	// Delete the property
	if err := db.DB.Delete(&property).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting property"})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "Property deleted successfully"})
}
