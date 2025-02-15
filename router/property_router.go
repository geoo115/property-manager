package router

import (
	"github.com/geoo115/property-manager/api/property"
	"github.com/gin-gonic/gin"
)

func PropertyRouter(r *gin.Engine) {
	r.GET("/properties", property.GetProperties)
	r.GET("/properties/:id", property.GetPropertyByID)
	r.POST("/properties", property.CreateProperty)
	r.PUT("/properties/:id", property.UpdateProperty)
	r.DELETE("/properties/:id", property.DeleteProperty)
}
