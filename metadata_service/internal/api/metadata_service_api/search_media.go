package metadata_service_api

import (
	"context"
	"errors"

	"github.com/vbncursed/medialog/metadata_service/internal/models"
	pb_models "github.com/vbncursed/medialog/metadata_service/internal/pb/models"
	"github.com/vbncursed/medialog/metadata_service/internal/services/metadata_service"
	"google.golang.org/grpc/codes"
)

func (m *MetadataServiceAPI) SearchMedia(ctx context.Context, req *pb_models.SearchMediaRequest) (*pb_models.SearchMediaResponse, error) {
	var mediaType *models.MediaType
	if req.Type != nil {
		mt := convertMediaType(*req.Type)
		mediaType = &mt
	}

	var externalID *models.ExternalID
	if req.ExternalId != nil {
		externalID = convertExternalIDFromProto(req.ExternalId)
	}

	input := models.SearchMediaInput{
		Query:      req.Query,
		Type:       mediaType,
		ExternalID: externalID,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	res, err := m.metadataService.SearchMedia(ctx, input)
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

	results := make([]*pb_models.Media, 0, len(res.Results))
	for _, media := range res.Results {
		results = append(results, convertMediaToProto(media))
	}

	totalPages := uint32(0)
	if res.Total > 0 && res.PageSize > 0 {
		totalPages = (res.Total + res.PageSize - 1) / res.PageSize
	}

	return &pb_models.SearchMediaResponse{
		Results:    results,
		Total:      res.Total,
		Page:       res.Page,
		PageSize:   res.PageSize,
		TotalPages: totalPages,
	}, nil
}

