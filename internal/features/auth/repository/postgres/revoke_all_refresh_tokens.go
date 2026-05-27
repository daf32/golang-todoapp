package auth_postgres_repository

import (
	"context"
	"fmt"
)

func (r *AuthRepository) RevokeAllRefreshTokensForUser(
	ctx context.Context,
	userID int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE todoapp.refresh_tokens
		SET revoked = true
		WHERE user_id = $1 AND revoked = false
	`

	if _, err := r.pool.Exec(
		ctx,
		query,
		userID,
	); err != nil {
		return fmt.Errorf("revoke all refresh tokens for user_id=%d: %w", userID, err)
	}

	return nil
}
