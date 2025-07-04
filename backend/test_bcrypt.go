package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// The stored hash from the database
	storedHash := "$2a$12$qAp1oCkZ2xtg//aihZSqJeC0iWYIva3u/lmt05OWyW1cQK1.7o8gG"

	// Test different passwords
	passwords := []string{
		"password123",
		"Password123",
		"password123 ",
		" password123",
		"password123\n",
		"password123\r",
		"password123\r\n",
	}

	for _, password := range passwords {
		err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
		if err == nil {
			fmt.Printf("✓ Password '%s' (length: %d) MATCHES the hash\n", password, len(password))
		} else {
			fmt.Printf("✗ Password '%s' (length: %d) does NOT match the hash: %v\n", password, len(password), err)
		}
	}
}
