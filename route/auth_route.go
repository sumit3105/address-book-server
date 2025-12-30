package route

import (
	"address-book-server/controller"

	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.Engine, authController controller.AuthController) {
	api := router.Group("/api/v1/auth")
	{
		api.POST("/register", authController.Register)
		api.POST("/login", authController.Login)
	}
}
