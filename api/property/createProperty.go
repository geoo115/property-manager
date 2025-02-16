package property

import (
	"net/http"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

func CreateProperty(c *gin.Context) {
	var input struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Bedrooms    uint    `json:"bedrooms" binding:"gte=0"`
		Bathrooms   uint    `json:"bathrooms" binding:"gte=0"`
		Price       float64 `json:"price" binding:"gte=0"`
		SquareFeet  uint    `json:"square_feet" binding:"gte=0"`
		Address     string  `json:"address" binding:"required"`
		City        string  `json:"city" binding:"required"`
		PostCode    string  `json:"post_code"`
		OwnerID     uint    `json:"owner_id" binding:"required"`
		Available   bool    `json:"available"`
	}

	// Bind JSON request
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property data", "details": err.Error()})
		return
	}

	// Check if Owner exists
	var owner models.User
	if err := db.DB.First(&owner, input.OwnerID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Owner does not exist"})
		return
	}

	// Create new property object
	property := models.Property{
		Name:        input.Name,
		Description: input.Description,
		Bedrooms:    input.Bedrooms,
		Bathrooms:   input.Bathrooms,
		Price:       input.Price,
		SquareFeet:  input.SquareFeet,
		Address:     input.Address,
		City:        input.City,
		PostCode:    input.PostCode,
		OwnerID:     input.OwnerID,
		Available:   input.Available,
	}

	// Insert into the database
	if err := db.DB.Create(&property).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating property"})
		return
	}

	if err := db.DB.Preload("Owner").First(&property, property.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching property owner"})
		return
	}

	// Return response with HTTP 201 Created
	c.JSON(http.StatusCreated, gin.H{
		"message":  "Property created successfully",
		"property": property,
	})
}
