package middleware

import (
	"strings"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/logger"
	"github.com/geoo115/property-manager/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.LogWarning("Missing authorization header", logrus.Fields{
				"ip":     c.ClientIP(),
				"method": c.Request.Method,
				"path":   c.Request.URL.Path,
			})
			response.Unauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			logger.LogWarning("Invalid authorization header format", logrus.Fields{
				"ip":     c.ClientIP(),
				"method": c.Request.Method,
				"path":   c.Request.URL.Path,
			})
			response.Unauthorized(c, "Bearer token required")
			c.Abort()
			return
		}

		// Check if token is blacklisted in Redis
		blacklisted, _ := db.RedisClient.Get(db.Ctx, "blacklist:"+tokenString).Result()
		if blacklisted == "blacklisted" {
			logger.LogWarning("Blacklisted token used", logrus.Fields{
				"ip":    c.ClientIP(),
				"token": tokenString[:10] + "...", // Log only first 10 chars
			})
			response.Unauthorized(c, "Token has been revoked")
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil {
			logger.LogWarning("Invalid token", logrus.Fields{
				"ip":    c.ClientIP(),
				"error": err.Error(),
			})
			response.Unauthorized(c, "Invalid token")
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if exp, ok := claims["exp"].(float64); ok {
				if time.Now().Unix() > int64(exp) {
					logger.LogWarning("Expired token used", logrus.Fields{
						"ip":         c.ClientIP(),
						"expired_at": time.Unix(int64(exp), 0),
					})
					response.Unauthorized(c, "Token expired")
					c.Abort()
					return
				}
			}

			userID, ok1 := claims["userID"].(float64)
			role, ok2 := claims["role"].(string)
			username, ok3 := claims["username"].(string)

			if !ok1 || !ok2 || !ok3 {
				logger.LogWarning("Invalid token claims", logrus.Fields{
					"ip":     c.ClientIP(),
					"claims": claims,
				})
				response.Unauthorized(c, "Invalid token claims")
				c.Abort()
				return
			}

			c.Set("user_id", uint(userID))
			c.Set("user_role", role)
			c.Set("username", username)
			c.Next()
			return
		}

		logger.LogWarning("Token validation failed", logrus.Fields{
			"ip": c.ClientIP(),
		})
		response.Unauthorized(c, "Invalid token")
		c.Abort()
	}
}
