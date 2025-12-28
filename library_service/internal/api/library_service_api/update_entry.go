package library_service_api

import (
	"context"
	"errors"

	"github.com/vbncursed/medialog/library_service/internal/models"
	pb_models "github.com/vbncursed/medialog/library_service/internal/pb/models"
	"github.com/vbncursed/medialog/library_service/internal/services/library_service"
	"google.golang.org/grpc/codes"
)

func (l *LibraryServiceAPI) UpdateEntry(ctx context.Context, req *pb_models.UpdateEntryRequest) (*pb_models.EntryResponse, error) {
	// Извлекаем user_id из JWT токена (обязательно)
	userID, err := l.getUserIDFromContext(ctx)
	if err != nil {
		return nil, newError(codes.Unauthenticated, ErrCodeUnauthorized, "Authentication required. Invalid or missing JWT token.")
	}

	input := models.UpdateEntryInput{
		EntryID: req.GetEntryId(),
		UserID:  userID,
	}

	if req.Status != nil {
		status := convertEntryStatus(*req.Status)
		input.Status = &status
	}
	if req.Rating != nil {
		rating := req.GetRating()
		input.Rating = &rating
	}
	if req.Review != nil {
		review := req.GetReview()
		input.Review = &review
	}
	if len(req.GetTags()) > 0 {
		input.Tags = req.GetTags()
	}
	if req.StartedAt != nil && req.GetStartedAt() != nil {
		startedAt := req.GetStartedAt().AsTime().Unix()
		input.StartedAt = &startedAt
	}
	if req.FinishedAt != nil && req.GetFinishedAt() != nil {
		finishedAt := req.GetFinishedAt().AsTime().Unix()
		input.FinishedAt = &finishedAt
	}

	res, err := l.libraryService.UpdateEntry(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, library_service.ErrEntryNotFound):
			return nil, newError(codes.NotFound, ErrCodeEntryNotFound, "Entry not found.")
		case errors.Is(err, library_service.ErrUnauthorized):
			return nil, newError(codes.PermissionDenied, ErrCodeUnauthorized, "You don't have permission to update this entry.")
		case errors.Is(err, library_service.ErrInvalidInput):
			return nil, newError(codes.InvalidArgument, ErrCodeInvalidInput, "Invalid input parameters.")
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
