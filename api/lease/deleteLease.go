package lease

import (
	"context"
	"net/http"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

func DeleteLease(c *gin.Context) {
	id := c.Param("id")

	var lease models.Lease
	if err := db.DB.First(&lease, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lease not found"})
		return
	}

	if err := db.DB.Delete(&lease).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting lease"})
		return
	}
	db.RedisClient.FlushDB(context.Background())
	c.JSON(http.StatusOK, gin.H{"message": "Lease deleted successfully"})
}
