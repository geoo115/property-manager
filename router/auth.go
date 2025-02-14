package router

import (
	"github.com/geoo115/property-manager/api/auth"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("/login", auth.LoginHandler)
	r.POST("/register", auth.RegisterHandler)
}
