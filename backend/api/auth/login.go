package auth

import (
	"context"
	"net/http"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/middleware"
	"github.com/geoo115/property-manager/models"
	"github.com/geoo115/property-manager/utils"
	"github.com/gin-gonic/gin"
)

// Redis client for rate limiting
var ctx = context.Background()
var redisClient = db.RedisClient

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

	if !utils.Comparepassword(user.Password, Credentials.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate Access Token
	accessToken, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating access token"})
		return
	}

	// Generate Refresh Token
	refreshToken, err := middleware.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating refresh token"})
		return
	}

	// Store refresh token in secure HttpOnly cookie
	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "localhost", false, true)

	// Respond with access token
	c.JSON(http.StatusOK, gin.H{
		"message":      "Login successful",
		"access_token": accessToken,
	})
}
