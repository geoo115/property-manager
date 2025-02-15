package router

import (
	"github.com/geoo115/property-manager/api/lease"
	"github.com/gin-gonic/gin"
)

func LeaseRouter(r *gin.Engine) {
	r.GET("/leases", lease.GetLeases)
	r.GET("/leases/:id", lease.GetLeaseByID)
	r.POST("/leases/", lease.CreateLease)
	r.PUT("/leases/:id", lease.UpdateLease)
	r.DELETE("leases/:id", lease.DeleteLease)
}
