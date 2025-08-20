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
func validGender(fl validator.FieldLevel) bool {
	if gender, ok := fl.Field().Interface().(string); ok {
		return gender == "M" || gender == "F"
	}
	return false
}
