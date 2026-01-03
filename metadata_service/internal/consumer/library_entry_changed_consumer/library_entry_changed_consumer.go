package library_entry_changed_consumer

import (
	"context"

	"github.com/vbncursed/medialog/metadata_service/internal/services/processors/library_entry_processor"
)

type libraryEntryProcessor interface {
	Handle(ctx context.Context, event *library_entry_processor.LibraryEntryEvent) error
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

