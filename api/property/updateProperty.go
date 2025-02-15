package property

import (
	"net/http"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

func UpdateProperty(c *gin.Context) {
	// Get property ID from URL parameter
	id := c.Param("id")

	// Validate incoming JSON request body
	var input struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Bedrooms    uint    `json:"bedrooms"`
		Bathrooms   uint    `json:"bathrooms"`
		Price       float64 `json:"price"`
		SquareFeet  uint    `json:"square_feet"`
		Address     string  `json:"address"`
		City        string  `json:"city"`
		PostCode    string  `json:"post_code"`
		OwnerID     uint    `json:"owner_id"`
		Available   bool    `json:"available"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property data"})
		return
	}

	// Find the property to be updated
	var property models.Property
	if err := db.DB.First(&property, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}

	// Update property fields
	property.Name = input.Name
	property.Description = input.Description
	property.Bedrooms = input.Bedrooms
	property.Bathrooms = input.Bathrooms
	property.Price = input.Price
	property.SquareFeet = input.SquareFeet
	property.Address = input.Address
	property.City = input.City
	property.PostCode = input.PostCode
	property.OwnerID = input.OwnerID
	property.Available = input.Available

	// Save updated property to the database
	if err := db.DB.Save(&property).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating property"})
		return
	}

	// Respond with the updated property
	c.JSON(http.StatusOK, gin.H{
		"message":  "Property updated successfully",
		"property": property,
	})
}
