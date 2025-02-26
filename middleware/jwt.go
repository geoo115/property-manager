package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			return
		}

		// Check if token is blacklisted in Redis
		blacklisted, _ := db.RedisClient.Get(db.Ctx, "blacklist:"+tokenString).Result()
		if blacklisted == "blacklisted" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if exp, ok := claims["exp"].(float64); ok {
				if time.Now().Unix() > int64(exp) {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
					return
				}
			}

			userID, ok1 := claims["userID"].(float64)
			role, ok2 := claims["role"].(string)
			username, ok3 := claims["username"].(string)

			if !ok1 || !ok2 || !ok3 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				return
			}

			c.Set("user_id", uint(userID))
			c.Set("user_role", role)
			c.Set("username", username)
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
	}
}
