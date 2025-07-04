package auth

import (
	"context"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/logger"
	"github.com/geoo115/property-manager/middleware"
	"github.com/geoo115/property-manager/models"
	"github.com/geoo115/property-manager/response"
	"github.com/geoo115/property-manager/utils"
	"github.com/geoo115/property-manager/validator"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// Redis client for rate limiting
var ctx = context.Background()
var redisClient = db.RedisClient

func LoginHandler(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		logger.LogError(err, "Invalid login request", logrus.Fields{
			"ip": c.ClientIP(),
		})
		response.BadRequest(c, "Invalid request format", err)
		return
	}

	// Require either email or username
	if credentials.Email == "" && credentials.Username == "" {
		logger.LogError(nil, "Login validation failed", logrus.Fields{
			"ip": c.ClientIP(),
		})
		response.BadRequest(c, "Email or username is required", nil)
		return
	}

	// Validate input
	var validationErrors validator.ValidationErrors
	if err := validator.ValidateRequired(credentials.Password, "password"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := validator.ValidateMinLength(credentials.Password, 6, "password"); err != nil {
		validationErrors = append(validationErrors, *err)
	}

	if len(validationErrors) > 0 {
		logger.LogError(nil, "Login validation failed", logrus.Fields{
			"email":    credentials.Email,
			"username": credentials.Username,
			"ip":       c.ClientIP(),
			"errors":   validationErrors,
		})
		response.ValidationError(c, validationErrors)
		return
	}

	// Sanitize inputs
	credentials.Email = validator.SanitizeString(credentials.Email)
	credentials.Username = validator.SanitizeString(credentials.Username)

	var user models.User
	var err error

	// Try to find user by email first, then by username
	if credentials.Email != "" {
		err = db.DB.Where("email = ?", credentials.Email).First(&user).Error
	} else {
		err = db.DB.Where("username = ?", credentials.Username).First(&user).Error
	}

	if err != nil {
		logger.LogError(err, "Login attempt with invalid credentials", logrus.Fields{
			"email":    credentials.Email,
			"username": credentials.Username,
			"ip":       c.ClientIP(),
		})
		response.Unauthorized(c, "Invalid credentials")
		return
	}

	if !utils.ComparePassword(user.Password, credentials.Password) {
		// Debug logging to help identify the issue
		logger.LogError(nil, "Password comparison failed", logrus.Fields{
			"email":             credentials.Email,
			"username":          credentials.Username,
			"user_id":           user.ID,
			"ip":                c.ClientIP(),
			"hashed_password":   user.Password,
			"provided_password": credentials.Password,
			"password_length":   len(credentials.Password),
			"hash_length":       len(user.Password),
		})

		// Additional test: try bcrypt comparison directly
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
		logger.LogError(err, "Direct bcrypt comparison result", logrus.Fields{
			"email":        credentials.Email,
			"bcrypt_error": err,
		})

		logger.LogError(nil, "Login attempt with invalid password", logrus.Fields{
			"email":    credentials.Email,
			"username": credentials.Username,
			"user_id":  user.ID,
			"ip":       c.ClientIP(),
		})
		response.Unauthorized(c, "Invalid credentials")
		return
	}

	// Generate Access Token
	accessToken, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		logger.LogError(err, "Failed to generate access token", logrus.Fields{
			"user_id":  user.ID,
			"username": user.Username,
		})
		response.InternalServerError(c, "Authentication failed", err)
		return
	}

	// Generate Refresh Token
	refreshToken, err := middleware.GenerateRefreshToken(user.ID)
	if err != nil {
		logger.LogError(err, "Failed to generate refresh token", logrus.Fields{
			"user_id":  user.ID,
			"username": user.Username,
		})
		response.InternalServerError(c, "Authentication failed", err)
		return
	}

	// Store refresh token in secure HttpOnly cookie
	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "localhost", false, true)

	logger.LogInfo("User logged in successfully", logrus.Fields{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"ip":       c.ClientIP(),
	})

	// Respond with access token and user info
	response.Success(c, gin.H{
		"access_token": accessToken,
		"user":         user.ToResponse(),
	}, "Login successful")
}
