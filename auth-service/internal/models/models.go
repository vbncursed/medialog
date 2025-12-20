package models

import "time"

type User struct {
	ID           uint64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

type Session struct {
	ID          uint64
	UserID      uint64
	RefreshHash []byte
	ExpiresAt   time.Time
	RevokedAt   *time.Time
	CreatedAt   time.Time
	UserAgent   string
	IP          string
}

// ---- Input DTOs (service-layer request models) ----

// RegisterInput — данные для регистрации пользователя.
// Это не "пользовательская" доменная модель, т.к. содержит plaintext пароль.
type RegisterInput struct {
	Email     string
	Password  string
	UserAgent string
	IP        string
}

// LoginInput — данные для логина пользователя.
type LoginInput struct {
	Email     string
	Password  string
	UserAgent string
	IP        string
}

// RefreshInput — данные для refresh операции.
type RefreshInput struct {
	RefreshToken string
	UserAgent    string
	IP           string
}
