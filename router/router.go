package router

import (
	"time"

	"github.com/geoo115/property-manager/api/accounting"
	"github.com/geoo115/property-manager/api/lease"
	"github.com/geoo115/property-manager/api/maintenance"
	"github.com/geoo115/property-manager/api/property"
	"github.com/geoo115/property-manager/api/user"
	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	// Public routes
	AuthRoutes(r)

	// Admin group: full access to all endpoints.
	admin := r.Group("/admin")
	admin.Use(db.RateLimit(100, time.Minute), middleware.JWTMiddleware(), middleware.RoleMiddleware("admin"))
	{
		UserRouter(admin)
		PropertyRouter(admin)
		LeaseRouter(admin)
		MaintenanceRoutes(admin)
		// Mount accounting endpoints under "/admin/accounting"
		accountingGroup := admin.Group("/accounting")
		AccountingRouter(accountingGroup)
	}

	// Landlord group: restricted access to their properties and leases.
	landlord := r.Group("/landlord")
	landlord.Use(db.RateLimit(100, time.Minute), middleware.JWTMiddleware(), middleware.RoleMiddleware("landlord"))
	{
		landlord.GET("/properties", property.GetProperties)
		landlord.GET("/properties/:id", property.GetPropertyByID)
		landlord.GET("/leases", lease.GetLeases)
		landlord.GET("/leases/:id", lease.GetLeaseByID)
		landlord.GET("/properties/:id/maintenances", maintenance.GetLandlordMaintenances)
		landlord.POST("/properties/:id/maintenances", maintenance.CreateMaintenanceByProperty)
		landlord.GET("/invoices", accounting.GetInvoicesForLandlord)
		landlord.GET("/expenses", accounting.GetExpensesForLandlord)
	}

	// Tenant group: can only access their leases.
	tenant := r.Group("/tenant")
	tenant.Use(db.RateLimit(100, time.Minute), middleware.JWTMiddleware(), middleware.RoleMiddleware("tenant"))
	{
		tenant.GET("/leases", lease.GetLeasesForTenant)
		tenant.GET("/leases/:id", lease.GetLeaseByID)
		tenant.GET("/leases/active", lease.GetActiveLeaseForTenant) // Ensure this route is defined
		tenant.GET("/leases/:id/maintenance", maintenance.GetMaintenances)
		tenant.POST("/leases/:id/maintenance", maintenance.CreateMaintenanceByLease)
		tenant.GET("/invoices", accounting.GetInvoicesForTenant)
	}

	// MaintenanceTeam group: for maintenance staff
	maintenanceTeam := r.Group("/maintenanceTeam")
	maintenanceTeam.Use(db.RateLimit(100, time.Minute), middleware.JWTMiddleware(), middleware.RoleMiddleware("maintenanceTeam"))
	{
		maintenanceTeam.GET("/maintenances", maintenance.GetMaintenances)
		maintenanceTeam.GET("/maintenance/:id", maintenance.GetMaintenance)
		maintenanceTeam.PUT("/maintenance/:id", maintenance.UpdateMaintenance)
		maintenanceTeam.GET("/users", user.GetUsers)
		maintenanceTeam.GET("/properties", property.GetProperties)
	}
}
