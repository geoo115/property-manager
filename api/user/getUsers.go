package user

import (
	"net/http"
	"strconv"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	var users []models.User
	var responseUsers []gin.H

	// Pagination: Get limit & offset from query parameters (default limit 10)
	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
	}
	if o := c.Query("offset"); o != "" {
		offset, _ = strconv.Atoi(o)
	}

	// Filtering: Get role from query parameters
	query := db.DB.Model(&models.User{})
	if role := c.Query("role"); role != "" {
		query = query.Where("role = ?", role)
	}

	// Sorting: Allow sorting by "username" or "created_at"
	sortBy := c.DefaultQuery("sort_by", "created_at") // Default: Sort by created_at
	order := c.DefaultQuery("order", "desc")          // Default: Descending order

	if err := query.Order(sortBy + " " + order).Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users"})
		return
	}

	// Remove sensitive data (passwords) before sending response
	for _, user := range users {
		responseUsers = append(responseUsers, gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
			"role":       user.Role,
			"phone":      user.Phone,
			"created_at": user.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"users":  responseUsers,
		"count":  len(responseUsers),
		"limit":  limit,
		"offset": offset,
	})
}
