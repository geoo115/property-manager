package user

import (
	"net/http"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/geoo115/property-manager/utils"
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var input struct {
		Username  string `json:"username" binding:"required"`
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		Password  string `json:"password" binding:"required,min=6"`
		Email     string `json:"email" binding:"required,email"`
		Role      string `json:"role" binding:"required,oneof=admin tenant landlord maintenanceTeam"`
		Phone     string `json:"phone" binding:"required"`
	}

	// Bind JSON request to the temporary struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	// Hash password before storing it
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	// Create user object from validated input
	user := models.User{
		Username:  input.Username,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Password:  hashedPassword, // Store the hashed password
		Email:     input.Email,
		Role:      input.Role,
		Phone:     input.Phone,
	}

	// Insert into database
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	// Send response (omit password for security)
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
			"role":       user.Role,
			"phone":      user.Phone,
			"created_at": user.CreatedAt,
		},
	})
}
