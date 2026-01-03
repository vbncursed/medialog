package library_entry_changed_consumer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/vbncursed/medialog/shared/events"
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
			continue
		}

		var event events.LibraryEntryEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			continue
		}

		err = c.libraryEntryProcessor.Handle(ctx, &event)
	}
}
