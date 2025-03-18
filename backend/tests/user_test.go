package tests

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/geoo115/property-manager/api/user"
	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/geoo115/property-manager/utils"
	"github.com/gin-gonic/gin"
)

// randomPhone generates a random 10-digit phone number as a string.
func randomPhone() string {
	rand.Seed(time.Now().UnixNano())
	// Ensure a 10-digit number (leading digit non-zero).
	return strconv.Itoa(rand.Intn(9000000000) + 1000000000)
}

// randomUsername generates a unique username.
func randomUsername(base string) string {
	rand.Seed(time.Now().UnixNano())
	return base + strconv.Itoa(rand.Intn(100000))
}

// randomEmail generates a unique email.
func randomEmail(base string) string {
	rand.Seed(time.Now().UnixNano())
	return base + strconv.Itoa(rand.Intn(100000)) + "@example.com"
}

// TestCreateUser tests the POST /users endpoint.
func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Step 1: Ensure unique email & phone to avoid constraint violations
	username := randomUsername("testuser")
	email := randomEmail(username)
	phone := randomPhone()

	// Step 2: Prepare request payload
	payload := map[string]interface{}{
		"username":   username,
		"password":   "password123",
		"first_name": "Test",
		"last_name":  "User",
		"email":      email,
		"role":       "tenant",
		"phone":      phone,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Step 3: Create HTTP request (better than getTestContext)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Step 4: Setup router and execute the request
	router := gin.Default()
	router.POST("/users", user.CreateUser) // Attach the route handler
	router.ServeHTTP(w, req)               // Simulate the HTTP request

	// Step 5: Verify response
	if w.Code != http.StatusCreated { // Expect 201 Created
		t.Fatalf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Step 6: Verify JSON response contains user data
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Step 7: Extract user object from response
	userData, exists := response["user"].(map[string]interface{})
	if !exists {
		t.Fatalf("Response does not contain user data")
	}

	// Step 8: Validate response values
	if userData["username"] != payload["username"] {
		t.Errorf("Expected username %s, got %s", payload["username"], userData["username"])
	}

	// Step 9: Verify that the user was actually inserted in the database
	var createdUser models.User
	if err := db.DB.Where("username = ?", username).First(&createdUser).Error; err != nil {
		t.Fatalf("User not found in database: %v", err)
	}
}

// TestGetUsers tests the GET /users endpoint.
func TestGetUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Ensure at least one user exists in the DB.
	var count int64
	db.DB.Model(&models.User{}).Count(&count)
	if count == 0 {
		// Create a dummy user if none exists.
		dummyUser := models.User{
			Username:  "dummyuser",
			FirstName: "Dummy",
			LastName:  "User",
			Password:  "dummy123",
			Email:     "dummy@example.com",
			Role:      "tenant",
			Phone:     "1234567890",
		}
		hashed, _ := utils.HashPassword(dummyUser.Password)
		dummyUser.Password = hashed
		if err := db.DB.Create(&dummyUser).Error; err != nil {
			t.Fatalf("Error creating dummy user: %v", err)
		}
	}

	// Create test context and call the handler.
	c, w := getTestContext("GET", "/users", nil)
	user.GetUsers(c)

	// Check that we got an HTTP 200 OK status.
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Unmarshal the response body into a map.
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Error unmarshalling response: %v", err)
	}

	// Verify that the "users" key exists and is an array.
	users, ok := response["users"].([]interface{})
	if !ok {
		t.Fatalf("Expected 'users' key in response to be an array, got: %v", response["users"])
	}

	if len(users) == 0 {
		t.Error("Expected at least one user, got none")
	}

	// Optionally verify additional keys like "count", "limit", and "offset".
	if _, exists := response["count"]; !exists {
		t.Error("Expected 'count' key in response")
	}
	if _, exists := response["limit"]; !exists {
		t.Error("Expected 'limit' key in response")
	}
	if _, exists := response["offset"]; !exists {
		t.Error("Expected 'offset' key in response")
	}
}

func setParams(c *gin.Context, key, value string) {
	c.Params = []gin.Param{{Key: key, Value: value}}
}

func TestGetUserByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Step 1: Create a test user in the database.
	username := randomUsername("testuser")
	email := randomEmail(username)
	phone := randomPhone()

	testUser := models.User{
		Username:  username,
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Email:     email,
		Role:      "tenant",
		Phone:     phone,
	}

	// Hash the password before saving.
	hashed, err := utils.HashPassword(testUser.Password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}
	testUser.Password = hashed

	// Insert into the database.
	if err := db.DB.Create(&testUser).Error; err != nil {
		t.Fatalf("Error creating test user: %v", err)
	}

	// Step 2: Use the actual ID of the created user in the request.
	url := "/users/" + strconv.Itoa(int(testUser.ID))
	c, w := getTestContext("GET", url, nil)
	setParams(c, "id", strconv.Itoa(int(testUser.ID))) // Ensure ID is set in the context.
	user.GetUserByID(c)

	// Step 3: Validate response.
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Step 4: Deserialize response and validate the returned user.
	var fetchedUser models.User
	if err := json.Unmarshal(w.Body.Bytes(), &fetchedUser); err != nil {
		t.Fatalf("Error unmarshalling response: %v", err)
	}
	if fetchedUser.ID != testUser.ID {
		t.Errorf("Expected user ID %d, got %d", testUser.ID, fetchedUser.ID)
	}
}

func TestUpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Generate unique username, email, and phone.
	username := randomUsername("updateuser")
	email := randomEmail(username)
	phone := randomPhone()

	testUser := models.User{
		Username:  username,
		Password:  "password123",
		FirstName: "Update",
		LastName:  "User",
		Email:     email,
		Role:      "tenant",
		Phone:     phone,
	}

	hashed, err := utils.HashPassword(testUser.Password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}
	testUser.Password = hashed

	if err := db.DB.Create(&testUser).Error; err != nil {
		t.Fatalf("Error creating test user: %v", err)
	}

	// Prepare update payload.
	updatePayload := map[string]interface{}{
		"username":   randomUsername("updateduser"), // Unique username
		"password":   "newpassword123",
		"first_name": "Updated",
		"last_name":  "User",
		"email":      randomEmail("updateduser"), // Unique email
		"role":       "tenant",
		"phone":      testUser.Phone, // Keep the same phone
	}
	payloadBytes, err := json.Marshal(updatePayload)
	if err != nil {
		t.Fatalf("Error marshalling update payload: %v", err)
	}

	url := "/users/" + strconv.Itoa(int(testUser.ID))
	c, w := getTestContext("PUT", url, payloadBytes)

	// ✅ FIX: Manually set the "id" parameter before calling UpdateUser.
	setParams(c, "id", strconv.Itoa(int(testUser.ID)))

	user.UpdateUser(c)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var updatedUser models.User
	if err := json.Unmarshal(w.Body.Bytes(), &updatedUser); err != nil {
		t.Fatalf("Error unmarshalling response: %v", err)
	}
	if updatedUser.Username != updatePayload["username"] {
		t.Errorf("Expected username '%s', got '%s'", updatePayload["username"], updatedUser.Username)
	}
	if updatedUser.Password == "newpassword123" {
		t.Error("Expected password to be hashed, but it remains in plain text")
	}
}

func TestDeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Step 1: Create a test user in the database.
	username := randomUsername("testuser")
	email := randomEmail(username)
	phone := randomPhone()

	testUser := models.User{
		Username:  username,
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Email:     email,
		Role:      "tenant",
		Phone:     phone,
	}

	// Hash the password before saving.
	hashed, err := utils.HashPassword(testUser.Password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}
	testUser.Password = hashed

	// Insert into the database.
	if err := db.DB.Create(&testUser).Error; err != nil {
		t.Fatalf("Error creating test user: %v", err)
	}

	// Step 2: Use the actual ID of the created user in the request.
	url := "/users/" + strconv.Itoa(int(testUser.ID))
	c, w := getTestContext("DELETE", url, nil)
	setParams(c, "id", strconv.Itoa(int(testUser.ID))) // Ensure ID is set in the context.
	user.DeleteUser(c)

	// Step 3: Validate response.
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Step 4: Attempt to fetch the user again (should return 404).
	c, w = getTestContext("GET", url, nil)
	setParams(c, "id", strconv.Itoa(int(testUser.ID))) // Ensure ID is set in the context.
	user.GetUserByID(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Expected status 404, got %d. Body: %s", w.Code, w.Body.String())
	}
}
