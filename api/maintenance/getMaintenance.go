package maintenance

import (
	"net/http"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

func GetMaintenance(c *gin.Context) {
	id := c.Param("id")

	var maintenance models.Maintenance
	if err := db.DB.Preload("Tenant").Preload("Property.Owner").
		First(&maintenance, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Maintenance not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"maintenance": maintenance})
}
