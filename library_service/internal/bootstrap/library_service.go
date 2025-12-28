package bootstrap

import (
	libraryentryeventproducer "github.com/vbncursed/medialog/library_service/internal/producer/library_entry_event_producer"
	"github.com/vbncursed/medialog/library_service/internal/services/library_service"
	"github.com/vbncursed/medialog/library_service/internal/storage/library_storage"
)

func InitLibraryService(
	storage *library_storage.LibraryStorage,
	producer *libraryentryeventproducer.LibraryEntryEventProducer,
) *library_service.LibraryService {
	return library_service.NewLibraryService(storage, producer)
}
