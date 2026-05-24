package auth_postgres_repository

import (
	"context"
	"fmt"
	"time"

	core_auth "github.com/daf32/golang-todoapp/internal/core/auth"
	"github.com/daf32/golang-todoapp/internal/core/domain"
)

func (r *AuthRepository) CreateEmailConfirmationToken(
	ctx context.Context,
	userID int,
	ttl time.Duration,
) (domain.EmailConfirmationToken, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	plainText, err := core_auth.GenerateSecureToken(32)
	if err != nil {
		return domain.EmailConfirmationToken{}, fmt.Errorf("generate token: %w", err)
	}

	tokenHash := core_auth.HashToken(plainText)

	expiresAt := time.Now().Add(ttl)

	if _, err := r.pool.Exec(
		ctx,
		`
		DELETE FROM todoapp.email_confirmation_tokens 
		WHERE user_id = $1 OR expires_at < now()
		`,
		userID,
	); err != nil {
		return domain.EmailConfirmationToken{}, fmt.Errorf("cleanup tokens: %w", err)
	}

	query := `
		INSERT INTO todoapp.email_confirmation_tokens (token_hash, user_id, expires_at)
		VALUES ($1, $2, $3)
	`

	if _, err := r.pool.Exec(
		ctx,
		query,
		tokenHash,
		userID,
		expiresAt,
	); err != nil {
		return domain.EmailConfirmationToken{}, fmt.Errorf("insert email confirmation token: %w", err)
	}

	return domain.NewEmailConfirmationToken(
		plainText,
		userID,
		expiresAt,
	), nil
}
