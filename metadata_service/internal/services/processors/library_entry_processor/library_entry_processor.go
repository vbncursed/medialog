package library_entry_processor

import (
	"context"

	"github.com/vbncursed/medialog/metadata_service/internal/models"
)

type metadataService interface {
	GetMedia(ctx context.Context, mediaID uint64) (*models.Media, error)
	GetMediaByExternalID(ctx context.Context, source, externalID string) (*models.Media, error)
}

type LibraryEntryProcessor struct {
	metadataService metadataService
}

func NewLibraryEntryProcessor(metadataService metadataService) *LibraryEntryProcessor {
	return &LibraryEntryProcessor{
		metadataService: metadataService,
	}
}
