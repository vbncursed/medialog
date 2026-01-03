package library_service

import (
	"context"
	"errors"
	"time"

	"github.com/vbncursed/medialog/library_service/internal/models"
	"github.com/vbncursed/medialog/library_service/internal/storage/library_storage"
)

func (s *LibraryService) UpdateEntry(ctx context.Context, in models.UpdateEntryInput) (*models.Entry, error) {
	entry, err := s.storage.GetEntry(ctx, in.EntryID, in.UserID)
	if err != nil {
		if errors.Is(err, library_storage.ErrEntryNotFound) {
			return nil, ErrEntryNotFound
		}
		return nil, err
	}

	if entry.UserID != in.UserID {
		return nil, ErrUnauthorized
	}

	if in.Status != nil {
		if err := validateEntryStatus(*in.Status); err != nil {
			return nil, err
		}
		entry.Status = *in.Status
	}

	if in.Rating != nil {
		if err := validateRating(*in.Rating); err != nil {
			return nil, err
		}
		entry.Rating = *in.Rating
	}

	if in.Review != nil {
		entry.Review = *in.Review
	}

	if in.Tags != nil {
		entry.Tags = in.Tags
	}
	if entry.Tags == nil {
		entry.Tags = []string{}
	}

	if in.StartedAt != nil {
		if *in.StartedAt > 0 {
			t := time.Unix(*in.StartedAt, 0)
			entry.StartedAt = &t
		} else {
			entry.StartedAt = nil
		}
	}

	if in.FinishedAt != nil {
		if *in.FinishedAt > 0 {
			t := time.Unix(*in.FinishedAt, 0)
			entry.FinishedAt = &t
		} else {
			entry.FinishedAt = nil
		}
	}

	entry.UpdatedAt = time.Now()

	if err := s.storage.UpdateEntry(ctx, entry); err != nil {
		if errors.Is(err, library_storage.ErrEntryNotFound) {
			return nil, ErrEntryNotFound
		}
		return nil, err
	}

	if err := s.producer.PublishEntryChanged(ctx, entry); err != nil {
		// Логируем ошибку, но не возвращаем её, так как основная операция выполнена
	}

	return entry, nil
}
