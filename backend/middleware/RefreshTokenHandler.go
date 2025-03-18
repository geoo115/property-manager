package middleware

import (
	"net/http"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RefreshTokenHandler(c *gin.Context) {
	// Read the refresh token from the cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token not found"})
		return
	}

	// Validate the refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	userID := uint(claims["userID"].(float64))

	// Fetch user from database
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	newAccessToken, err := GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}
