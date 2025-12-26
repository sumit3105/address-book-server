package validator

import (
	"unicode"
	"github.com/go-playground/validator/v10"
)

func PhoneValidator(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	if len(phone) < 10 {
		return false
	}

	for _, ch := range phone {
		if !unicode.IsDigit(ch) {
			return false
		}
	}

	return true
}
