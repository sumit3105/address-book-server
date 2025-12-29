package main

import (
	"address-book-server/controller"
	"address-book-server/logger"
	"address-book-server/middleware"
	"address-book-server/repository"
	"address-book-server/route"
	"address-book-server/service"
	"address-book-server/utils"
	"address-book-server/validator"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)


func main() {
	logger.InitLogger()
	defer logger.Log.Sync()

	err := godotenv.Load()
	if err != nil {
		logger.Log.Error("Error loading the .env file", zap.Error(err))
	}

	db := utils.Connect()
	utils.PerformMigration(db)

	validator.InitValidator()	

	userRepo := repository.NewUserRepository(db)
	userService := service.NewAuthService(userRepo)
	authController := controller.NewAuthController(userService)

	addressRepo := repository.NewAddressRepository(db)
	addressService := service.NewAddressService(addressRepo)
	addressController := controller.NewAddressController(addressService)

	r := gin.New()
	r.Use(middleware.ReuqestLogger())
	r.Use(gin.Recovery())
	r.Use(middleware.ErrorHandler())
	
	route.AuthRoute(r, authController)
	route.AddressRoute(r, addressController)
	
	r.Run(":8080")
}