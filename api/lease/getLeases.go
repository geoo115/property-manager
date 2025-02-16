package lease

import (
	"net/http"
	"strconv"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

func GetLeases(c *gin.Context) {
	var leases []models.Lease

	limit := 10
	offset := 0

	// Parse limit & offset
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	// Retrieve user role and ID from context
	userRole, _ := c.Get("user_role")
	userID, _ := c.Get("user_id")

	query := db.DB.Model(&models.Lease{}).Preload("Tenant").Preload("Property.Owner")

	// Role-based access control
	if userRole == "tenant" {
		query = query.Where("tenant_id = ?", userID)
	} else if userRole == "landlord" {
		query = query.Joins("JOIN properties ON properties.id = leases.property_id").
			Where("properties.owner_id = ?", userID)
	}

	// Filtering
	if propertyID := c.Query("property_id"); propertyID != "" {
		query = query.Where("leases.property_id = ?", propertyID)
	}

	if startAfter := c.Query("start_after"); startAfter != "" {
		if t, err := time.Parse(time.RFC3339, startAfter); err == nil {
			query = query.Where("leases.start_date >= ?", t)
		}
	}

	// Sorting
	sortBy := c.DefaultQuery("sort_by", "leases.start_date")
	order := c.DefaultQuery("order", "desc")
	query = query.Order(sortBy + " " + order)

	// Fetch the leases
	if err := query.Limit(limit).Offset(offset).Find(&leases).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching leases"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"leases": leases,
		"limit":  limit,
		"offset": offset,
	})
}
