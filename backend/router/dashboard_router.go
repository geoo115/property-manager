package router

import (
	"github.com/geoo115/property-manager/api/dashboard"
	"github.com/gin-gonic/gin"
)

func DashboardRouter(r gin.IRouter) {
	// Dashboard statistics endpoint
	r.GET("/dashboard/stats", dashboard.GetDashboardStats)

	// Dashboard activities endpoint
	r.GET("/dashboard/activities", dashboard.GetRecentActivities)
}
