package auth_storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/vbncursed/medialog/auth_service/internal/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

func (s *AuthStorage) CreateUser(ctx context.Context, email string, passwordHash string) (uint64, error) {
	row := s.db.QueryRow(ctx,
		`INSERT INTO `+usersTable+` (email, password_hash, role) VALUES ($1, $2, $3) RETURNING id`,
		email, passwordHash, models.RoleUser,
	)

	var id uint64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (s *AuthStorage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	row := s.db.QueryRow(ctx,
		`SELECT id, email, password_hash, role, created_at FROM `+usersTable+` WHERE email = $1`,
		email,
	)

	var u models.User
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (s *AuthStorage) GetUserByID(ctx context.Context, userID uint64) (*models.User, error) {
	row := s.db.QueryRow(ctx,
		`SELECT id, email, password_hash, role, created_at FROM `+usersTable+` WHERE id = $1`,
		userID,
	)

	var u models.User
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}
