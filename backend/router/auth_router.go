package router

import (
	"time"

	"github.com/geoo115/property-manager/api/auth"
	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	// Allow 10 requests per minute for the login endpoint.
	r.POST("/login", db.RateLimit(30, time.Minute), auth.LoginHandler)
	r.POST("/register", db.RateLimit(10, time.Minute), auth.RegisterHandler)
	r.POST("/refresh-token", db.RateLimit(10, time.Minute), middleware.RefreshTokenHandler)
	r.POST("/logout", auth.LogoutHandler) // Optional: Add RateLimit if needed
}
