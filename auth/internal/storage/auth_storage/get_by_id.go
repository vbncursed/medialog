package auth_storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	pkgerrors "github.com/pkg/errors"
	"github.com/vbncursed/medialog/auth/internal/models"
)

func (s *AuthStorage) GetUserByID(ctx context.Context, userID uint64) (*models.User, error) {
	var u models.User
	err := s.db.QueryRow(ctx, fmt.Sprintf(`
		SELECT %s, %s, %s, %s, %s
		FROM %s
		WHERE %s = $1
	`, idColumn, emailColumn, passwordHashColumn, roleColumn, createdAtColumn, tableName, idColumn),
		userID,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, pkgerrors.Wrap(err, "failed to get user by id")
	}

	return &u, nil
}
