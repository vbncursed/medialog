package library_entry_processor

import (
	"context"
	"log/slog"
)

// LibraryEntryEvent представляет событие об изменении записи в библиотеке
type LibraryEntryEvent struct {
	EntryID  uint64 `json:"entry_id"`
	UserID   uint64 `json:"user_id"`
	MediaID  uint64 `json:"media_id"`
	Type     int    `json:"type"`     // MediaType как int
	Status   int    `json:"status"`   // EntryStatus как int
	Rating   uint32 `json:"rating"`
	UpdatedAt int64 `json:"updated_at"`
}

func (p *LibraryEntryProcessor) Handle(ctx context.Context, event *LibraryEntryEvent) error {
	if event.MediaID == 0 {
		slog.Warn("library entry event has no media_id", "entry_id", event.EntryID)
		return nil
	}

	// Проверяем, есть ли метаданные для этого media_id
	media, err := p.metadataService.GetMedia(ctx, event.MediaID)
	if err == nil && media != nil {
		// Метаданные уже есть, ничего не делаем
		slog.Debug("metadata already exists", "media_id", event.MediaID)
		return nil
	}

	// Метаданных нет - пытаемся обогатить из внешних API
	// Это будет реализовано позже, когда добавим external API клиенты
	slog.Info("metadata enrichment requested", "media_id", event.MediaID, "entry_id", event.EntryID)
	
	// TODO: реализовать обогащение метаданных из внешних API
	// 1. Определить тип контента (movie/tv/book)
	// 2. Вызвать соответствующий внешний API (TMDB для movie/tv, Open Library для book)
	// 3. Сохранить метаданные в БД

	return nil
}

