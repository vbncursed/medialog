package library_storage

import (
	"context"

	"github.com/pkg/errors"
)

func (s *LibraryStorage) DeleteEntry(ctx context.Context, entryID, userID uint64) error {
	result, err := s.db.Exec(ctx,
		`DELETE FROM `+entriesTable+` 
		WHERE entry_id = $1 AND user_id = $2`,
		entryID, userID,
	)

	if err != nil {
		return errors.Wrap(err, "failed to delete entry")
	}

	if result.RowsAffected() == 0 {
		return ErrEntryNotFound
	}

	return nil
}

