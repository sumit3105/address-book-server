package route

import (
	"address-book-server/controller"
	"address-book-server/middleware"	

	"github.com/gin-gonic/gin"
)

func AddressRoute(router *gin.Engine, addressController controller.AddressController) {
	addressApi := router.Group("/api/v1/address")
	addressApi.Use(middleware.AuthMiddleware())
	{
		addressApi.GET("/", addressController.List)
		addressApi.POST("/", addressController.Create)
		addressApi.PUT("/:id", addressController.Update)
		addressApi.DELETE("/:id", addressController.Delete)
		addressApi.POST("/export", addressController.Export)
	}
}
