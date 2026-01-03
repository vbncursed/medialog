package metadata_service_api

import (
	"context"
	"errors"

	"github.com/vbncursed/medialog/metadata_service/internal/models"
	pb_models "github.com/vbncursed/medialog/metadata_service/internal/pb/models"
	"github.com/vbncursed/medialog/metadata_service/internal/services/metadata_service"
	"google.golang.org/grpc/codes"
)

func (m *MetadataServiceAPI) CreateMedia(ctx context.Context, req *pb_models.CreateMediaRequest) (*pb_models.GetMediaResponse, error) {
	if req.Title == "" {
		return nil, newFieldError(codes.InvalidArgument, ErrCodeMissingField, "title", "Title is required.")
	}

	externalIDs := make([]models.ExternalID, 0, len(req.ExternalIds))
	for _, eid := range req.ExternalIds {
		if extID := convertExternalIDFromProto(eid); extID != nil {
			externalIDs = append(externalIDs, *extID)
		}
	}

	input := models.CreateMediaInput{
		Type:        convertMediaType(req.Type),
		Title:       req.Title,
		Year:        req.Year,
		Genres:      req.Genres,
		PosterURL:   req.PosterUrl,
		CoverURL:    req.CoverUrl,
		ExternalIDs: externalIDs,
	}

	media, err := m.metadataService.CreateMedia(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, metadata_service.ErrInvalidInput):
			return nil, newError(codes.InvalidArgument, ErrCodeInvalidInput, "Invalid input parameters.")
		case errors.Is(err, metadata_service.ErrInvalidMediaType):
			return nil, newFieldError(codes.InvalidArgument, ErrCodeInvalidMediaType, "type", "Invalid media type.")
		default:
			return nil, newError(codes.Internal, ErrCodeInternal, "An internal error occurred. Please try again later.")
		}
	}

	return &pb_models.GetMediaResponse{
		Media: convertMediaToProto(media),
	}, nil
}

