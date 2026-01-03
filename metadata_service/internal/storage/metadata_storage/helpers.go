package metadata_storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/vbncursed/medialog/metadata_service/internal/models"
)

func (s *MetadataStorage) scanMediaRow(row pgx.Row) (*models.Media, error) {
	var media models.Media
	var year sql.NullInt32
	var posterURL, coverURL sql.NullString

	err := row.Scan(
		&media.MediaID,
		&media.Type,
		&media.Title,
		&year,
		&media.Genres,
		&posterURL,
		&coverURL,
		&media.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if year.Valid {
		y := uint32(year.Int32)
		media.Year = &y
	}

	if posterURL.Valid {
		media.PosterURL = &posterURL.String
	}

	if coverURL.Valid {
		media.CoverURL = &coverURL.String
	}

	return &media, nil
}

func (s *MetadataStorage) scanMedia(ctx context.Context, rows pgx.Rows) (*models.Media, error) {
	var media models.Media
	var year sql.NullInt32
	var posterURL, coverURL sql.NullString

	err := rows.Scan(
		&media.MediaID,
		&media.Type,
		&media.Title,
		&year,
		&media.Genres,
		&posterURL,
		&coverURL,
		&media.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if year.Valid {
		y := uint32(year.Int32)
		media.Year = &y
	}

	if posterURL.Valid {
		media.PosterURL = &posterURL.String
	}

	if coverURL.Valid {
		media.CoverURL = &coverURL.String
	}

	externalIDs, err := s.getExternalIDs(ctx, media.MediaID)
	if err != nil {
		return nil, fmt.Errorf("failed to get external ids: %w", err)
	}
	media.ExternalIDs = externalIDs

	return &media, nil
}

func (s *MetadataStorage) getExternalIDs(ctx context.Context, mediaID uint64) ([]models.ExternalID, error) {
	rows, err := s.db.Query(ctx,
		fmt.Sprintf(`SELECT %s, %s
		FROM %s
		WHERE %s = $1`,
			sourceColumn, externalIDColumn,
			metadataExternalIDsTable,
			mediaIDColumn),
		mediaID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var externalIDs []models.ExternalID
	for rows.Next() {
		var extID models.ExternalID
		if err := rows.Scan(&extID.Source, &extID.ExternalID); err != nil {
			return nil, err
		}
		externalIDs = append(externalIDs, extID)
	}

	return externalIDs, nil
}
