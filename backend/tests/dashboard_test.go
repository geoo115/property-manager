package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/geoo115/property-manager/api/dashboard"
	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/geoo115/property-manager/utils"
	"github.com/gin-gonic/gin"
)

func TestGetDashboardStats(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test users
	adminUser := createTestUser("admin", "admin@test.com", "admin")
	landlordUser := createTestUser("landlord", "landlord@test.com", "landlord")
	tenantUser := createTestUser("tenant", "tenant@test.com", "tenant")

	// Create test properties
	property1 := models.Property{
		Name:      "Test Property 1",
		Address:   "123 Test St",
		City:      "Test City",
		Price:     1500.00,
		OwnerID:   landlordUser.ID,
		Available: true,
	}
	db.DB.Create(&property1)

	property2 := models.Property{
		Name:      "Test Property 2",
		Address:   "456 Test Ave",
		City:      "Test City",
		Price:     2000.00,
		OwnerID:   landlordUser.ID,
		Available: true,
	}
	db.DB.Create(&property2)

	// Create test lease
	lease := models.Lease{
		PropertyID:  property1.ID,
		TenantID:    tenantUser.ID,
		StartDate:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		MonthlyRent: 1500.00,
		Status:      "active",
	}
	db.DB.Create(&lease)

	// Create test maintenance request
	maintenance := models.Maintenance{
		PropertyID:    property1.ID,
		RequestedByID: tenantUser.ID,
		Description:   "Test maintenance request",
		Status:        "pending",
		Priority:      "medium",
	}
	db.DB.Create(&maintenance)

	// Create test invoice
	invoice := models.Invoice{
		TenantID:      tenantUser.ID,
		PropertyID:    property1.ID,
		Amount:        1500.00,
		PaidAmount:    1500.00,
		InvoiceDate:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		DueDate:       time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
		Category:      "rent",
		PaymentStatus: "paid",
	}
	db.DB.Create(&invoice)

	tests := []struct {
		name           string
		userID         uint
		userRole       string
		expectedStatus int
		checkFields    []string
	}{
		{
			name:           "Admin Dashboard Stats",
			userID:         adminUser.ID,
			userRole:       "admin",
			expectedStatus: http.StatusOK,
			checkFields:    []string{"totalProperties", "activeLeases", "pendingMaintenance"},
		},
		{
			name:           "Landlord Dashboard Stats",
			userID:         landlordUser.ID,
			userRole:       "landlord",
			expectedStatus: http.StatusOK,
			checkFields:    []string{"myProperties", "activeLeases", "maintenanceRequests"},
		},
		{
			name:           "Tenant Dashboard Stats",
			userID:         tenantUser.ID,
			userRole:       "tenant",
			expectedStatus: http.StatusOK,
			checkFields:    []string{"myLeases", "outstandingInvoices", "maintenanceRequests"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set user context
			c.Set("user_id", tt.userID)
			c.Set("user_role", tt.userRole)

			// Call the handler
			dashboard.GetDashboardStats(c)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Parse response
			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			// Check required fields
			for _, field := range tt.checkFields {
				if _, exists := response[field]; !exists {
					t.Errorf("Expected field '%s' in response", field)
				}
			}
		})
	}
}

func TestGetRecentActivities(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test users
	adminUser := createTestUser("admin", "admin@test.com", "admin")
	landlordUser := createTestUser("landlord", "landlord@test.com", "landlord")
	tenantUser := createTestUser("tenant", "tenant@test.com", "tenant")

	// Create test property
	property := models.Property{
		Name:      "Test Property",
		Address:   "123 Test St",
		City:      "Test City",
		Price:     1500.00,
		OwnerID:   landlordUser.ID,
		Available: true,
	}
	db.DB.Create(&property)

	// Create test lease
	lease := models.Lease{
		PropertyID:  property.ID,
		TenantID:    tenantUser.ID,
		StartDate:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		MonthlyRent: 1500.00,
		Status:      "active",
	}
	db.DB.Create(&lease)

	// Create test maintenance request
	maintenance := models.Maintenance{
		PropertyID:    property.ID,
		RequestedByID: tenantUser.ID,
		Description:   "Test maintenance request",
		Status:        "pending",
		Priority:      "medium",
	}
	db.DB.Create(&maintenance)

	tests := []struct {
		name           string
		userID         uint
		userRole       string
		expectedStatus int
	}{
		{
			name:           "Admin Activities",
			userID:         adminUser.ID,
			userRole:       "admin",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Landlord Activities",
			userID:         landlordUser.ID,
			userRole:       "landlord",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Tenant Activities",
			userID:         tenantUser.ID,
			userRole:       "tenant",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set user context
			c.Set("user_id", tt.userID)
			c.Set("user_role", tt.userRole)

			// Call the handler
			dashboard.GetRecentActivities(c)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Parse response
			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			// Check activities field exists
			if _, exists := response["activities"]; !exists {
				t.Error("Expected 'activities' field in response")
			}
		})
	}
}

func TestDashboardUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		handler        func(*gin.Context)
		expectedStatus int
	}{
		{
			name:           "GetDashboardStats Unauthorized",
			handler:        dashboard.GetDashboardStats,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "GetRecentActivities Unauthorized",
			handler:        dashboard.GetRecentActivities,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test context without user context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Call the handler
			tt.handler(c)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Parse response
			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			// Check error message
			if _, exists := response["error"]; !exists {
				t.Error("Expected 'error' field in response")
			}
		})
	}
}

// Helper function to create test users
func createTestUser(username, email, role string) models.User {
	hashedPassword, _ := utils.HashPassword("password123")
	user := models.User{
		Username:  username,
		FirstName: "Test",
		LastName:  "User",
		Email:     email,
		Password:  hashedPassword,
		Role:      role,
		Phone:     "1234567890",
		IsActive:  true,
	}
	db.DB.Create(&user)
	return user
}

// Helper function to clean up test data
func cleanupTestData() {
	db.DB.Exec("DELETE FROM maintenance_requests")
	db.DB.Exec("DELETE FROM invoices")
	db.DB.Exec("DELETE FROM leases")
	db.DB.Exec("DELETE FROM properties")
	db.DB.Exec("DELETE FROM users")
}
