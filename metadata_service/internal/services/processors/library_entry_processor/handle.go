package library_entry_processor

import (
	"context"

	"github.com/vbncursed/medialog/metadata_service/internal/models"
	"github.com/vbncursed/medialog/shared/events"
)

func (p *LibraryEntryProcessor) Handle(ctx context.Context, event *events.LibraryEntryEvent) error {
	if event.MediaID == 0 {
		return nil
	}

	media, err := p.metadataService.GetMedia(ctx, event.MediaID)
	if err == nil && media != nil {
		return nil
	}

	mediaType := convertIntToMediaType(event.Type)
	if mediaType == models.MediaTypeUnspecified {
		return nil
	}

	if event.ExternalID != nil && event.ExternalID.Source != "" && event.ExternalID.ExternalID != "" {
		enrichedMedia, err := p.metadataService.GetMediaByExternalID(ctx, event.ExternalID.Source, event.ExternalID.ExternalID)
		if err == nil && enrichedMedia != nil {
			if enrichedMedia.MediaID != event.MediaID {
				return nil
			}
			return nil
		}

		_, err = p.metadataService.GetMediaByExternalID(ctx, event.ExternalID.Source, event.ExternalID.ExternalID)
		if err != nil {
			return nil
		}
	}

	return nil
}

func convertIntToMediaType(t int) models.MediaType {
	switch t {
	case 1:
		return models.MediaTypeMovie
	case 2:
		return models.MediaTypeTV
	case 3:
		return models.MediaTypeBook
	default:
		return models.MediaTypeUnspecified
	}
}
