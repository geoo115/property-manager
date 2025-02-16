package property

import (
	"net/http"
	"strconv"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

func GetProperties(c *gin.Context) {
	var properties []models.Property

	// Pagination: Set default values
	limit := 10
	offset := 0

	// Parse limit & offset from query parameters
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	// Get user role and ID from context (set by JWT middleware)
	userRole, _ := c.Get("user_role")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	landlordID, ok := userID.(uint) // Ensure correct type
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Start query builder
	query := db.DB.Model(&models.Property{})

	// Role-based filtering
	if userRole == "landlord" {
		// Landlords can only access their own properties
		query = query.Where("owner_id = ?", landlordID)
	}

	// Apply optional filters
	if available := c.Query("available"); available != "" {
		query = query.Where("available = ?", available == "true")
	}
	if city := c.Query("city"); city != "" {
		query = query.Where("city = ?", city)
	}
	if ownerID := c.Query("owner_id"); ownerID != "" && userRole == "admin" {
		// Only admins can filter by owner_id
		query = query.Where("owner_id = ?", ownerID)
	}

	// Execute query with pagination
	if err := query.Preload("Units").Preload("Owner").
		Limit(limit).Offset(offset).
		Find(&properties).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching properties"})
		return
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"properties": properties,
		"limit":      limit,
		"offset":     offset,
		"count":      len(properties),
	})
}
