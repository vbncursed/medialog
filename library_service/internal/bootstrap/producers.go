package bootstrap

import (
	"fmt"

	"github.com/vbncursed/medialog/library_service/config"
	libraryentryeventproducer "github.com/vbncursed/medialog/library_service/internal/producer/library_entry_event_producer"
)

func InitLibraryEntryEventProducer(cfg *config.Config) *libraryentryeventproducer.LibraryEntryEventProducer {
	kafkaBrokers := []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)}
	return libraryentryeventproducer.NewLibraryEntryEventProducer(kafkaBrokers, cfg.Kafka.LibraryEntryEventTopic)
}

