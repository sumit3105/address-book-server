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
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading the .env file")
	}

	db := utils.Connect()

	utils.PerformMigration(db)

	validator.InitValidator()

	logger.InitLogger()
	defer logger.Log.Sync()

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