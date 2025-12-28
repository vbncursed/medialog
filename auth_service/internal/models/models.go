package models

import "time"

type User struct {
	ID           uint64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

type Session struct {
	ID          uint64 `json:"-"`
	UserID      uint64
	RefreshHash []byte
	ExpiresAt   time.Time
	RevokedAt   *time.Time
	CreatedAt   time.Time
	UserAgent   string
	IP          string
}

type AuthInput struct {
	Email     string
	Password  string
	UserAgent string
	IP        string
}

type RegisterInput = AuthInput

type LoginInput = AuthInput

type RefreshInput struct {
	RefreshToken string
	UserAgent    string
	IP           string
}
