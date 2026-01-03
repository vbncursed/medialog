package metadata_storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	pkgerrors "github.com/pkg/errors"
	"github.com/vbncursed/medialog/metadata_service/internal/models"
)

var (
	ErrMediaNotFound = errors.New("media not found")
)

func (s *MetadataStorage) SearchMedia(ctx context.Context, input models.SearchMediaInput) (*models.SearchMediaResult, error) {
	query := s.buildSearchQuery(input)

	// Строим COUNT запрос отдельно
	countQuery := s.buildCountQuery(input)
	countSQL, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to build count query")
	}

	var total uint32
	err = s.db.QueryRow(ctx, countSQL, countArgs...).Scan(&total)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to count media")
	}

	query = query.Limit(uint64(input.PageSize)).Offset(uint64((input.Page - 1) * input.PageSize))
	querySQL, queryArgs, err := query.ToSql()
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to build search query")
	}

	rows, err := s.db.Query(ctx, querySQL, queryArgs...)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to query media")
	}
	defer rows.Close()

	var results []*models.Media
	for rows.Next() {
		media, err := s.scanMedia(ctx, rows)
		if err != nil {
			return nil, pkgerrors.Wrap(err, "failed to scan media")
		}
		results = append(results, media)
	}

	return &models.SearchMediaResult{
		Results:  results,
		Total:    total,
		Page:     input.Page,
		PageSize: input.PageSize,
	}, nil
}

func (s *MetadataStorage) buildSearchQuery(input models.SearchMediaInput) squirrel.SelectBuilder {
	query := squirrel.Select(
		fmt.Sprintf("m.%s", mediaIDColumn),
		fmt.Sprintf("m.%s", typeColumn),
		fmt.Sprintf("m.%s", titleColumn),
		fmt.Sprintf("m.%s", yearColumn),
		fmt.Sprintf("m.%s", genresColumn),
		fmt.Sprintf("m.%s", posterURLColumn),
		fmt.Sprintf("m.%s", coverURLColumn),
		fmt.Sprintf("m.%s", updatedAtColumn),
	).
		From(fmt.Sprintf("%s m", metadataMediaTable)).
		PlaceholderFormat(squirrel.Dollar)

	if input.Query != "" {
		query = query.Where(squirrel.ILike{fmt.Sprintf("m.%s", titleColumn): fmt.Sprintf("%%%s%%", input.Query)})
	}

	if input.Type != nil {
		query = query.Where(squirrel.Eq{fmt.Sprintf("m.%s", typeColumn): int(*input.Type)})
	}

	if input.ExternalID != nil {
		query = query.Join(fmt.Sprintf("%s e ON m.%s = e.%s", metadataExternalIDsTable, mediaIDColumn, mediaIDColumn)).
			Where(squirrel.Eq{
				fmt.Sprintf("e.%s", sourceColumn):     input.ExternalID.Source,
				fmt.Sprintf("e.%s", externalIDColumn): input.ExternalID.ExternalID,
			})
	}

	query = query.OrderBy(fmt.Sprintf("m.%s DESC", updatedAtColumn))

	return query
}

func (s *MetadataStorage) buildCountQuery(input models.SearchMediaInput) squirrel.SelectBuilder {
	query := squirrel.Select("COUNT(*)").
		From(fmt.Sprintf("%s m", metadataMediaTable)).
		PlaceholderFormat(squirrel.Dollar)

	if input.Query != "" {
		query = query.Where(squirrel.ILike{fmt.Sprintf("m.%s", titleColumn): fmt.Sprintf("%%%s%%", input.Query)})
	}

	if input.Type != nil {
		query = query.Where(squirrel.Eq{fmt.Sprintf("m.%s", typeColumn): int(*input.Type)})
	}

	if input.ExternalID != nil {
		query = query.Join(fmt.Sprintf("%s e ON m.%s = e.%s", metadataExternalIDsTable, mediaIDColumn, mediaIDColumn)).
			Where(squirrel.Eq{
				fmt.Sprintf("e.%s", sourceColumn):     input.ExternalID.Source,
				fmt.Sprintf("e.%s", externalIDColumn): input.ExternalID.ExternalID,
			})
	}

	return query
}
