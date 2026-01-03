package metadata_service

import (
	"context"
	"time"

	"github.com/vbncursed/medialog/metadata_service/internal/models"
)

func (s *MetadataService) CreateMedia(ctx context.Context, input models.CreateMediaInput) (*models.Media, error) {
	if err := s.validateCreateInput(input); err != nil {
		return nil, err
	}

	media := &models.Media{
		Type:        input.Type,
		Title:       input.Title,
		Year:        input.Year,
		Genres:      input.Genres,
		PosterURL:   input.PosterURL,
		CoverURL:    input.CoverURL,
		ExternalIDs: input.ExternalIDs,
		UpdatedAt:   time.Now(),
	}

	if err := s.storage.CreateMedia(ctx, media); err != nil {
		return nil, err
	}

	return media, nil
}

