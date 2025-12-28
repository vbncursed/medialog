package auth_storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type AuthStorage struct {
	db *pgxpool.Pool
}

func NewAuthStorage(connString string) (*AuthStorage, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка парсинга конфига")
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка подключения")
	}

	storage := &AuthStorage{db: db}
	if err := storage.initTables(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *AuthStorage) initTables() error {
	sql := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
  id BIGSERIAL PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS %s (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES %s(id) ON DELETE CASCADE,
  refresh_hash BYTEA UNIQUE NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  user_agent TEXT NOT NULL DEFAULT '',
  ip TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS sessions_user_id_idx ON %s(user_id);
CREATE INDEX IF NOT EXISTS sessions_expires_at_idx ON %s(expires_at);
CREATE INDEX IF NOT EXISTS sessions_revoked_at_idx ON %s(revoked_at) WHERE revoked_at IS NOT NULL;
`, usersTable, sessionsTable, usersTable, sessionsTable, sessionsTable, sessionsTable)

	_, err := s.db.Exec(context.Background(), sql)
	if err != nil {
		return errors.Wrap(err, "init tables")
	}

	return nil
}
