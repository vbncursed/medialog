package bootstrap

import (
	"github.com/vbncursed/medialog/library_service/internal/producer/library_entry_event_producer"
	"github.com/vbncursed/medialog/library_service/internal/services/library_service"
	"github.com/vbncursed/medialog/library_service/internal/storage/library_storage"
)

func InitLibraryService(
	storage *library_storage.LibraryStorage,
	producer *library_entry_event_producer.LibraryEntryEventProducer,
) *library_service.LibraryService {
	return library_service.NewLibraryService(storage, producer)
}
