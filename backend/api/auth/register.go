package auth

import (
	"net/http"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/logger"
	"github.com/geoo115/property-manager/models"
	"github.com/geoo115/property-manager/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RegisterHandler(c *gin.Context) {
	var req models.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if req.Username == "" || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username, Password, and Email are required"})
		return
	}

	var userExists models.User
	if err := db.DB.Where("email = ?", req.Email).First(&userExists).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	if err := db.DB.Where("username = ?", req.Username).First(&userExists).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	// Debug logging
	logger.LogInfo("Password hashing successful", logrus.Fields{
		"email":             req.Email,
		"username":          req.Username,
		"original_password": req.Password,
		"hashed_password":   hashedPassword,
		"password_length":   len(req.Password),
	})

	user := models.User{
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  hashedPassword,
		Email:     req.Email,
		Role:      req.Role,
		Phone:     req.Phone,
		Avatar:    req.Avatar,
		IsActive:  true,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
