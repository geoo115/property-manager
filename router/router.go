package router

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	AuthRoutes(r)
	UserRouter(r)
	PropertyRouter(r)
}
