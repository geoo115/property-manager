package auth

import (
	"net/http"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/middleware"
	"github.com/geoo115/property-manager/models"
	"github.com/geoo115/property-manager/utils"
	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	var Credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&Credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if Credentials.Username == "" || Credentials.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and Password are required"})
		return
	}
	var user models.User
	if err := db.DB.Where("username=?", Credentials.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		return
	}
	if !utils.ComparePassword(user.Password, Credentials.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		return
	}
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Login successfuly", "token": token})
}
