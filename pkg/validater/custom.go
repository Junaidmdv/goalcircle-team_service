package validater

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func phoneValidation(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	// E.164 format: + followed by 1–15 digits
	regex := `^\+[1-9]\d{1,14}$`

	match, _ := regexp.MatchString(regex, phone)
	return match
}

func passwordValidation(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[@$!%*?&]`).MatchString(password)
	hasMinLen := len(password) >= 8

	return hasUpper && hasLower && hasNumber && hasSpecial && hasMinLen
}
