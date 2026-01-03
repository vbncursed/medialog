package library_entry_changed_consumer

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/vbncursed/medialog/metadata_service/internal/services/processors/library_entry_processor"
)

func (c *LibraryEntryChangedConsumer) Consume(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           c.kafkaBrokers,
		GroupID:           "MetadataService_group",
		Topic:             c.topicName,
		HeartbeatInterval: 3 * time.Second,
		SessionTimeout:    30 * time.Second,
	})
	defer r.Close()

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			slog.Error("LibraryEntryChangedConsumer.consume error", "error", err)
			continue
		}

		var event library_entry_processor.LibraryEntryEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			slog.Error("failed to unmarshal library entry event", "error", err)
			continue
		}

		err = c.libraryEntryProcessor.Handle(ctx, &event)
		if err != nil {
			slog.Error("failed to handle library entry event", "error", err, "entry_id", event.EntryID)
		} else {
			slog.Info("library entry event processed", "entry_id", event.EntryID, "media_id", event.MediaID)
		}
	}
}

