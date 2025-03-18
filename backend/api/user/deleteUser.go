package user

import (
	"context"
	"net/http"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := db.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting user"})
	}
	db.RedisClient.FlushDB(context.Background())
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
