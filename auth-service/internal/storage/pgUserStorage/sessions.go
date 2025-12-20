package pgstorage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/vbncursed/medialog/auth-service/internal/models"
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

func (s *PGstorage) CreateSession(ctx context.Context, userID uint64, refreshHash []byte, expiresAt time.Time, userAgent, ip string) (uint64, error) {
	row := s.db.QueryRow(ctx,
		`INSERT INTO `+sessionsTable+` (user_id, refresh_hash, expires_at, user_agent, ip)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		userID, refreshHash, expiresAt, userAgent, ip,
	)

	var id uint64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (s *PGstorage) GetSessionByRefreshHash(ctx context.Context, refreshHash []byte) (*models.Session, error) {
	row := s.db.QueryRow(ctx,
		`SELECT id, user_id, refresh_hash, expires_at, revoked_at, created_at, user_agent, ip
		 FROM `+sessionsTable+`
		 WHERE refresh_hash = $1`,
		refreshHash,
	)

	var sess models.Session
	if err := row.Scan(
		&sess.ID,
		&sess.UserID,
		&sess.RefreshHash,
		&sess.ExpiresAt,
		&sess.RevokedAt,
		&sess.CreatedAt,
		&sess.UserAgent,
		&sess.IP,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	return &sess, nil
}

func (s *PGstorage) RevokeSessionByID(ctx context.Context, sessionID uint64, revokedAt time.Time) error {
	ct, err := s.db.Exec(ctx,
		`UPDATE `+sessionsTable+` SET revoked_at = $1 WHERE id = $2 AND revoked_at IS NULL`,
		revokedAt, sessionID,
	)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrSessionNotFound
	}
	return nil
}

func (s *PGstorage) RevokeAllSessionsByUserID(ctx context.Context, userID uint64, revokedAt time.Time) error {
	_, err := s.db.Exec(ctx,
		`UPDATE `+sessionsTable+` SET revoked_at = $1 WHERE user_id = $2 AND revoked_at IS NULL`,
		revokedAt, userID,
	)
	return err
}
