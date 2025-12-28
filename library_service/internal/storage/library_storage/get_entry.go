package library_storage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	pkgerrors "github.com/pkg/errors"
	"github.com/vbncursed/medialog/library_service/internal/models"
)

var (
	ErrEntryNotFound = errors.New("entry not found")
)

func (s *LibraryStorage) GetEntry(ctx context.Context, entryID, userID uint64) (*models.Entry, error) {
	row := s.db.QueryRow(ctx,
		`SELECT entry_id, user_id, media_id, type, status, rating, review, tags, 
		started_at, finished_at, created_at, updated_at 
		FROM `+entriesTable+` 
		WHERE entry_id = $1 AND user_id = $2`,
		entryID, userID,
	)

	var entry models.Entry
	var startedAt, finishedAt *time.Time

	err := row.Scan(
		&entry.EntryID,
		&entry.UserID,
		&entry.MediaID,
		&entry.Type,
		&entry.Status,
		&entry.Rating,
		&entry.Review,
		&entry.Tags,
		&startedAt,
		&finishedAt,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEntryNotFound
		}
		return nil, pkgerrors.Wrap(err, "failed to get entry")
	}

	entry.StartedAt = startedAt
	entry.FinishedAt = finishedAt

	return &entry, nil
}
