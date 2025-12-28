package libraryentryeventproducer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/segmentio/kafka-go"
	"github.com/vbncursed/medialog/library_service/internal/models"
)

// PublishEntryChanged публикует событие об изменении записи
func (p *LibraryEntryEventProducer) PublishEntryChanged(ctx context.Context, entry *models.Entry) error {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(p.kafkaBrokers...),
		Topic:    p.topicName,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	event := map[string]interface{}{
		"entry_id":   entry.EntryID,
		"user_id":    entry.UserID,
		"media_id":   entry.MediaID,
		"type":       entry.Type,
		"status":     entry.Status,
		"rating":     entry.Rating,
		"updated_at": entry.UpdatedAt.Unix(),
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed to marshal entry event", "err", err)
		return err
	}

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", entry.EntryID)),
		Value: eventJSON,
	}

	err = writer.WriteMessages(ctx, msg)
	if err != nil {
		slog.Error("failed to publish entry event", "err", err, "entry_id", entry.EntryID)
		return err
	}

	slog.Info("entry event published", "entry_id", entry.EntryID, "topic", p.topicName)
	return nil
}
