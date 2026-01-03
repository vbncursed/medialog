package metadata_service_api

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vbncursed/medialog/metadata_service/internal/models"
	pb_models "github.com/vbncursed/medialog/metadata_service/internal/pb/models"
)

func convertMediaType(t pb_models.MediaType) models.MediaType {
	switch t {
	case pb_models.MediaType_MEDIA_TYPE_MOVIE:
		return models.MediaTypeMovie
	case pb_models.MediaType_MEDIA_TYPE_TV:
		return models.MediaTypeTV
	case pb_models.MediaType_MEDIA_TYPE_BOOK:
		return models.MediaTypeBook
	default:
		return models.MediaTypeUnspecified
	}
}

func convertMediaTypeToProto(t models.MediaType) pb_models.MediaType {
	switch t {
	case models.MediaTypeMovie:
		return pb_models.MediaType_MEDIA_TYPE_MOVIE
	case models.MediaTypeTV:
		return pb_models.MediaType_MEDIA_TYPE_TV
	case models.MediaTypeBook:
		return pb_models.MediaType_MEDIA_TYPE_BOOK
	default:
		return pb_models.MediaType_MEDIA_TYPE_UNSPECIFIED
	}
}

func convertExternalIDToProto(e *models.ExternalID) *pb_models.ExternalID {
	if e == nil {
		return nil
	}
	return &pb_models.ExternalID{
		Source:     e.Source,
		ExternalId: e.ExternalID,
	}
}

func convertExternalIDFromProto(e *pb_models.ExternalID) *models.ExternalID {
	if e == nil {
		return nil
	}
	return &models.ExternalID{
		Source:     e.Source,
		ExternalID: e.ExternalId,
	}
}

func convertMediaToProto(m *models.Media) *pb_models.Media {
	if m == nil {
		return nil
	}

	var year *uint32
	if m.Year != nil {
		year = m.Year
	}

	externalIDs := make([]*pb_models.ExternalID, 0, len(m.ExternalIDs))
	for _, eid := range m.ExternalIDs {
		externalIDs = append(externalIDs, convertExternalIDToProto(&eid))
	}

	return &pb_models.Media{
		MediaId:    m.MediaID,
		Type:       convertMediaTypeToProto(m.Type),
		Title:      m.Title,
		Year:       year,
		Genres:     m.Genres,
		PosterUrl:  m.PosterURL,
		CoverUrl:   m.CoverURL,
		ExternalIds: externalIDs,
		UpdatedAt:  timestamppb.New(m.UpdatedAt),
	}
}

func convertMediaFromProto(m *pb_models.Media) *models.Media {
	if m == nil {
		return nil
	}

	externalIDs := make([]models.ExternalID, 0, len(m.ExternalIds))
	for _, eid := range m.ExternalIds {
		if extID := convertExternalIDFromProto(eid); extID != nil {
			externalIDs = append(externalIDs, *extID)
		}
	}

	return &models.Media{
		MediaID:    m.MediaId,
		Type:       convertMediaType(m.Type),
		Title:      m.Title,
		Year:       m.Year,
		Genres:     m.Genres,
		PosterURL:  m.PosterUrl,
		CoverURL:   m.CoverUrl,
		ExternalIDs: externalIDs,
		UpdatedAt:  m.UpdatedAt.AsTime(),
	}
}

