package library_service_api

import (
	"context"
	"errors"

	"github.com/vbncursed/medialog/library_service/internal/models"
	pb_models "github.com/vbncursed/medialog/library_service/internal/pb/models"
	"github.com/vbncursed/medialog/library_service/internal/services/library_service"
	"google.golang.org/grpc/codes"
)

func (l *LibraryServiceAPI) CreateEntry(ctx context.Context, req *pb_models.CreateEntryRequest) (*pb_models.EntryResponse, error) {
	// Извлекаем user_id из JWT токена (обязательно)
	userID, err := l.getUserIDFromContext(ctx)
	if err != nil {
		return nil, newError(codes.Unauthenticated, ErrCodeUnauthorized, "Authentication required. Invalid or missing JWT token.")
	}

	var startedAt, finishedAt int64
	if req.GetStartedAt() != nil {
		startedAt = req.GetStartedAt().AsTime().Unix()
	}
	if req.GetFinishedAt() != nil {
		finishedAt = req.GetFinishedAt().AsTime().Unix()
	}

	res, err := l.libraryService.CreateEntry(ctx, models.CreateEntryInput{
		UserID:     userID,
		MediaID:    req.GetMediaId(),
		Type:       convertMediaType(req.GetType()),
		Status:     convertEntryStatus(req.GetStatus()),
		Rating:     req.GetRating(),
		Review:     req.GetReview(),
		Tags:       req.GetTags(),
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
	})
	if err != nil {
		switch {
		case errors.Is(err, library_service.ErrInvalidInput):
			return nil, newError(codes.InvalidArgument, ErrCodeInvalidInput, "Invalid input parameters.")
		case errors.Is(err, library_service.ErrEntryAlreadyExists):
			return nil, newError(codes.AlreadyExists, ErrCodeEntryAlreadyExists, "Entry already exists for this media.")
		case errors.Is(err, library_service.ErrInvalidMediaID):
			return nil, newFieldError(codes.InvalidArgument, ErrCodeInvalidMediaID, "media_id", "Invalid media ID.")
		case errors.Is(err, library_service.ErrInvalidMediaType):
			return nil, newFieldError(codes.InvalidArgument, ErrCodeInvalidMediaType, "type", "Invalid media type.")
		case errors.Is(err, library_service.ErrInvalidStatus):
			return nil, newFieldError(codes.InvalidArgument, ErrCodeInvalidStatus, "status", "Invalid entry status.")
		case errors.Is(err, library_service.ErrInvalidRating):
			return nil, newFieldError(codes.InvalidArgument, ErrCodeInvalidRating, "rating", "Rating must be between 0 and 10.")
		default:
			return nil, newError(codes.Internal, ErrCodeInternal, "An internal error occurred. Please try again later.")
		}
	}

	return &pb_models.EntryResponse{
		Entry: convertEntryToProto(res),
	}, nil
}
