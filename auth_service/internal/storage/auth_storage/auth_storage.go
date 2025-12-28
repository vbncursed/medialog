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
`, usersTable)

	_, err := s.db.Exec(context.Background(), sql)
	if err != nil {
		return errors.Wrap(err, "init tables")
	}

	return nil
}
