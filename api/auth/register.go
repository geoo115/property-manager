package auth

import (
	"net/http"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/geoo115/property-manager/utils"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if user.Username == "" || user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username, Password, and Email are required"})
		return
	}

	var userExists models.User
	if err := db.DB.Where("email=?", user.Email).First(&userExists).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	if err := db.DB.Where("username=?", user.Username).First(&userExists).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	user.Password = hashedPassword

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}
}
