package library_service

import (
	"context"
	"errors"

	"github.com/vbncursed/medialog/library_service/internal/models"
	"github.com/vbncursed/medialog/library_service/internal/storage/library_storage"
)

func (s *LibraryService) GetEntry(ctx context.Context, entryID, userID uint64) (*models.Entry, error) {
	entry, err := s.storage.GetEntry(ctx, entryID, userID)
	if err != nil {
		if errors.Is(err, library_storage.ErrEntryNotFound) {
			return nil, ErrEntryNotFound
		}
		return nil, err
	}

	if entry.UserID != userID {
		return nil, ErrUnauthorized
	}

	return entry, nil
}
