package router

import (
	"net/http"
	"time"

	"github.com/geoo115/property-manager/api/accounting"
	"github.com/geoo115/property-manager/api/lease"
	"github.com/geoo115/property-manager/api/maintenance"
	"github.com/geoo115/property-manager/api/property"
	"github.com/geoo115/property-manager/api/user"
	"github.com/geoo115/property-manager/config"
	"github.com/geoo115/property-manager/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, cfg *config.Config) {
	// Add global middleware
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CORSFromConfig(cfg))

	// Handle 404 and 405 errors
	r.NoRoute(middleware.NotFoundHandler())
	r.NoMethod(middleware.MethodNotAllowedHandler())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "property-manager",
			"version": "1.0.0",
		})
	})

	// Public routes with development-friendly rate limiting
	public := r.Group("/api/v1")
	public.Use(middleware.IPRateLimit(500, time.Minute)) // 500 requests per minute per IP for development
	{
		AuthRoutes(public)
	}

	// Admin group: full access to all endpoints
	admin := r.Group("/api/v1/admin")
	admin.Use(
		middleware.IPRateLimit(200, time.Minute), // Higher limit for authenticated users
		middleware.JWTMiddleware(),
		middleware.RoleMiddleware("admin"),
	)
	{
		UserRouter(admin)
		PropertyRouter(admin)
		LeaseRouter(admin)
		MaintenanceRoutes(admin)
		// Mount accounting endpoints under "/admin/accounting"
		accountingGroup := admin.Group("/accounting")
		AccountingRouter(accountingGroup)
		// Mount dashboard endpoints
		DashboardRouter(admin)
	}

	// Landlord group: restricted access to their properties and leases
	landlord := r.Group("/api/v1/landlord")
	landlord.Use(
		middleware.IPRateLimit(150, time.Minute),
		middleware.JWTMiddleware(),
		middleware.RoleMiddleware("landlord"),
	)
	{
		landlord.GET("/properties", property.GetProperties)
		landlord.GET("/properties/:id", property.GetPropertyByID)
		landlord.GET("/leases", lease.GetLeases)
		landlord.GET("/leases/:id", lease.GetLeaseByID)
		landlord.GET("/properties/:id/maintenances", maintenance.GetLandlordMaintenances)
		landlord.POST("/properties/:id/maintenances", maintenance.CreateMaintenanceByProperty)
		landlord.GET("/invoices", accounting.GetInvoicesForLandlord)
		landlord.GET("/expenses", accounting.GetExpensesForLandlord)
		// Mount dashboard endpoints for landlords
		DashboardRouter(landlord)
	}

	// Tenant group: can only access their leases.
	tenant := r.Group("/tenant")
	tenant.Use(middleware.IPRateLimit(100, time.Minute), middleware.JWTMiddleware(), middleware.RoleMiddleware("tenant"))
	{
		tenant.GET("/leases", lease.GetLeasesForTenant)
		tenant.GET("/leases/:id", lease.GetLeaseByID)
		tenant.GET("/leases/active", lease.GetActiveLeaseForTenant) // Ensure this route is defined
		tenant.GET("/leases/:id/maintenance", maintenance.GetMaintenances)
		tenant.POST("/leases/:id/maintenance", maintenance.CreateMaintenanceByLease)
		tenant.GET("/invoices", accounting.GetInvoicesForTenant)
		// Mount dashboard endpoints for tenants
		DashboardRouter(tenant)
	}

	// MaintenanceTeam group: for maintenance staff
	maintenanceTeam := r.Group("/maintenanceTeam")
	maintenanceTeam.Use(middleware.IPRateLimit(100, time.Minute), middleware.JWTMiddleware(), middleware.RoleMiddleware("maintenanceTeam"))
	{
		maintenanceTeam.GET("/maintenances", maintenance.GetMaintenances)
		maintenanceTeam.GET("/maintenance/:id", maintenance.GetMaintenance)
		maintenanceTeam.PUT("/maintenance/:id", maintenance.UpdateMaintenance)
		maintenanceTeam.GET("/users", user.GetUsers)
		maintenanceTeam.GET("/properties", property.GetProperties)
		// Mount dashboard endpoints for maintenance team
		DashboardRouter(maintenanceTeam)
	}
}
