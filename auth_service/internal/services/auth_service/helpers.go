package auth_service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidArgument     = errors.New("invalid argument")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrSessionRevoked      = errors.New("session revoked")
	ErrSessionExpired      = errors.New("session expired")
)

// test hooks
var (
	bcryptGenerate = bcrypt.GenerateFromPassword
	randRead       = rand.Read
)

func passwordHash(password string) (string, error) {
	b, err := bcryptGenerate([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func comparePassword(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func newRefreshToken(refreshTTL time.Duration) (token string, hash []byte, exp time.Time, err error) {
	b := make([]byte, 32)
	if _, err = randRead(b); err != nil {
		return "", nil, time.Time{}, err
	}

	token = base64.RawURLEncoding.EncodeToString(b)
	sum := sha256.Sum256([]byte(token))
	hash = sum[:]
	exp = time.Now().Add(refreshTTL)
	return token, hash, exp, nil
}

func newAccessToken(jwtSecret string, userID uint64, role string, accessTTL time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"iat":  now.Unix(),
		"exp":  now.Add(accessTTL).Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(jwtSecret))
}

func tokenToHash(refreshToken string) (token string, hash []byte, exp time.Time, err error) {
	// refresh токен не содержит exp внутри (opaque).
	// Это helper только для хеширования единообразно.
	sum := sha256.Sum256([]byte(refreshToken))
	return refreshToken, sum[:], time.Time{}, nil
}
