package library_service

import "errors"

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrEntryNotFound      = errors.New("entry not found")
	ErrEntryAlreadyExists = errors.New("entry already exists")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidMediaID     = errors.New("invalid media id")
	ErrInvalidMediaType   = errors.New("invalid media type")
	ErrInvalidStatus      = errors.New("invalid status")
	ErrInvalidRating      = errors.New("invalid rating")
)
