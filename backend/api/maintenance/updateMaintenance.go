package maintenance

import (
	"context"
	"fmt"
	"net/http"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

// UpdateMaintenance updates a maintenance request and invalidates Redis caches.
func UpdateMaintenance(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Description *string `json:"description"`
		PropertyID  *uint   `json:"property_id"`
		Status      *string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid maintenance data", "details": err.Error()})
		return
	}

	var maintenance models.Maintenance
	if err := db.DB.First(&maintenance, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Maintenance not found"})
		return
	}

	// Update only provided fields
	if input.Description != nil {
		maintenance.Description = *input.Description
	}
	if input.PropertyID != nil {
		maintenance.PropertyID = *input.PropertyID
	}
	if input.Status != nil {
		maintenance.Status = *input.Status
	}

	if err := db.DB.Save(&maintenance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating maintenance request"})
		return
	}

	// Reload with associations
	if err := db.DB.Preload("RequestedBy").Preload("Property.Owner").First(&maintenance, maintenance.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching updated details"})
		return
	}

	// Invalidate Redis caches
	ctx := context.Background()
	cacheKeys := []string{
		"maintenances:all",
		"maintenances:team",
		fmt.Sprintf("maintenance:%s", id),
		// Additional invalidation based on role/property could be added if tracked
	}
	for _, key := range cacheKeys {
		if err := db.RedisClient.Del(ctx, key).Err(); err != nil {
			fmt.Printf("Failed to delete Redis key %s: %v\n", key, err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Maintenance updated successfully",
		"maintenance": maintenance,
	})
}
