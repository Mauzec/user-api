package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/mauzec/user-api/internal/util"
)

func validPhone(fl validator.FieldLevel) bool {
	if phone, ok := fl.Field().Interface().(string); ok {
		return util.IsValidPhoneNumber(phone)
	}
	return false
}
func validSex(fl validator.FieldLevel) bool {
	if sex, ok := fl.Field().Interface().(string); ok {
		return sex == "M" || sex == "F"
	}
	return false
}
