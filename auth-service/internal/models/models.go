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


