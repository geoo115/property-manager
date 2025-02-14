package router

import (
	"github.com/geoo115/property-manager/api/user"
	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.Engine) {
	r.GET("users", user.GetUsers)
	r.GET("/users/id", user.GetUserByID)
	r.POST("/users", user.CreateUser)
	r.PUT("/users/id", user.UpdateUser)
	r.DELETE("/users/id", user.DeleteUser)
}
