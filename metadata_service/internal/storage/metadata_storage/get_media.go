package metadata_storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	pkgerrors "github.com/pkg/errors"
	"github.com/vbncursed/medialog/metadata_service/internal/models"
)

func (s *MetadataStorage) GetMedia(ctx context.Context, mediaID uint64) (*models.Media, error) {
	row := s.db.QueryRow(ctx,
		fmt.Sprintf(`SELECT m.%s, m.%s, m.%s, m.%s, m.%s, m.%s, m.%s, m.%s
		FROM %s m
		WHERE m.%s = $1`,
			mediaIDColumn, typeColumn, titleColumn, yearColumn, genresColumn,
			posterURLColumn, coverURLColumn, updatedAtColumn,
			metadataMediaTable, mediaIDColumn),
		mediaID,
	)

	media, err := s.scanMediaRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMediaNotFound
		}
		return nil, pkgerrors.Wrap(err, "failed to get media")
	}

	externalIDs, err := s.getExternalIDs(ctx, mediaID)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to get external ids")
	}
	media.ExternalIDs = externalIDs

	return media, nil
}
