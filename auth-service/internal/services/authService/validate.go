package authService

import (
	"net/mail"
	"strings"
	"unicode"
)

func validateEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}
	_, err := mail.ParseAddress(email)
	return err == nil
}

func validatePassword(password string) bool {
	p := strings.TrimSpace(password)
	if len(p) < 8 || len(p) > 128 {
		return false
	}

	var hasUpper, hasLower, hasDigit bool
	for _, r := range p {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}
