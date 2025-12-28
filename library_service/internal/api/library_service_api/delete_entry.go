package library_service_api

import (
	"context"
	"errors"

	pb_models "github.com/vbncursed/medialog/library_service/internal/pb/models"
	"github.com/vbncursed/medialog/library_service/internal/services/library_service"
	"google.golang.org/grpc/codes"
)

func (l *LibraryServiceAPI) DeleteEntry(ctx context.Context, req *pb_models.DeleteEntryRequest) (*pb_models.DeleteEntryResponse, error) {
	// Извлекаем user_id из JWT токена (обязательно)
	userID, err := l.getUserIDFromContext(ctx)
	if err != nil {
		return nil, newError(codes.Unauthenticated, ErrCodeUnauthorized, "Authentication required. Invalid or missing JWT token.")
	}

	err = l.libraryService.DeleteEntry(ctx, req.GetEntryId(), userID)
	if err != nil {
		switch {
		case errors.Is(err, library_service.ErrEntryNotFound):
			return nil, newError(codes.NotFound, ErrCodeEntryNotFound, "Entry not found.")
		case errors.Is(err, library_service.ErrUnauthorized):
			return nil, newError(codes.PermissionDenied, ErrCodeUnauthorized, "You don't have permission to delete this entry.")
		default:
			return nil, newError(codes.Internal, ErrCodeInternal, "An internal error occurred. Please try again later.")
		}
	}

	return &pb_models.DeleteEntryResponse{}, nil
}
