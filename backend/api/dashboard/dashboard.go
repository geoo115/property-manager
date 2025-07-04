package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/logger"
	"github.com/geoo115/property-manager/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// DashboardStats represents dashboard statistics
type DashboardStats struct {
	TotalProperties     int64   `json:"totalProperties"`
	ActiveLeases        int64   `json:"activeLeases"`
	PendingMaintenance  int64   `json:"pendingMaintenance"`
	TotalRevenue        float64 `json:"totalRevenue"`
	MyProperties        int64   `json:"myProperties"`
	MonthlyRevenue      float64 `json:"monthlyRevenue"`
	MaintenanceRequests int64   `json:"maintenanceRequests"`
	MyLeases            int64   `json:"myLeases"`
	OutstandingInvoices int64   `json:"outstandingInvoices"`
	TotalPaid           float64 `json:"totalPaid"`
	// Additional admin metrics
	TotalUsers           int64   `json:"totalUsers,omitempty"`
	TotalTenants         int64   `json:"totalTenants,omitempty"`
	TotalLandlords       int64   `json:"totalLandlords,omitempty"`
	ExpiredLeases        int64   `json:"expiredLeases,omitempty"`
	OverdueInvoices      int64   `json:"overdueInvoices,omitempty"`
	CompletedMaintenance int64   `json:"completedMaintenance,omitempty"`
	TotalExpenses        float64 `json:"totalExpenses,omitempty"`
	MonthlyExpenses      float64 `json:"monthlyExpenses,omitempty"`
	OccupancyRate        float64 `json:"occupancyRate,omitempty"`
}

// Activity represents a recent activity
type Activity struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Timestamp   time.Time `json:"timestamp"`
	Type        string    `json:"type"`
	EntityID    uint      `json:"entityId"`
	EntityType  string    `json:"entityType"`
}

