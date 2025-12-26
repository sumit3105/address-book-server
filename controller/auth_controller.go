package controller

import (
	"address-book-server/dto"
	appError "address-book-server/error"
	"address-book-server/service"
	"address-book-server/utils"
	"address-book-server/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type authController struct {
	userService service.AuthService
}

func NewAuthController(userService service.AuthService) AuthController {
	return &authController{userService: userService}
}

func (c *authController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(
			appError.BadRequest(
				"Invalid request body",
				err,
			),
		)
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		ctx.Error(
			appError.NewValidationError(
				utils.FormatValidationErrors(err),
			),
		)
		return
	}

	if err := c.userService.Register(req.Email, req.Password); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"message": "User registered successfully",
	})
}

func (c *authController) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(
			appError.NewValidationError(
				utils.FormatValidationErrors(err),
			),
		)
		return
	}

	token, err := c.userService.Login(req.Email, req.Password)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token": token,
	})

}