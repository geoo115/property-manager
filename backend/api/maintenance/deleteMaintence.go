package maintenance

import (
	"context"
	"fmt"
	"net/http"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

// DeleteMaintenance deletes a maintenance request and invalidates Redis caches.
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

	// Invalidate Redis caches
	ctx := context.Background()
	cacheKeys := []string{
		"maintenances:all",
		"maintenances:team",
		fmt.Sprintf("maintenance:%s", id),
		// Additional keys based on role/property could be added if needed
	}
	for _, key := range cacheKeys {
		if err := db.RedisClient.Del(ctx, key).Err(); err != nil {
			fmt.Printf("Failed to delete Redis key %s: %v\n", key, err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Maintenance request deleted successfully"})
}