// GetDashboardStats returns dashboard statistics based on user role
func GetDashboardStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRole, _ := c.Get("user_role")
	ctx := context.Background()
	cacheKey := "dashboard_stats:" + userRole.(string) + ":" + strconv.FormatUint(uint64(userID.(uint)), 10)

	// Check Redis cache
	cachedData, err := db.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var stats DashboardStats
		if json.Unmarshal([]byte(cachedData), &stats) == nil {
			c.JSON(http.StatusOK, stats)
			return
		}
	}

	stats := DashboardStats{}

	switch userRole {
	case "admin":
		// Admin sees comprehensive system statistics
		db.DB.Model(&models.Property{}).Count(&stats.TotalProperties)
		db.DB.Model(&models.Lease{}).Where("status = ?", "active").Count(&stats.ActiveLeases)
		db.DB.Model(&models.Maintenance{}).Where("status = ?", "pending").Count(&stats.PendingMaintenance)

		// Calculate total revenue from invoices
		var totalRevenue float64
		db.DB.Model(&models.Invoice{}).Select("COALESCE(SUM(paid_amount), 0)").Scan(&totalRevenue)
		stats.TotalRevenue = totalRevenue

		// Additional admin-specific metrics
		db.DB.Model(&models.User{}).Count(&stats.TotalUsers)
		db.DB.Model(&models.User{}).Where("role = ?", "tenant").Count(&stats.TotalTenants)
		db.DB.Model(&models.User{}).Where("role = ?", "landlord").Count(&stats.TotalLandlords)
		db.DB.Model(&models.Lease{}).Where("status = ? OR end_date < ?", "expired", time.Now()).Count(&stats.ExpiredLeases)
		db.DB.Model(&models.Invoice{}).Where("payment_status = ? AND due_date < ?", "pending", time.Now()).Count(&stats.OverdueInvoices)
		db.DB.Model(&models.Maintenance{}).Where("status = ?", "completed").Count(&stats.CompletedMaintenance)

		// Calculate total expenses
		var totalExpenses float64
		db.DB.Model(&models.Expense{}).Select("COALESCE(SUM(amount), 0)").Scan(&totalExpenses)
		stats.TotalExpenses = totalExpenses

		// Calculate monthly expenses (current month)
		var monthlyExpenses float64
		currentMonth := time.Now().Format("2006-01")
		db.DB.Model(&models.Expense{}).
			Where("TO_CHAR(expense_date, 'YYYY-MM') = ?", currentMonth).
			Select("COALESCE(SUM(amount), 0)").Scan(&monthlyExpenses)
		stats.MonthlyExpenses = monthlyExpenses

		// Calculate occupancy rate
		var totalOccupiableProperties int64
		var occupiedProperties int64
		db.DB.Model(&models.Property{}).Where("status = ?", "available").Count(&totalOccupiableProperties)
		db.DB.Model(&models.Property{}).
			Joins("JOIN leases ON properties.id = leases.property_id").
			Where("leases.status = ?", "active").
			Count(&occupiedProperties)

		if totalOccupiableProperties > 0 {
			stats.OccupancyRate = float64(occupiedProperties) / float64(totalOccupiableProperties) * 100
		}

	case "landlord":
		// Landlord sees their own statistics
		db.DB.Model(&models.Property{}).Where("owner_id = ?", userID).Count(&stats.MyProperties)
		db.DB.Model(&models.Lease{}).Joins("JOIN properties ON properties.id = leases.property_id").
			Where("properties.owner_id = ? AND leases.status = ?", userID, "active").Count(&stats.ActiveLeases)
		db.DB.Model(&models.Maintenance{}).Joins("JOIN properties ON properties.id = maintenance_requests.property_id").
			Where("properties.owner_id = ? AND maintenance_requests.status = ?", userID, "pending").Count(&stats.MaintenanceRequests)

		// Calculate monthly revenue (PostgreSQL compatible)
		var monthlyRevenue float64
		currentMonth := time.Now().Format("2006-01")
		db.DB.Model(&models.Invoice{}).Joins("JOIN properties ON properties.id = invoices.property_id").
			Where("properties.owner_id = ? AND TO_CHAR(invoices.invoice_date, 'YYYY-MM') = ?", userID, currentMonth).
			Select("COALESCE(SUM(paid_amount), 0)").Scan(&monthlyRevenue)
		stats.MonthlyRevenue = monthlyRevenue

	case "tenant":
		// Tenant sees their own statistics
		db.DB.Model(&models.Lease{}).Where("tenant_id = ? AND status = ?", userID, "active").Count(&stats.MyLeases)
		db.DB.Model(&models.Invoice{}).Where("tenant_id = ? AND payment_status = ?", userID, "unpaid").Count(&stats.OutstandingInvoices)
		db.DB.Model(&models.Maintenance{}).Where("requested_by_id = ?", userID).Count(&stats.MaintenanceRequests)

		// Calculate total paid
		var totalPaid float64
		db.DB.Model(&models.Invoice{}).Where("tenant_id = ?", userID).
			Select("COALESCE(SUM(paid_amount), 0)").Scan(&totalPaid)
		stats.TotalPaid = totalPaid
	}

	// Cache the result for 5 minutes
	jsonData, _ := json.Marshal(stats)
	db.RedisClient.Set(ctx, cacheKey, jsonData, 5*time.Minute)

	c.JSON(http.StatusOK, stats)
}

// GetRecentActivities returns recent activities for the dashboard
func GetRecentActivities(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRole, _ := c.Get("user_role")
	ctx := context.Background()
	cacheKey := "dashboard_activities:" + userRole.(string) + ":" + strconv.FormatUint(uint64(userID.(uint)), 10)

	// Check Redis cache
	cachedData, err := db.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var activities []Activity
		if json.Unmarshal([]byte(cachedData), &activities) == nil {
			c.JSON(http.StatusOK, gin.H{"activities": activities})
			return
		}
	}

	activities := []Activity{}

	switch userRole {
	case "admin":
		// Admin sees all recent activities
		activities = getAdminActivities()
	case "landlord":
		// Landlord sees activities related to their properties
		activities = getLandlordActivities(userID.(uint))
	case "tenant":
		// Tenant sees their own activities
		activities = getTenantActivities(userID.(uint))
	case "maintenanceTeam":
		// Maintenance team sees maintenance-related activities
		activities = getMaintenanceTeamActivities()
	}

	// Cache the result for 2 minutes
	jsonData, _ := json.Marshal(activities)
	db.RedisClient.Set(ctx, cacheKey, jsonData, 2*time.Minute)

	c.JSON(http.StatusOK, gin.H{"activities": activities})
}

