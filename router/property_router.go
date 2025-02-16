package router

import (
	"github.com/geoo115/property-manager/api/property"
	"github.com/gin-gonic/gin"
)

func PropertyRouter(rg *gin.RouterGroup) {
	rg.GET("/properties", property.GetProperties)
	rg.GET("/properties/:id", property.GetPropertyByID)
	rg.POST("/properties", property.CreateProperty)
	rg.PUT("/properties/:id", property.UpdateProperty)
	rg.DELETE("/properties/:id", property.DeleteProperty)
}
