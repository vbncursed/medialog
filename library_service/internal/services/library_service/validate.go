package library_service

import (
	"github.com/vbncursed/medialog/library_service/internal/models"
)

func validateMediaType(t models.MediaType) error {
	switch t {
	case models.MediaTypeMovie, models.MediaTypeTV, models.MediaTypeBook:
		return nil
	default:
		return ErrInvalidMediaType
	}
}

func validateEntryStatus(s models.EntryStatus) error {
	switch s {
	case models.EntryStatusPlanned, models.EntryStatusInProgress, models.EntryStatusDone:
		return nil
	default:
		return ErrInvalidStatus
	}
}

func validateRating(rating uint32) error {
	if rating > 10 {
		return ErrInvalidRating
	}
	return nil
}

func validateMediaID(mediaID uint64) error {
	if mediaID == 0 {
		return ErrInvalidMediaID
	}
	return nil
}
