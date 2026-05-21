package auth_postgres_repository

import (
	"context"
	"fmt"
	"time"

	core_auth "github.com/daf32/golang-todoapp/internal/core/auth"
	"github.com/google/uuid"
)

func (r *AuthRepository) CreateRefreshToken(
	ctx context.Context,
	userID int,
	ttl time.Duration,
) (core_auth.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tokenID := uuid.New()
	expiresAt := time.Now().Add(ttl)

	token := &RefreshToken{
		ID:        tokenID,
		UserID:    userID,
		Token:     tokenID.String(),
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		Revoked:   false,
	}

	query := `
        INSERT INTO todoapp.refresh_tokens (id, user_id, token, expires_at, created_at, revoked)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

	if _, err := r.pool.Exec(
		ctx,
		query,
		token.ID,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
		token.Revoked,
	); err != nil {
		return core_auth.RefreshToken{}, fmt.Errorf("scan error: %w", err)
	}

	return core_auth.NewRefreshToken(
		tokenID,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
		token.Revoked,
	), nil
}
