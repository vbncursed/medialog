package metadata_storage

import (
	"context"
	"fmt"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/vbncursed/medialog/metadata_service/internal/models"
)

func (s *MetadataStorage) CreateMedia(ctx context.Context, media *models.Media) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return pkgerrors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	var mediaID uint64
	now := time.Now()
	err = tx.QueryRow(ctx,
		fmt.Sprintf(`INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING %s`,
			metadataMediaTable, typeColumn, titleColumn, yearColumn, genresColumn,
			posterURLColumn, coverURLColumn, updatedAtColumn, mediaIDColumn),
		int(media.Type), media.Title, media.Year, media.Genres,
		media.PosterURL, media.CoverURL, now,
	).Scan(&mediaID)

	if err != nil {
		return pkgerrors.Wrap(err, "failed to insert media")
	}

	media.MediaID = mediaID

	if len(media.ExternalIDs) > 0 {
		for _, extID := range media.ExternalIDs {
			_, err = tx.Exec(ctx,
				fmt.Sprintf(`INSERT INTO %s (%s, %s, %s)
				VALUES ($1, $2, $3)
				ON CONFLICT (%s, %s, %s) DO NOTHING`,
					metadataExternalIDsTable, mediaIDColumn, sourceColumn, externalIDColumn,
					mediaIDColumn, sourceColumn, externalIDColumn),
				mediaID, extID.Source, extID.ExternalID,
			)
			if err != nil {
				return pkgerrors.Wrap(err, "failed to insert external id")
			}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return pkgerrors.Wrap(err, "failed to commit transaction")
	}

	return nil
}
