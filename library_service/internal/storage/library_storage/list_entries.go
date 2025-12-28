package library_storage

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/vbncursed/medialog/library_service/internal/models"
)

func (s *LibraryStorage) ListEntries(ctx context.Context, input models.ListEntriesInput) (*models.ListEntriesResult, error) {
	// Строим запрос для подсчета общего количества
	countQuery := s.buildListQuery(input, true)
	countSQL, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build count query")
	}

	var total uint32
	err = s.db.QueryRow(ctx, countSQL, countArgs...).Scan(&total)
	if err != nil {
		return nil, errors.Wrap(err, "failed to count entries")
	}

	// Строим запрос для получения записей
	query := s.buildListQuery(input, false)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build list query")
	}

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query entries")
	}
	defer rows.Close()

	var entries []*models.Entry
	for rows.Next() {
		var entry models.Entry
		var startedAt, finishedAt *time.Time

		err := rows.Scan(
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
			return nil, errors.Wrap(err, "failed to scan entry")
		}

		entry.StartedAt = startedAt
		entry.FinishedAt = finishedAt
		entries = append(entries, &entry)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating rows")
	}

	return &models.ListEntriesResult{
		Entries:  entries,
		Total:    total,
		Page:     input.Page,
		PageSize: input.PageSize,
	}, nil
}

func (s *LibraryStorage) buildListQuery(input models.ListEntriesInput, countOnly bool) squirrel.Sqlizer {
	baseQuery := squirrel.Select().
		From(entriesTable).
		Where(squirrel.Eq{"user_id": input.UserID}).
		PlaceholderFormat(squirrel.Dollar)

	if countOnly {
		baseQuery = baseQuery.Columns("COUNT(*)")
	} else {
		baseQuery = baseQuery.Columns(
			"entry_id", "user_id", "media_id", "type", "status",
			"rating", "review", "tags", "started_at", "finished_at",
			"created_at", "updated_at",
		)
	}

	// Фильтры
	if len(input.Types) > 0 {
		baseQuery = baseQuery.Where(squirrel.Eq{"type": input.Types})
	}

	if len(input.Statuses) > 0 {
		baseQuery = baseQuery.Where(squirrel.Eq{"status": input.Statuses})
	}

	if input.MinRating > 0 {
		baseQuery = baseQuery.Where(squirrel.GtOrEq{"rating": input.MinRating})
	}

	if input.MaxRating > 0 && input.MaxRating <= 10 {
		baseQuery = baseQuery.Where(squirrel.LtOrEq{"rating": input.MaxRating})
	}

	if input.FinishedFrom > 0 {
		baseQuery = baseQuery.Where(squirrel.GtOrEq{"finished_at": time.Unix(input.FinishedFrom, 0)})
	}

	if input.FinishedTo > 0 {
		baseQuery = baseQuery.Where(squirrel.LtOrEq{"finished_at": time.Unix(input.FinishedTo, 0)})
	}

	if len(input.Tags) > 0 {
		for _, tag := range input.Tags {
			baseQuery = baseQuery.Where("? = ANY(tags)", tag)
		}
	}

	if !countOnly {
		sortBy := input.SortBy
		if sortBy == "" {
			sortBy = "created_at"
		}

		order := input.SortOrder
		if order != "asc" && order != "desc" {
			order = "desc"
		}

		baseQuery = baseQuery.OrderBy(sortBy + " " + order)

		offset := (input.Page - 1) * input.PageSize
		baseQuery = baseQuery.Limit(uint64(input.PageSize)).Offset(uint64(offset))
	}

	return baseQuery
}
