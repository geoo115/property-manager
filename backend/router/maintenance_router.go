package router

import (
	"github.com/geoo115/property-manager/api/maintenance"
	"github.com/gin-gonic/gin"
)

func MaintenanceRoutes(rg *gin.RouterGroup) {
	rg.GET("/maintenances", maintenance.GetMaintenances)
	rg.GET("/maintenance/:id", maintenance.GetMaintenance)
	rg.POST("/leases/:leaseID/maintenance", maintenance.CreateMaintenanceByLease)            // Tenant
	rg.POST("/properties/:propertyID/maintenances", maintenance.CreateMaintenanceByProperty) // Admin/Landlord
	rg.PUT("/maintenance/:id", maintenance.UpdateMaintenance)
	rg.DELETE("/maintenance/:id", maintenance.DeleteMaintenance)
}
