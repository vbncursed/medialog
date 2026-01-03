package metadata_storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type MetadataStorage struct {
	db *pgxpool.Pool
}

func NewMetadataStorage(connString string) (*MetadataStorage, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка парсинга конфига")
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка подключения")
	}

	storage := &MetadataStorage{db: db}
	if err := storage.initTables(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *MetadataStorage) initTables() error {
	sql := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
  media_id BIGSERIAL PRIMARY KEY,
  type INT NOT NULL,
  title TEXT NOT NULL,
  year INT,
  genres TEXT[] NOT NULL DEFAULT '{}',
  poster_url TEXT,
  cover_url TEXT,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS %s (
  media_id BIGINT NOT NULL REFERENCES %s(media_id) ON DELETE CASCADE,
  source TEXT NOT NULL,
  external_id TEXT NOT NULL,
  PRIMARY KEY (media_id, source, external_id)
);

CREATE INDEX IF NOT EXISTS idx_%s_type ON %s(type);
CREATE INDEX IF NOT EXISTS idx_%s_title ON %s(title);
CREATE INDEX IF NOT EXISTS idx_%s_year ON %s(year);
CREATE INDEX IF NOT EXISTS idx_%s_source_external_id ON %s(source, external_id);
`, metadataMediaTable, metadataExternalIDsTable, metadataMediaTable,
		metadataMediaTable, metadataMediaTable,
		metadataMediaTable, metadataMediaTable,
		metadataMediaTable, metadataMediaTable,
		metadataExternalIDsTable, metadataExternalIDsTable)

	_, err := s.db.Exec(context.Background(), sql)
	if err != nil {
		return errors.Wrap(err, "init tables")
	}

	return nil
}
