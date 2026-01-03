package metadata_service

import (
	"context"
	"fmt"

	"github.com/vbncursed/medialog/metadata_service/internal/models"
)

func (s *MetadataService) GetMediaByExternalID(ctx context.Context, source, externalID string) (*models.Media, error) {
	if source == "" || externalID == "" {
		return nil, ErrInvalidInput
	}

	cacheKey := fmt.Sprintf("metadata:external:%s:%s", source, externalID)
	if s.cache != nil {
		if cached, err := s.cache.GetMedia(ctx, cacheKey); err == nil && cached != nil {
			return cached, nil
		}
	}

	media, err := s.storage.GetMediaByExternalID(ctx, source, externalID)
	if err == nil && media != nil {
		if s.cache != nil {
			_ = s.cache.SetMedia(ctx, cacheKey, media, s.mediaTTL)
		}
		return media, nil
	}

	if s.externalAPI != nil {
		media, err := s.externalAPI.GetMediaByExternalID(ctx, source, externalID)
		if err == nil && media != nil {
			_ = s.storage.CreateMedia(ctx, media)
			if s.cache != nil {
				_ = s.cache.SetMedia(ctx, cacheKey, media, s.mediaTTL)
			}
			return media, nil
		}
	}

	return nil, ErrMediaNotFound
}
