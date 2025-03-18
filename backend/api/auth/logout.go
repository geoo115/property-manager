package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func LogoutHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
		return
	}

	// Blacklist token in Redis
	redisClient.Set(ctx, "blacklist:"+tokenString, "blacklisted", 1*time.Hour)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
