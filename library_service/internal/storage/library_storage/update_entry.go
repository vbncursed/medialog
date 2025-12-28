package library_storage

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/vbncursed/medialog/library_service/internal/models"
)

func (s *LibraryStorage) UpdateEntry(ctx context.Context, entry *models.Entry) error {
	var startedAt, finishedAt *time.Time
	if entry.StartedAt != nil {
		startedAt = entry.StartedAt
	}
	if entry.FinishedAt != nil {
		finishedAt = entry.FinishedAt
	}

	result, err := s.db.Exec(ctx,
		`UPDATE `+entriesTable+` 
		SET status = $1, rating = $2, review = $3, tags = $4, 
		started_at = $5, finished_at = $6, updated_at = $7
		WHERE entry_id = $8 AND user_id = $9`,
		entry.Status,
		entry.Rating,
		entry.Review,
		entry.Tags,
		startedAt,
		finishedAt,
		entry.UpdatedAt,
		entry.EntryID,
		entry.UserID,
	)

	if err != nil {
		return errors.Wrap(err, "failed to update entry")
	}

	if result.RowsAffected() == 0 {
		return ErrEntryNotFound
	}

	return nil
}
