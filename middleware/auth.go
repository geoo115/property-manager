package middleware

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

var jwtSecret []byte

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
}

// GenerateToken creates a JWT with user information
func GenerateToken(userID uint, username string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   userID,
		"role":     role,
		"username": username,
		"exp":      time.Now().Add(1 * time.Hour).Unix(), // Token expires in 1 hour
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateRefreshToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(24 * time.Hour).Unix(), // Refresh token expires in 24 hours
	})
	return token.SignedString(jwtSecret)
}
