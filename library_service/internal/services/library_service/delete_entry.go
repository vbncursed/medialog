package library_service

import (
	"context"
	"errors"

	"github.com/vbncursed/medialog/library_service/internal/storage/library_storage"
)

func (s *LibraryService) DeleteEntry(ctx context.Context, entryID, userID uint64) error {
	entry, err := s.storage.GetEntry(ctx, entryID, userID)
	if err != nil {
		if errors.Is(err, library_storage.ErrEntryNotFound) {
			return ErrEntryNotFound
		}
		return err
	}

	if entry.UserID != userID {
		return ErrUnauthorized
	}

	if err := s.storage.DeleteEntry(ctx, entryID, userID); err != nil {
		if errors.Is(err, library_storage.ErrEntryNotFound) {
			return ErrEntryNotFound
		}
		return err
	}

	// Публикуем событие об удалении записи
	if err := s.producer.PublishEntryChanged(ctx, entry); err != nil {
		// Логируем ошибку, но не возвращаем её, так как основная операция выполнена
	}

	return nil
}
