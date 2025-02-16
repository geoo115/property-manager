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
		c.JSON(http.StatusNotFound, gin.H{"error": "Lease not found"})
		return
	}

	// Retrieve user role and ID from context
	userRole, _ := c.Get("user_role")
	userID, _ := c.Get("user_id") // Assuming stored as uint in context

	// Enforce role-based access
	if userRole == "tenant" && lease.TenantID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if userRole == "landlord" && lease.Property.OwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"lease": lease,
	})
}
