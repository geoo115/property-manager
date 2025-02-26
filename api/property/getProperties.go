package property

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

// GetProperties godoc
// @Summary List properties
// @Description Get a paginated list of properties based on user role and filters
// @Tags Properties
// @Accept json
// @Produce json
// @Param limit query int false "Number of items per page" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Param available query boolean false "Filter by availability"
// @Param city query string false "Filter by city"
// @Param owner_id query int false "Filter by owner ID (admin only)"
// @Security BearerAuth
// @Success 200 {object} PropertyListResponse "List of properties"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /admin/properties [get]
// @Router /landlord/properties [get]
func GetProperties(c *gin.Context) {
	var properties []models.Property

	limit := 10
	offset := 0

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

	userRole, _ := c.Get("user_role")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	landlordID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	cacheKey := "properties:role=" + userRole.(string) + ":limit=" + strconv.Itoa(limit) + ":offset=" + strconv.Itoa(offset)
	if userRole == "landlord" {
		cacheKey += ":owner_id=" + strconv.Itoa(int(landlordID))
	}
	if available := c.Query("available"); available != "" {
		cacheKey += ":available=" + available
	}
	if city := c.Query("city"); city != "" {
		cacheKey += ":city=" + city
	}
	if ownerID := c.Query("owner_id"); ownerID != "" && userRole == "admin" {
		cacheKey += ":owner_id=" + ownerID
	}

	ctx := context.Background()
	cachedData, err := db.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedData), &properties); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"properties": properties,
				"limit":      limit,
				"offset":     offset,
				"count":      len(properties),
				"cache":      "hit",
			})
			return
		}
	}

	query := db.DB.Model(&models.Property{})
	if userRole == "landlord" {
		query = query.Where("owner_id = ?", landlordID)
	}

	if available := c.Query("available"); available != "" {
		query = query.Where("available = ?", available == "true")
	}
	if city := c.Query("city"); city != "" {
		query = query.Where("city = ?", city)
	}
	if ownerID := c.Query("owner_id"); ownerID != "" && userRole == "admin" {
		query = query.Where("owner_id = ?", ownerID)
	}

	if err := query.Preload("Units").Preload("Owner").
		Limit(limit).Offset(offset).
		Find(&properties).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching properties"})
		return
	}

	jsonData, _ := json.Marshal(properties)
	db.RedisClient.Set(ctx, cacheKey, jsonData, 10*time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"properties": properties,
		"limit":      limit,
		"offset":     offset,
		"count":      len(properties),
		"cache":      "miss",
	})
}

// PropertyListResponse defines the response structure
type PropertyListResponse struct {
	Properties []models.Property `json:"properties"`
	Limit      int               `json:"limit" example:"10"`
	Offset     int               `json:"offset" example:"0"`
	Count      int               `json:"count" example:"5"`
	Cache      string            `json:"cache" example:"miss"`
}
