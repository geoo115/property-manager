package router

import (
	"github.com/geoo115/property-manager/api/user"
	"github.com/gin-gonic/gin"
)

func UserRouter(rg *gin.RouterGroup) {
	rg.GET("/users", user.GetUsers)
	rg.GET("/users/active", user.GetActiveUsers)
	rg.GET("/users/:id", user.GetUserByID)
	rg.POST("/users", user.CreateUser)
	rg.PUT("/users/:id", user.UpdateUser)
	rg.DELETE("/users/:id", user.DeleteUser)
}
