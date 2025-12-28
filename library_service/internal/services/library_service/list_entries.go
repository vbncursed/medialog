package library_service

import (
	"context"

	"github.com/vbncursed/medialog/library_service/internal/models"
)

func (s *LibraryService) ListEntries(ctx context.Context, in models.ListEntriesInput) (*models.ListEntriesResult, error) {
	if in.Page == 0 {
		in.Page = 1
	}
	if in.PageSize == 0 {
		in.PageSize = 20
	}
	if in.PageSize > 100 {
		in.PageSize = 100
	}

	if in.SortBy == "" {
		in.SortBy = "created_at"
	}
	if in.SortOrder != "asc" && in.SortOrder != "desc" {
		in.SortOrder = "desc"
	}

	return s.storage.ListEntries(ctx, in)
}
