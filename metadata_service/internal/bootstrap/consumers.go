package bootstrap

import (
	"fmt"

	"github.com/vbncursed/medialog/metadata_service/config"
	library_entry_changed_consumer "github.com/vbncursed/medialog/metadata_service/internal/consumer/library_entry_changed_consumer"
	library_entry_processor "github.com/vbncursed/medialog/metadata_service/internal/services/processors/library_entry_processor"
)

func InitLibraryEntryChangedConsumer(
	cfg *config.Config,
	libraryEntryProcessor *library_entry_processor.LibraryEntryProcessor,
) *library_entry_changed_consumer.LibraryEntryChangedConsumer {
	kafkaBrokers := []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)}
	return library_entry_changed_consumer.NewLibraryEntryChangedConsumer(
		libraryEntryProcessor,
		kafkaBrokers,
		cfg.Kafka.LibraryEntryChangedTopic,
	)
}

