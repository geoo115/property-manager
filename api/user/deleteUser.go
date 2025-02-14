package user

import (
	"net/http"
	"strconv"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
)

func DeleteUser(c *gin.Context) {
	strid := c.Param("id")
	id, err := strconv.Atoi(strid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	if err := db.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting user"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
