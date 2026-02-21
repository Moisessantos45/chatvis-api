package pkg

import (
	"errors"
	"regexp"
)

// isValidEmail verifica si el email tiene un formato válido
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// validatePasswordStrength valida la fortaleza de la contraseña
func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("la contraseña debe tener al menos 8 caracteres")
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

	if !hasUpper {
		return errors.New("la contraseña debe contener al menos una letra mayúscula")
	}
	if !hasLower {
		return errors.New("la contraseña debe contener al menos una letra minúscula")
	}
	if !hasNumber {
		return errors.New("la contraseña debe contener al menos un número")
	}
	if !hasSpecial {
		return errors.New("la contraseña debe contener al menos un carácter especial")
	}

	return nil
}
