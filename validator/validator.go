package validator

import "github.com/go-playground/validator/v10"

var Validate *validator.Validate

func InitValidator() {
	Validate = validator.New()

	Validate.RegisterValidation("password", PasswordValidator)
	Validate.RegisterValidation("phone", PhoneValidator)
}