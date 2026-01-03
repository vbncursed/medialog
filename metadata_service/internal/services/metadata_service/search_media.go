package metadata_service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/vbncursed/medialog/metadata_service/internal/models"
)

func (s *MetadataService) SearchMedia(ctx context.Context, input models.SearchMediaInput) (*models.SearchMediaResult, error) {
	if err := s.validateSearchInput(input); err != nil {
		return nil, err
	}

	cacheKey := s.generateSearchCacheKey(input)
	if s.cache != nil {
		if cached, err := s.cache.GetSearchResults(ctx, cacheKey); err == nil && cached != nil {
			return cached, nil
		}
	}

	if input.ExternalID != nil {
		media, err := s.GetMediaByExternalID(ctx, input.ExternalID.Source, input.ExternalID.ExternalID)
		if err == nil && media != nil {
			result := &models.SearchMediaResult{
				Results:  []*models.Media{media},
				Total:    1,
				Page:     input.Page,
				PageSize: input.PageSize,
			}
			if s.cache != nil {
				_ = s.cache.SetSearchResults(ctx, cacheKey, result, s.searchTTL)
			}
			return result, nil
		}
	}

	result, err := s.storage.SearchMedia(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(result.Results) == 0 && input.Query != "" && s.externalAPI != nil {
		externalResults, err := s.externalAPI.SearchMedia(ctx, input.Query, input.Type)
		if err != nil {
		} else if len(externalResults) > 0 {
			for _, media := range externalResults {
				if len(media.ExternalIDs) > 0 {
					existingMedia, err := s.storage.GetMediaByExternalID(ctx, media.ExternalIDs[0].Source, media.ExternalIDs[0].ExternalID)
					if err == nil && existingMedia != nil {
						media.MediaID = existingMedia.MediaID
						continue
					}
				}
				if err := s.storage.CreateMedia(ctx, media); err != nil {
					if len(media.ExternalIDs) > 0 {
						existingMedia, err := s.storage.GetMediaByExternalID(ctx, media.ExternalIDs[0].Source, media.ExternalIDs[0].ExternalID)
						if err == nil && existingMedia != nil {
							media.MediaID = existingMedia.MediaID
						}
					}
				}
			}
			result.Results = externalResults
			result.Total = uint32(len(externalResults))
		}
	}

	if s.cache != nil {
		_ = s.cache.SetSearchResults(ctx, cacheKey, result, s.searchTTL)
	}

	return result, nil
}

func (s *MetadataService) generateSearchCacheKey(input models.SearchMediaInput) string {
	key := fmt.Sprintf("search:%s", input.Query)
	if input.Type != nil {
		key += fmt.Sprintf(":type:%d", *input.Type)
	}
	if input.ExternalID != nil {
		key += fmt.Sprintf(":ext:%s:%s", input.ExternalID.Source, input.ExternalID.ExternalID)
	}
	key += fmt.Sprintf(":page:%d:size:%d", input.Page, input.PageSize)

	hash := md5.Sum([]byte(key))
	return fmt.Sprintf("metadata:search:%s", hex.EncodeToString(hash[:]))
}
