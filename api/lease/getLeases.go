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

	query := db.DB.Model(&models.Lease{})

	if tenantID := c.Query("tenant_id"); tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if propertyID := c.Query("property_id"); propertyID != "" {
		query = query.Where("property_id = ?", propertyID)
	}

	if startAfter := c.Query("start_after"); startAfter != "" {
		if t, err := time.Parse(time.RFC3339, startAfter); err == nil {
			query = query.Where("start_date >= ?", t)
		}
	}

	sortBy := c.DefaultQuery("sort_by", "start_date")
	order := c.DefaultQuery("order", "desc")
	query = query.Order(sortBy + " " + order)

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
