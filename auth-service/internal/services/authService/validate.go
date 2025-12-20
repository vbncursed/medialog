package authService

import (
	"net/mail"
	"strings"
	"unicode"

	"github.com/vbncursed/medialog/auth-service/internal/models"
)

func validateEmail(email string) bool {
	email = strings.TrimSpace(email)
	if len(email) < 3 || len(email) > 254 {
		return false
	}

	// Базовая RFC-проверка.
	if _, err := mail.ParseAddress(email); err != nil {
		return false
	}

	// Дополнительные проверки как в students-service:
	// - ровно один '@'
	// - доменная часть не пустая и не слишком длинная
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if len(parts[1]) == 0 || len(parts[1]) > 253 {
		return false
	}

	return true
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

func normalizeAndValidateAuthInput(in models.AuthInput) (models.AuthInput, error) {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	if !validateEmail(in.Email) || !validatePassword(in.Password) {
		return models.AuthInput{}, ErrInvalidArgument
	}
	return in, nil
}
