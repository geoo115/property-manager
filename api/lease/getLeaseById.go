package lease

import (
	"net/http"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

func GetLeaseByID(c *gin.Context) {
	id := c.Param("id")

	var lease models.Lease
	if err := db.DB.Preload("Property").Preload("Tenant").First(&lease, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "lease not Found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"lease": lease,
	})
}
