package validator

import (
	"gopkg.in/go-playground/validator.v10"
)

var validate = validator.New()

func ValidateStruct(payload any) error {
	return validate.Struct(payload)
}
