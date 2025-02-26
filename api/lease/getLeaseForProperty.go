package lease

import (
	"net/http"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetLeaseForProperty(c *gin.Context) {
	propertyID := c.Param("id") // Get the property ID from URL

	// Logic to retrieve the lease for this property
	lease, err := GetLeaseByPropertyID(propertyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lease not found"})
		return
	}

	c.JSON(http.StatusOK, lease)
}

func GetLeaseByPropertyID(propertyID string) (*models.Lease, error) {
	var lease models.Lease

	// Fetch lease associated with the given property ID from the database
	result := db.DB.Where("property_id = ?", propertyID).First(&lease)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // No lease found for this property
		}
		return nil, result.Error // Some other error occurred
	}

	return &lease, nil
}
