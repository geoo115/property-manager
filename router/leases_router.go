package router

import (
	"github.com/geoo115/property-manager/api/lease"
	"github.com/gin-gonic/gin"
)

func LeaseRouter(rg *gin.RouterGroup) {
	rg.GET("/leases", lease.GetLeases)
	rg.GET("/leases/:id", lease.GetLeaseByID)
	rg.GET("/leases/active", lease.GetActiveLeaseForTenant)
	rg.GET("/properties/:id/lease", lease.GetLeaseForProperty)
	rg.POST("/leases", lease.CreateLease)
	rg.PUT("/leases/:id", lease.UpdateLease)
	rg.DELETE("/leases/:id", lease.DeleteLease)
}
