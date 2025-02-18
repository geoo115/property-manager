package router

import (
	"github.com/geoo115/property-manager/api/lease"
	"github.com/geoo115/property-manager/api/maintenance"
	"github.com/geoo115/property-manager/api/property"
	"github.com/geoo115/property-manager/api/user"
	"github.com/geoo115/property-manager/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	// Public routes
	AuthRoutes(r)

	// Admin group: full access to all endpoints.
	admin := r.Group("/admin")
	admin.Use(middleware.JWTMiddleware(), middleware.RoleMiddleware("admin"))
	{
		UserRouter(admin)
		PropertyRouter(admin)
		LeaseRouter(admin)
		MaintenanceRoutes(admin)
	}

	// Landlord group: restricted access to their properties and leases.
	landlord := r.Group("/landlord")
	landlord.Use(middleware.JWTMiddleware(), middleware.RoleMiddleware("landlord"))
	{
		landlord.GET("/properties", property.GetProperties)
		landlord.GET("/properties/:id", property.GetPropertyByID)
		landlord.GET("/leases", lease.GetLeases)
		landlord.GET("/leases/:id", lease.GetLeaseByID)
		landlord.GET("/properties/:id/maintenances", maintenance.GetLandlordMaintenances)
		landlord.POST("/properties/:id/maintenances", maintenance.CreateLandlordMaintenance)
	}

	// Tenant group: can only access their leases.
	tenant := r.Group("/tenant")
	tenant.Use(middleware.JWTMiddleware(), middleware.RoleMiddleware("tenant"))
	{
		tenant.GET("/leases", lease.GetLeases)
		tenant.GET("/leases/:id", lease.GetLeaseByID)
		tenant.GET("/leases/:id/maintenance", maintenance.GetMaintenances)
		tenant.POST("/leases/:id/maintenance", maintenance.CreateMaintenance)
	}

	maintenanceTeam := r.Group("/maintenanceTeam")
	maintenanceTeam.Use(middleware.JWTMiddleware(), middleware.RoleMiddleware("maintenanceTeam"))
	{
		maintenanceTeam.GET("/maintenances", maintenance.GetMaintenances)
		maintenanceTeam.GET("/maintenance/:id", maintenance.GetMaintenance)
		maintenanceTeam.PUT("/maintenance/:id", maintenance.UpdateMaintenance)
		maintenanceTeam.GET("/users", user.GetUsers)
		maintenanceTeam.GET("/properties", property.GetProperties)
	}
}
