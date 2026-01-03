package clients

import (
	"context"

	"github.com/vbncursed/medialog/metadata_service/internal/clients/open_library"
	"github.com/vbncursed/medialog/metadata_service/internal/clients/tmdb"
	"github.com/vbncursed/medialog/metadata_service/internal/models"
)

type ExternalAPIClient struct {
	tmdbClient    *tmdb.TMDBClient
	openLibClient *open_library.OpenLibraryClient
}

func NewExternalAPIClient(
	tmdbClient *tmdb.TMDBClient,
	openLibClient *open_library.OpenLibraryClient,
) *ExternalAPIClient {
	return &ExternalAPIClient{
		tmdbClient:    tmdbClient,
		openLibClient: openLibClient,
	}
}

func (c *ExternalAPIClient) SearchMedia(ctx context.Context, query string, mediaType *models.MediaType) ([]*models.Media, error) {
	var results []*models.Media

	if mediaType == nil || *mediaType == models.MediaTypeMovie || *mediaType == models.MediaTypeTV {
		if c.tmdbClient != nil {
			tmdbResults, err := c.tmdbClient.SearchMedia(ctx, query, mediaType)
			if err == nil {
				results = append(results, tmdbResults...)
			}
		}
	}

	if mediaType == nil || *mediaType == models.MediaTypeBook {
		if c.openLibClient != nil {
			bookResults, err := c.openLibClient.SearchMedia(ctx, query, mediaType)
			if err == nil {
				results = append(results, bookResults...)
			}
		}
	}

	return results, nil
}

func (c *ExternalAPIClient) GetMediaByExternalID(ctx context.Context, source, externalID string) (*models.Media, error) {
	switch source {
	case "tmdb":
		if c.tmdbClient != nil {
			return c.tmdbClient.GetMediaByExternalID(ctx, source, externalID)
		}
	case "isbn", "openlibrary":
		if c.openLibClient != nil {
			return c.openLibClient.GetMediaByExternalID(ctx, source, externalID)
		}
	}

	return nil, nil
}
