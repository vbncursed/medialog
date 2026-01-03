package metadata_service_api

import (
	"context"
	"errors"

	pb_models "github.com/vbncursed/medialog/metadata_service/internal/pb/models"
	"github.com/vbncursed/medialog/metadata_service/internal/services/metadata_service"
	"google.golang.org/grpc/codes"
)

func (m *MetadataServiceAPI) GetMedia(ctx context.Context, req *pb_models.GetMediaRequest) (*pb_models.GetMediaResponse, error) {
	if req.MediaId == 0 {
		return nil, newFieldError(codes.InvalidArgument, ErrCodeInvalidMediaID, "media_id", "Media ID is required.")
	}

	media, err := m.metadataService.GetMedia(ctx, req.MediaId)
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

