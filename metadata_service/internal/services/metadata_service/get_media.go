package metadata_service

import (
	"context"
	"fmt"

	"github.com/vbncursed/medialog/metadata_service/internal/models"
)

func (s *MetadataService) GetMedia(ctx context.Context, mediaID uint64) (*models.Media, error) {
	if mediaID == 0 {
		return nil, ErrInvalidInput
	}

	cacheKey := fmt.Sprintf("metadata:media:%d", mediaID)
	if cached, err := s.cache.GetMedia(ctx, cacheKey); err == nil && cached != nil {
		return cached, nil
	}

	media, err := s.storage.GetMedia(ctx, mediaID)
	if err != nil {
		return nil, err
	}

	_ = s.cache.SetMedia(ctx, cacheKey, media, s.mediaTTL)

	return media, nil
}