// getAdminActivities returns recent activities for admin users
func getAdminActivities() []Activity {
	activities := []Activity{}

	// Recent property additions
	var properties []models.Property
	db.DB.Order("created_at DESC").Limit(3).Find(&properties)
	for _, prop := range properties {
		activities = append(activities, Activity{
			ID:          prop.ID,
			Title:       "New Property Added",
			Description: prop.Name + " in " + prop.City,
			Icon:        "home",
			Timestamp:   prop.CreatedAt,
			Type:        "property",
			EntityID:    prop.ID,
			EntityType:  "property",
		})
	}

	// Recent lease activities
	var leases []models.Lease
	db.DB.Preload("Property").Preload("Tenant").Order("created_at DESC").Limit(3).Find(&leases)
	for _, lease := range leases {
		activities = append(activities, Activity{
			ID:          lease.ID,
			Title:       "New Lease Created",
			Description: lease.Property.Name + " - " + lease.Tenant.FirstName + " " + lease.Tenant.LastName,
			Icon:        "file-contract",
			Timestamp:   lease.CreatedAt,
			Type:        "lease",
			EntityID:    lease.ID,
			EntityType:  "lease",
		})
	}

	// Recent maintenance requests
	var maintenance []models.Maintenance
	db.DB.Preload("Property").Preload("RequestedBy").Order("created_at DESC").Limit(3).Find(&maintenance)
	for _, maint := range maintenance {
		activities = append(activities, Activity{
			ID:          maint.ID,
			Title:       "Maintenance Request",
			Description: maint.Property.Name + " - " + maint.Description,
			Icon:        "tools",
			Timestamp:   maint.CreatedAt,
			Type:        "maintenance",
			EntityID:    maint.ID,
			EntityType:  "maintenance",
		})
	}

	// Recent user registrations
	var users []models.User
	db.DB.Where("role IN ?", []string{"tenant", "landlord"}).Order("created_at DESC").Limit(3).Find(&users)
	for _, user := range users {
		activities = append(activities, Activity{
			ID:          user.ID,
			Title:       "New User Registered",
			Description: fmt.Sprintf("%s %s (%s)", user.FirstName, user.LastName, user.Role),
			Icon:        "user-plus",
			Timestamp:   user.CreatedAt,
			Type:        "user",
			EntityID:    user.ID,
			EntityType:  "user",
		})
	}

	// Recent invoices
	var invoices []models.Invoice
	db.DB.Preload("Property").Preload("Tenant").Order("created_at DESC").Limit(3).Find(&invoices)
	for _, invoice := range invoices {
		activities = append(activities, Activity{
			ID:          invoice.ID,
			Title:       "Invoice Generated",
			Description: fmt.Sprintf("%s - %s %s - $%.2f", invoice.Property.Name, invoice.Tenant.FirstName, invoice.Tenant.LastName, invoice.Amount),
			Icon:        "file-invoice",
			Timestamp:   invoice.CreatedAt,
			Type:        "invoice",
			EntityID:    invoice.ID,
			EntityType:  "invoice",
		})
	}

	// Sort all activities by timestamp and limit to 15 most recent
	sortActivitiesByTimestamp(activities)
	if len(activities) > 15 {
		activities = activities[:15]
	}

	return activities
}

// getLandlordActivities returns recent activities for landlord users
func getLandlordActivities(userID uint) []Activity {
	activities := []Activity{}

	// Recent property activities
	var properties []models.Property
	db.DB.Where("owner_id = ?", userID).Order("created_at DESC").Limit(5).Find(&properties)
	for _, prop := range properties {
		activities = append(activities, Activity{
			ID:          prop.ID,
			Title:       "Property Updated",
			Description: prop.Name + " in " + prop.City,
			Icon:        "home",
			Timestamp:   prop.UpdatedAt,
			Type:        "property",
			EntityID:    prop.ID,
			EntityType:  "property",
		})
	}

	// Recent lease activities for landlord's properties
	var leases []models.Lease
	db.DB.Preload("Property").Preload("Tenant").
		Joins("JOIN properties ON properties.id = leases.property_id").
		Where("properties.owner_id = ?", userID).
		Order("leases.created_at DESC").Limit(5).Find(&leases)
	for _, lease := range leases {
		activities = append(activities, Activity{
			ID:          lease.ID,
			Title:       "Lease Activity",
			Description: lease.Property.Name + " - " + lease.Tenant.FirstName + " " + lease.Tenant.LastName,
			Icon:        "file-contract",
			Timestamp:   lease.CreatedAt,
			Type:        "lease",
			EntityID:    lease.ID,
			EntityType:  "lease",
		})
	}

	return activities
}

