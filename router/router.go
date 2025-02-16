package router

import (
	"github.com/geoo115/property-manager/api/lease"
	"github.com/geoo115/property-manager/api/property"
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
	}

	// Landlord group: restricted access to their properties and leases.
	landlord := r.Group("/landlord")
	landlord.Use(middleware.JWTMiddleware(), middleware.RoleMiddleware("landlord"))
	{
		landlord.GET("/properties", property.GetProperties)
		landlord.GET("/properties/:id", property.GetPropertyByID)
		landlord.GET("/leases", lease.GetLeases)
		landlord.GET("/leases/:id", lease.GetLeaseByID)
	}

	// Tenant group: can only access their leases.
	tenant := r.Group("/tenant")
	tenant.Use(middleware.JWTMiddleware(), middleware.RoleMiddleware("tenant"))
	{
		tenant.GET("/leases", lease.GetLeases)
		tenant.GET("/leases/:id", lease.GetLeaseByID)
	}
}
