package authService

import (
	"crypto/sha256"
	"time"
)

func tokenToHash(refreshToken string) (token string, hash []byte, exp time.Time, err error) {
	// refresh токен не содержит exp внутри (opaque).
	// Это helper только для хеширования единообразно.
	sum := sha256.Sum256([]byte(refreshToken))
	return refreshToken, sum[:], time.Time{}, nil
}


