package bootstrap

import (
	"fmt"

	"github.com/vbncursed/medialog/library_service/config"
	"github.com/vbncursed/medialog/library_service/internal/producer/library_entry_event_producer"
)

func InitLibraryEntryEventProducer(cfg *config.Config) *library_entry_event_producer.LibraryEntryEventProducer {
	kafkaBrokers := []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)}
	metadataServiceAddr := cfg.MetadataService.GRPCAddr
	return library_entry_event_producer.NewLibraryEntryEventProducer(kafkaBrokers, cfg.Kafka.LibraryEntryEventTopic, metadataServiceAddr)
}
