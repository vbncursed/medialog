package authService

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type tokenPair struct {
	AccessToken  string
	RefreshToken string
	RefreshHash  []byte
	RefreshExp   time.Time
}

// test hooks
var randRead = rand.Read

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

func newAccessToken(jwtSecret string, userID uint64, accessTTL time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(accessTTL).Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(jwtSecret))
}
