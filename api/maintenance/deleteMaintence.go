package maintenance

import (
	"net/http"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

func DeleteMaintenance(c *gin.Context) {
	id := c.Param("id")

	var maintenance models.Maintenance
	if err := db.DB.First(&maintenance, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Maintenance request not found"})
		return
	}

	if err := db.DB.Delete(&maintenance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting maintenance request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Maintenance request deleted successfully"})
}
