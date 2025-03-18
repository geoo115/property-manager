package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/geoo115/property-manager/api/auth"
	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/geoo115/property-manager/utils"
	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	// Initialize the database before running tests.
	db.Init() // Ensure this sets db.DB properly.

	// Run the tests.
	code := m.Run()

	// Exit with the proper code.
	os.Exit(code)
}

// getTestContext returns a Gin context and a ResponseRecorder for testing.
func getTestContext(method, url string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	return c, w
}

// TestRegisterHandler tests the user registration endpoint.
func TestRegisterHandler(t *testing.T) {
	// Set Gin to test mode.
	gin.SetMode(gin.TestMode)

	// Clean up any pre-existing test user.
	db.DB.Where("username = ?", "testuser2").Delete(&models.User{})
	db.DB.Where("email = ?", "testuser2@example.com").Delete(&models.User{})

	// Prepare the registration payload.
	payload := map[string]string{
		"username": "testuser2",
		"password": "password123",
		"email":    "testuser2@example.com",
		"role":     "tenant",
		"phone":    "1234567893",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Create test context.
	c, w := getTestContext("POST", "/register", body)

	// Call the RegisterHandler.
	auth.RegisterHandler(c)

	// Check the response code.
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d. Response: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	// Optionally, you can check the response body if a success message is returned.
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}
	if msg, exists := resp["message"]; !exists || msg == "" {
		t.Error("Expected success message in response")
	}
}

// TestLoginHandler tests the login endpoint.
func TestLoginHandler(t *testing.T) {
	// Set Gin to test mode.
	gin.SetMode(gin.TestMode)

	// Create a test user in the database.
	// First, delete any existing user with the same username.
	db.DB.Where("username = ?", "testlogin").Delete(&models.User{})
	hashedPassword, err := utils.HashPassword("password123")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	testUser := models.User{
		Username: "testlogin",
		Email:    "testlogin@example.com",
		Password: hashedPassword,
		Role:     "tenant",
		Phone:    "0987654321",
	}
	if err := db.DB.Create(&testUser).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Prepare the login payload.
	credentials := map[string]string{
		"username": "testlogin",
		"password": "password123",
	}
	jsonPayload, err := json.Marshal(credentials)
	if err != nil {
		t.Fatalf("Failed to marshal credentials: %v", err)
	}

	// Create test context.
	c, w := getTestContext("POST", "/login", jsonPayload)

	// Call the LoginHandler.
	auth.LoginHandler(c)

	// Check the response code.
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Response: %s", http.StatusOK, w.Code, w.Body.String())
	}

	// Check if the token is present in the response.
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}
	if token, exists := resp["token"]; !exists || token == "" {
		t.Error("Expected token in response")
	}
}