// getTenantActivities returns recent activities for tenant users
func getTenantActivities(userID uint) []Activity {
	activities := []Activity{}

	// Recent lease activities
	var leases []models.Lease
	db.DB.Preload("Property").Where("tenant_id = ?", userID).Order("created_at DESC").Limit(5).Find(&leases)
	for _, lease := range leases {
		activities = append(activities, Activity{
			ID:          lease.ID,
			Title:       "Lease Information",
			Description: lease.Property.Name + " - " + lease.Property.Address,
			Icon:        "file-contract",
			Timestamp:   lease.UpdatedAt,
			Type:        "lease",
			EntityID:    lease.ID,
			EntityType:  "lease",
		})
	}

	// Recent maintenance requests
	var maintenance []models.Maintenance
	db.DB.Preload("Property").Where("requested_by_id = ?", userID).Order("created_at DESC").Limit(5).Find(&maintenance)
	for _, maint := range maintenance {
		activities = append(activities, Activity{
			ID:          maint.ID,
			Title:       "Maintenance Request",
			Description: maint.Property.Name + " - " + maint.Description,
			Icon:        "tools",
			Timestamp:   maint.CreatedAt,
			Type:        "maintenance",
			EntityID:    maint.ID,
			EntityType:  "maintenance",
		})
	}

	// Recent invoice activities
	var invoices []models.Invoice
	db.DB.Preload("Property").Where("tenant_id = ?", userID).Order("created_at DESC").Limit(5).Find(&invoices)
	for _, invoice := range invoices {
		activities = append(activities, Activity{
			ID:          invoice.ID,
			Title:       "Invoice Generated",
			Description: fmt.Sprintf("%s - $%.2f", invoice.Property.Name, invoice.Amount),
			Icon:        "file-invoice",
			Timestamp:   invoice.CreatedAt,
			Type:        "invoice",
			EntityID:    invoice.ID,
			EntityType:  "invoice",
		})
	}

	return activities
}

// getMaintenanceTeamActivities returns recent activities for maintenance team
func getMaintenanceTeamActivities() []Activity {
	activities := []Activity{}

	// Recent maintenance requests
	var maintenance []models.Maintenance
	db.DB.Preload("Property").Preload("RequestedBy").Order("created_at DESC").Limit(10).Find(&maintenance)
	for _, maint := range maintenance {
		activities = append(activities, Activity{
			ID:          maint.ID,
			Title:       "Maintenance Request",
			Description: maint.Property.Name + " - " + maint.Description,
			Icon:        "tools",
			Timestamp:   maint.CreatedAt,
			Type:        "maintenance",
			EntityID:    maint.ID,
			EntityType:  "maintenance",
		})
	}

	return activities
}

// InvalidateDashboardCache clears dashboard cache for a user
func InvalidateDashboardCache(userID uint, userRole string) {
	ctx := context.Background()
	userIDStr := strconv.FormatUint(uint64(userID), 10)
	statsKey := "dashboard_stats:" + userRole + ":" + userIDStr
	activitiesKey := "dashboard_activities:" + userRole + ":" + userIDStr

	db.RedisClient.Del(ctx, statsKey, activitiesKey)

	logger.LogInfo("Dashboard cache invalidated", logrus.Fields{
		"user_id": userID,
		"role":    userRole,
	})
}

// sortActivitiesByTimestamp sorts activities by timestamp in descending order
func sortActivitiesByTimestamp(activities []Activity) {
	sort.Slice(activities, func(i, j int) bool {
		return activities[i].Timestamp.After(activities[j].Timestamp)
	})
}
