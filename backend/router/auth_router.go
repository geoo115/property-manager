package router

import (
	"github.com/geoo115/property-manager/api/auth"
	"github.com/geoo115/property-manager/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r gin.IRouter) {
	// Allow 10 requests per minute for the login endpoint.
	r.POST("/login", auth.LoginHandler)
	r.POST("/register", auth.RegisterHandler)
	r.POST("/refresh-token", middleware.RefreshTokenHandler)
	r.POST("/logout", auth.LogoutHandler) // Optional: Add RateLimit if needed
}
