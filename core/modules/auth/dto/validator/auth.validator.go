package auth_validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// ValidateStrongPassword valida se a senha é forte
// Requer pelo menos: uma letra minúscula, uma letra maiúscula, um número e um caractere especial
func ValidateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	// Expressões regulares para cada requisito
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~]`).MatchString(password)

	return hasLower && hasUpper && hasNumber && hasSpecial
}
