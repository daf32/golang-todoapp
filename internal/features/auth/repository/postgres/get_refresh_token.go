package auth_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	core_postgres_pool "github.com/daf32/golang-todoapp/internal/core/repository/postgres/pool"
)

func (r *RefreshTokenRepository) GetRefreshToken(
	ctx context.Context,
	tokenString string,
) (domain.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, user_id, token, expires_at, created_at, revoked
		FROM todoapp.refresh_tokens
		WHERE token=$1
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		tokenString,
	)

	var refreshTokenModel RefreshToken

	err := row.Scan(
		&refreshTokenModel.ID,
		&refreshTokenModel.UserID,
		&refreshTokenModel.Token,
		&refreshTokenModel.ExpiresAt,
		&refreshTokenModel.CreatedAt,
		&refreshTokenModel.Revoked,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.RefreshToken{}, fmt.Errorf(
				"refresh token with token=%v: %w",
				tokenString,
				core_errors.ErrNotFound,
			)
		}

		return domain.RefreshToken{}, fmt.Errorf("scan error: %w", err)
	}

	return domain.NewRefreshToken(
		refreshTokenModel.ID,
		refreshTokenModel.UserID,
		refreshTokenModel.Token,
		refreshTokenModel.ExpiresAt,
		refreshTokenModel.CreatedAt,
		refreshTokenModel.Revoked,
	), nil
}
