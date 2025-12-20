package authService

import "errors"

var (
	ErrInvalidArgument     = errors.New("invalid argument")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrSessionRevoked      = errors.New("session revoked")
	ErrSessionExpired      = errors.New("session expired")
)


