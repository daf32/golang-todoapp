package auth_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	core_postgres_pool "github.com/daf32/golang-todoapp/internal/core/repository/postgres/pool"
)

func (r *RefreshTokenRepository) RevokeRefreshToken(
	ctx context.Context,
	tokenString string,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE todoapp.refresh_tokens
		SET revoked=true
		WHERE token=$1
	`

	_, err := r.pool.Exec(ctx, query, tokenString)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return fmt.Errorf(
				"refresh token with token=%v: %w",
				tokenString,
				err,
			)
		}

		return fmt.Errorf("scan error: %w", err)
	}

	return nil
}
