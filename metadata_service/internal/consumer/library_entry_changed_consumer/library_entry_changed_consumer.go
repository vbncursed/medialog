package library_entry_changed_consumer

import (
	"context"

	"github.com/vbncursed/medialog/shared/events"
)

type libraryEntryProcessor interface {
	Handle(ctx context.Context, event *events.LibraryEntryEvent) error
}

type LibraryEntryChangedConsumer struct {
	libraryEntryProcessor libraryEntryProcessor
	kafkaBrokers          []string
	topicName             string
}

func NewLibraryEntryChangedConsumer(
	libraryEntryProcessor libraryEntryProcessor,
	kafkaBrokers []string,
	topicName string,
) *LibraryEntryChangedConsumer {
	return &LibraryEntryChangedConsumer{
		libraryEntryProcessor: libraryEntryProcessor,
		kafkaBrokers:          kafkaBrokers,
		topicName:             topicName,
	}
}

