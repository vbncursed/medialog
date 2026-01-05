package auth_storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	pkgerrors "github.com/pkg/errors"
	"github.com/vbncursed/medialog/auth/internal/models"
)

func (s *AuthStorage) CreateUser(ctx context.Context, email string, passwordHash string) (uint64, error) {
	var userID uint64
	err := s.db.QueryRow(ctx, fmt.Sprintf(`
		INSERT INTO %s (%s, %s, %s)
		VALUES ($1, $2, $3)
		RETURNING %s
	`, tableName, emailColumn, passwordHashColumn, roleColumn, idColumn),
		email, passwordHash, models.RoleUser,
	).Scan(&userID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, pkgerrors.Wrap(err, "failed to create user")
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, pkgerrors.Wrap(err, "email already exists")
		}
		return 0, pkgerrors.Wrap(err, "failed to create user")
	}

	return userID, nil
}
