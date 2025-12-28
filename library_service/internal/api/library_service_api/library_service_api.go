package library_service_api

import (
	"context"

	"github.com/vbncursed/medialog/library_service/internal/models"
	"github.com/vbncursed/medialog/library_service/internal/pb/library_api"
	"github.com/vbncursed/medialog/library_service/internal/services/library_service"
)

type libraryService interface {
	CreateEntry(ctx context.Context, in models.CreateEntryInput) (*models.Entry, error)
	GetEntry(ctx context.Context, entryID, userID uint64) (*models.Entry, error)
	UpdateEntry(ctx context.Context, in models.UpdateEntryInput) (*models.Entry, error)
	DeleteEntry(ctx context.Context, entryID, userID uint64) error
	ListEntries(ctx context.Context, in models.ListEntriesInput) (*models.ListEntriesResult, error)
}

// LibraryServiceAPI реализует grpc LibraryServiceServer
type LibraryServiceAPI struct {
	library_api.UnimplementedLibraryServiceServer
	libraryService libraryService
	jwtSecret      string
}

func NewLibraryServiceAPI(libraryService *library_service.LibraryService, jwtSecret string) *LibraryServiceAPI {
	return &LibraryServiceAPI{
		libraryService: libraryService,
		jwtSecret:      jwtSecret,
	}
}
