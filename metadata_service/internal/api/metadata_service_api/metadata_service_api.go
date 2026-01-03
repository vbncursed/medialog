package metadata_service_api

import (
	"context"

	"github.com/vbncursed/medialog/metadata_service/internal/models"
	"github.com/vbncursed/medialog/metadata_service/internal/pb/metadata_api"
	"github.com/vbncursed/medialog/metadata_service/internal/services/metadata_service"
)

type metadataService interface {
	SearchMedia(ctx context.Context, in models.SearchMediaInput) (*models.SearchMediaResult, error)
	GetMedia(ctx context.Context, mediaID uint64) (*models.Media, error)
	GetMediaByExternalID(ctx context.Context, source, externalID string) (*models.Media, error)
	CreateMedia(ctx context.Context, in models.CreateMediaInput) (*models.Media, error)
}

// MetadataServiceAPI реализует grpc MetadataServiceServer
type MetadataServiceAPI struct {
	metadata_api.UnimplementedMetadataServiceServer
	metadataService metadataService
	jwtSecret       string
}

func NewMetadataServiceAPI(metadataService *metadata_service.MetadataService, jwtSecret string) *MetadataServiceAPI {
	return &MetadataServiceAPI{
		metadataService: metadataService,
		jwtSecret:       jwtSecret,
	}
}

