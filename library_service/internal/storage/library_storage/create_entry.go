package library_storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"github.com/vbncursed/medialog/library_service/internal/models"
)

func (s *LibraryStorage) CreateEntry(ctx context.Context, entry *models.Entry) error {
	var startedAt, finishedAt *time.Time
	if entry.StartedAt != nil {
		startedAt = entry.StartedAt
	}
	if entry.FinishedAt != nil {
		finishedAt = entry.FinishedAt
	}

	tags := entry.Tags
	if tags == nil {
		tags = []string{}
	}

	row := s.db.QueryRow(ctx,
		`INSERT INTO `+entriesTable+` 
		(user_id, media_id, type, status, rating, review, tags, started_at, finished_at, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
		RETURNING entry_id`,
		entry.UserID,
		entry.MediaID,
		entry.Type,
		entry.Status,
		entry.Rating,
		entry.Review,
		tags,
		startedAt,
		finishedAt,
		entry.CreatedAt,
		entry.UpdatedAt,
	)

	var entryID uint64
	if err := row.Scan(&entryID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return ErrEntryAlreadyExists
			}
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.Wrap(err, "failed to get entry_id after insert")
		}
		return errors.Wrap(err, "failed to create entry")
	}

	entry.EntryID = entryID
	return nil
}
