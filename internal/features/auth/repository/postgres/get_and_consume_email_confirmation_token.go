package auth_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	core_auth "github.com/daf32/golang-todoapp/internal/core/auth"
	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	core_postgres_pool "github.com/daf32/golang-todoapp/internal/core/repository/postgres/pool"
)

func (r *AuthRepository) GetAndConsumeEmailConfirmationToken(
	ctx context.Context,
	token string,
) (domain.EmailConfirmationToken, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.BeginTx(ctx)
	if err != nil {
		return domain.EmailConfirmationToken{}, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var result domain.EmailConfirmationToken

	tokenHash := core_auth.HashToken(token)
	result.Token = token

	row := tx.QueryRow(
		ctx,
		`
		SELECT user_id, expires_at
		FROM todoapp.email_confirmation_tokens
		WHERE token_hash = $1
		FOR UPDATE
		`,
		tokenHash,
	)

	if err := row.Scan(&result.UserID, &result.ExpiresAt); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.EmailConfirmationToken{}, fmt.Errorf(
				"token not found: %w", core_errors.ErrNotFound,
			)
		}
		return domain.EmailConfirmationToken{}, fmt.Errorf("scan token: %w", err)
	}

	if result.ExpiresAt.Before(time.Now()) {
		return domain.EmailConfirmationToken{}, fmt.Errorf(
			"token expired: %w", core_errors.ErrInvalidArgument,
		)
	}

	if _, err := tx.Exec(
		ctx,
		`DELETE FROM todoapp.email_confirmation_tokens WHERE token_hash=$1`,
		tokenHash,
	); err != nil {
		return domain.EmailConfirmationToken{}, fmt.Errorf("delete token: %w", err)
	}

	if _, err := tx.Exec(
		ctx,
		`
		UPDATE todoapp.users 
		SET email_verified=true, email_verified_at=now()
		WHERE id=$1
		`,
		result.UserID,
	); err != nil {
		return domain.EmailConfirmationToken{}, fmt.Errorf("commit tx: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.EmailConfirmationToken{}, fmt.Errorf("commit tx: %w", err)
	}

	return result, nil
}
