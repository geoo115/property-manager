package user

import (
	"net/http"
	"strconv"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/geoo115/property-manager/utils"
	"github.com/gin-gonic/gin"
)

func UpdateUser(c *gin.Context) {
	strid := c.Param("id")
	id, err := strconv.Atoi(strid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	var updatePayload models.User
	if err := c.ShouldBindJSON(&updatePayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}

	if updatePayload.Password != "" {
		hashedPassword, err := utils.HashPassword(updatePayload.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error hashing password"})
			return
		}
		updatePayload.Password = hashedPassword
	}

	if err := db.DB.Model(&user).Updates(updatePayload).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating user"})
		return
	}

	c.JSON(http.StatusOK, user)
}
