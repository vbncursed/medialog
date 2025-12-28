package library_service_api

import (
	"context"
	"errors"

	"github.com/vbncursed/medialog/library_service/internal/models"
	pb_models "github.com/vbncursed/medialog/library_service/internal/pb/models"
	"github.com/vbncursed/medialog/library_service/internal/services/library_service"
	"google.golang.org/grpc/codes"
)

func (l *LibraryServiceAPI) ListEntries(ctx context.Context, req *pb_models.ListEntriesRequest) (*pb_models.ListEntriesResponse, error) {
	// Извлекаем user_id из JWT токена (обязательно)
	userID, err := l.getUserIDFromContext(ctx)
	if err != nil {
		return nil, newError(codes.Unauthenticated, ErrCodeUnauthorized, "Authentication required. Invalid or missing JWT token.")
	}

	// Конвертируем типы и статусы
	types := make([]models.MediaType, 0, len(req.GetTypes()))
	for _, t := range req.GetTypes() {
		types = append(types, convertMediaType(t))
	}

	statuses := make([]models.EntryStatus, 0, len(req.GetStatuses()))
	for _, s := range req.GetStatuses() {
		statuses = append(statuses, convertEntryStatus(s))
	}

	var finishedFrom, finishedTo int64
	if req.FinishedFrom != nil && req.GetFinishedFrom() != nil {
		finishedFrom = req.GetFinishedFrom().AsTime().Unix()
	}
	if req.FinishedTo != nil && req.GetFinishedTo() != nil {
		finishedTo = req.GetFinishedTo().AsTime().Unix()
	}

	input := models.ListEntriesInput{
		UserID:       userID,
		Types:        types,
		Statuses:     statuses,
		Tags:         req.GetTags(),
		MinRating:    req.GetMinRating(),
		MaxRating:    req.GetMaxRating(),
		FinishedFrom: finishedFrom,
		FinishedTo:   finishedTo,
		SortBy:       req.GetSortBy(),
		SortOrder:    req.GetSortOrder(),
		Page:         req.GetPage(),
		PageSize:     req.GetPageSize(),
	}

	res, err := l.libraryService.ListEntries(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, library_service.ErrInvalidInput):
			return nil, newError(codes.InvalidArgument, ErrCodeInvalidInput, "Invalid input parameters.")
		default:
			return nil, newError(codes.Internal, ErrCodeInternal, "An internal error occurred. Please try again later.")
		}
	}

	// Конвертируем результаты
	entries := make([]*pb_models.Entry, 0, len(res.Entries))
	for _, e := range res.Entries {
		entries = append(entries, convertEntryToProto(e))
	}

	// Вычисляем total_pages
	totalPages := uint32(0)
	if res.Total > 0 && res.PageSize > 0 {
		totalPages = (res.Total + res.PageSize - 1) / res.PageSize
	}

	return &pb_models.ListEntriesResponse{
		Entries:    entries,
		Total:      res.Total,
		Page:       res.Page,
		PageSize:   res.PageSize,
		TotalPages: totalPages,
	}, nil
}
