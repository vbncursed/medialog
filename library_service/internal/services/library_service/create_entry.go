package library_service

import (
	"context"
	"errors"
	"time"

	"github.com/vbncursed/medialog/library_service/internal/models"
	"github.com/vbncursed/medialog/library_service/internal/storage/library_storage"
)

func (s *LibraryService) CreateEntry(ctx context.Context, in models.CreateEntryInput) (*models.Entry, error) {
	// Валидация
	if err := validateMediaID(in.MediaID); err != nil {
		return nil, err
	}
	if err := validateMediaType(in.Type); err != nil {
		return nil, err
	}
	if err := validateEntryStatus(in.Status); err != nil {
		return nil, err
	}
	if err := validateRating(in.Rating); err != nil {
		return nil, err
	}

	now := time.Now()

	var startedAt *time.Time
	if in.StartedAt > 0 {
		t := time.Unix(in.StartedAt, 0)
		startedAt = &t
	}

	var finishedAt *time.Time
	if in.FinishedAt > 0 {
		t := time.Unix(in.FinishedAt, 0)
		finishedAt = &t
	}

	tags := in.Tags
	if tags == nil {
		tags = []string{}
	}

	entry := &models.Entry{
		UserID:     in.UserID,
		MediaID:    in.MediaID,
		Type:       in.Type,
		Status:     in.Status,
		Rating:     in.Rating,
		Review:     in.Review,
		Tags:       tags,
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.storage.CreateEntry(ctx, entry); err != nil {
		if errors.Is(err, library_storage.ErrEntryAlreadyExists) {
			return nil, ErrEntryAlreadyExists
		}
		return nil, err
	}

	if err := s.producer.PublishEntryChanged(ctx, entry); err != nil {
		// Логируем ошибку, но не возвращаем её, так как основная операция выполнена
	}

	return entry, nil
}
