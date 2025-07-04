package maintenance

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/events"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

// CreateMaintenanceByProperty creates a maintenance request for a property and invalidates Redis caches.
func CreateMaintenanceByProperty(c *gin.Context) {
	propertyIDStr := c.Param("propertyID")
	propertyID, err := strconv.ParseUint(propertyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property ID"})
		return
	}

	var input struct {
		Description string `json:"description" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid maintenance data", "details": err.Error()})
		return
	}

	userRole, _ := c.Get("user_role")
	userID, _ := c.Get("user_id")

	if userRole == "landlord" {
		var property models.Property
		if err := db.DB.First(&property, propertyID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Property not found"})
			return
		}
		if property.OwnerID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not own this property"})
			return
		}
	}

	maintenance := models.Maintenance{
		RequestedByID: userID.(uint),
		PropertyID:    uint(propertyID),
		Description:   input.Description,
		RequestedAt:   time.Now(),
		Status:        "pending",
	}

	if err := db.DB.Create(&maintenance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating maintenance request"})
		return
	}

	if err := db.DB.Preload("RequestedBy").Preload("Property.Owner").First(&maintenance, maintenance.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching maintenance details"})
		return
	}

	if err := events.ProduceMaintenanceRequest(maintenance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send Kafka event"})
		return
	}

	// Invalidate Redis caches
	ctx := context.Background()
	cacheKeys := []string{
		"maintenances:all", // Admin view
		fmt.Sprintf("maintenances:landlord:%d:property:%d", userID.(uint), propertyID),
		"maintenances:team", // Maintenance team view
		fmt.Sprintf("maintenance:%d", maintenance.ID),
	}
	for _, key := range cacheKeys {
		if err := db.RedisClient.Del(ctx, key).Err(); err != nil {
			fmt.Printf("Failed to delete Redis key %s: %v\n", key, err)
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Maintenance request created successfully",
		"maintenance": maintenance,
	})
}

func CreateMaintenanceByLease(c *gin.Context) {
	// Log all URL parameters for debugging
	log.Printf("All URL params: %+v", c.Params)

	leaseID := c.Param("leaseID")
	if leaseID == "" {
		log.Printf("leaseID is empty, falling back to check other params")
		// Fallback check for common param names
		leaseID = c.Param("id") // Alternative if misnamed in route
		if leaseID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing lease ID in URL"})
			return
		}
	}
	log.Printf("Received leaseID from URL: %s", leaseID)

	var input struct {
		Description string `json:"description" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid maintenance data", "details": err.Error()})
		return
	}

	var lease models.Lease
	log.Printf("Fetching lease with ID: %s", leaseID)
	if err := db.DB.Where("id = ?", leaseID).First(&lease).Error; err != nil {
		log.Printf("Lease fetch error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lease not found"})
		return
	}
	log.Printf("Fetched lease: %+v", lease)

	if lease.EndDate.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No active lease found for tenant on this property"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
		return
	}
	log.Printf("User ID from token: %d", userID.(uint))

	// Validate tenant ownership
	if lease.TenantID != userID.(uint) {
		log.Printf("Authorization failed: Lease TenantID %d does not match UserID %d", lease.TenantID, userID.(uint))
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized for this lease"})
		return
	}

	maintenance := models.Maintenance{
		RequestedByID: userID.(uint),
		PropertyID:    lease.PropertyID,
		Description:   input.Description,
		RequestedAt:   time.Now(),
		Status:        "pending",
	}

	if err := db.DB.Create(&maintenance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating maintenance request"})
		return
	}

	if err := db.DB.Preload("RequestedBy").Preload("Property.Owner").First(&maintenance, maintenance.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching maintenance details"})
		return
	}

	if err := events.ProduceMaintenanceRequest(maintenance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send Kafka event"})
		return
	}

	// Invalidate Redis caches
	ctx := context.Background()
	cacheKeys := []string{
		"maintenances:all",
		fmt.Sprintf("maintenances:tenant:%d:lease:%s", userID.(uint), leaseID),
		"maintenances:team",
		fmt.Sprintf("maintenance:%d", maintenance.ID),
	}
	for _, key := range cacheKeys {
		if err := db.RedisClient.Del(ctx, key).Err(); err != nil {
			fmt.Printf("Failed to delete Redis key %s: %v\n", key, err)
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Maintenance request created successfully",
		"maintenance": maintenance,
	})
}
