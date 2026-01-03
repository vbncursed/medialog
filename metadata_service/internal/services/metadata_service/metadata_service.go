package metadata_service

import (
	"context"
	"errors"

	"github.com/vbncursed/medialog/metadata_service/internal/models"
)

var (
	ErrInvalidInput     = errors.New("invalid input")
	ErrMediaNotFound    = errors.New("media not found")
	ErrInvalidMediaType = errors.New("invalid media type")
)

type MetadataStorage interface {
	SearchMedia(ctx context.Context, input models.SearchMediaInput) (*models.SearchMediaResult, error)
	GetMedia(ctx context.Context, mediaID uint64) (*models.Media, error)
	GetMediaByExternalID(ctx context.Context, source, externalID string) (*models.Media, error)
	CreateMedia(ctx context.Context, media *models.Media) error
}

type MetadataCache interface {
	GetMedia(ctx context.Context, key string) (*models.Media, error)
	SetMedia(ctx context.Context, key string, media *models.Media, ttl int64) error
	GetSearchResults(ctx context.Context, key string) (*models.SearchMediaResult, error)
	SetSearchResults(ctx context.Context, key string, results *models.SearchMediaResult, ttl int64) error
}

type ExternalAPIClient interface {
	SearchMedia(ctx context.Context, query string, mediaType *models.MediaType) ([]*models.Media, error)
	GetMediaByExternalID(ctx context.Context, source, externalID string) (*models.Media, error)
}

type MetadataService struct {
	storage     MetadataStorage
	cache       MetadataCache
	externalAPI ExternalAPIClient
	mediaTTL    int64
	searchTTL   int64
}

func NewMetadataService(
	storage MetadataStorage,
	cache MetadataCache,
	externalAPI ExternalAPIClient,
	mediaTTL int64,
	searchTTL int64,
) *MetadataService {
	return &MetadataService{
		storage:     storage,
		cache:       cache,
		externalAPI: externalAPI,
		mediaTTL:    mediaTTL,
		searchTTL:   searchTTL,
	}
}
