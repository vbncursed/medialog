package library_storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type LibraryStorage struct {
	db *pgxpool.Pool
}

func NewLibraryStorage(connString string) (*LibraryStorage, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка парсинга конфига")
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка подключения")
	}

	storage := &LibraryStorage{db: db}
	if err := storage.initTables(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *LibraryStorage) initTables() error {
	sql := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
  entry_id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  media_id BIGINT NOT NULL,
  type INT NOT NULL,
  status INT NOT NULL,
  rating INT NOT NULL DEFAULT 0,
  review TEXT NOT NULL DEFAULT '',
  tags TEXT[] NOT NULL DEFAULT '{}',
  started_at TIMESTAMPTZ,
  finished_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(user_id, media_id, type)
);

CREATE INDEX IF NOT EXISTS idx_%s_user_id ON %s(user_id);
CREATE INDEX IF NOT EXISTS idx_%s_media_id ON %s(media_id);
CREATE INDEX IF NOT EXISTS idx_%s_type ON %s(type);
CREATE INDEX IF NOT EXISTS idx_%s_status ON %s(status);
CREATE INDEX IF NOT EXISTS idx_%s_rating ON %s(rating);
CREATE INDEX IF NOT EXISTS idx_%s_finished_at ON %s(finished_at);
`, entriesTable, entriesTable, entriesTable, entriesTable, entriesTable, entriesTable, entriesTable, entriesTable, entriesTable, entriesTable, entriesTable, entriesTable, entriesTable)

	_, err := s.db.Exec(context.Background(), sql)
	if err != nil {
		return errors.Wrap(err, "init tables")
	}

	return nil
}
