package library_service

import (
	"context"

	"github.com/vbncursed/medialog/library_service/internal/models"
	libraryentryeventproducer "github.com/vbncursed/medialog/library_service/internal/producer/library_entry_event_producer"
)

type LibraryStorage interface {
	CreateEntry(ctx context.Context, entry *models.Entry) error
	GetEntry(ctx context.Context, entryID, userID uint64) (*models.Entry, error)
	UpdateEntry(ctx context.Context, entry *models.Entry) error
	DeleteEntry(ctx context.Context, entryID, userID uint64) error
	ListEntries(ctx context.Context, input models.ListEntriesInput) (*models.ListEntriesResult, error)
}

type LibraryService struct {
	storage  LibraryStorage
	producer *libraryentryeventproducer.LibraryEntryEventProducer
}

func NewLibraryService(storage LibraryStorage, producer *libraryentryeventproducer.LibraryEntryEventProducer) *LibraryService {
	return &LibraryService{
		storage:  storage,
		producer: producer,
	}
}
