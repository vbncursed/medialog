package metadata_storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	pkgerrors "github.com/pkg/errors"
	"github.com/vbncursed/medialog/metadata_service/internal/models"
)

func (s *MetadataStorage) GetMediaByExternalID(ctx context.Context, source, externalID string) (*models.Media, error) {
	row := s.db.QueryRow(ctx,
		fmt.Sprintf(`SELECT m.%s, m.%s, m.%s, m.%s, m.%s, m.%s, m.%s, m.%s
		FROM %s m
		JOIN %s e ON m.%s = e.%s
		WHERE e.%s = $1 AND e.%s = $2`,
			mediaIDColumn, typeColumn, titleColumn, yearColumn, genresColumn,
			posterURLColumn, coverURLColumn, updatedAtColumn,
			metadataMediaTable, metadataExternalIDsTable,
			mediaIDColumn, mediaIDColumn,
			sourceColumn, externalIDColumn),
		source, externalID,
	)

	media, err := s.scanMediaRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMediaNotFound
		}
		return nil, pkgerrors.Wrap(err, "failed to get media by external id")
	}

	externalIDs, err := s.getExternalIDs(ctx, media.MediaID)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to get external ids")
	}
	media.ExternalIDs = externalIDs

	return media, nil
}
