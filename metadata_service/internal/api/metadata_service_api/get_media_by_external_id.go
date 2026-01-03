package metadata_service_api

import (
	"context"
	"errors"

	pb_models "github.com/vbncursed/medialog/metadata_service/internal/pb/models"
	"github.com/vbncursed/medialog/metadata_service/internal/services/metadata_service"
	"google.golang.org/grpc/codes"
)

func (m *MetadataServiceAPI) GetMediaByExternalID(ctx context.Context, req *pb_models.GetMediaByExternalIDRequest) (*pb_models.GetMediaResponse, error) {
	if req.Source == "" {
		return nil, newFieldError(codes.InvalidArgument, ErrCodeInvalidSource, "source", "Source is required.")
	}
	if req.ExternalId == "" {
		return nil, newFieldError(codes.InvalidArgument, ErrCodeMissingField, "external_id", "External ID is required.")
	}

	media, err := m.metadataService.GetMediaByExternalID(ctx, req.Source, req.ExternalId)
	if err != nil {
		switch {
		case errors.Is(err, metadata_service.ErrMediaNotFound):
			return nil, newError(codes.NotFound, ErrCodeMediaNotFound, "Media not found.")
		default:
			return nil, newError(codes.Internal, ErrCodeInternal, "An internal error occurred. Please try again later.")
		}
	}

	return &pb_models.GetMediaResponse{
		Media: convertMediaToProto(media),
	}, nil
}

