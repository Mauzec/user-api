package util

import (
	"regexp"
	"strings"
)

func IsValidPhoneNumber(phone_number string) bool {
	if phone_number == "" {
		return false
	}
	e164Regex := `^\+[1-9]\d{1,14}$`
	re := regexp.MustCompile(e164Regex)
	phone_number = strings.ReplaceAll(phone_number, " ", "")

	return re.Find([]byte(phone_number)) != nil
}
