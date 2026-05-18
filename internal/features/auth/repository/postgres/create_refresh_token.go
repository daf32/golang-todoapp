package auth_postgres_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/daf32/golang-todoapp/internal/core/domain"
	"github.com/google/uuid"
)

func (r *RefreshTokenRepository) CreateRefreshToken(
	ctx context.Context,
	userID int,
	ttl time.Duration,
) (domain.RefreshToken, error) {
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

	_, err := r.pool.Exec(
		ctx,
		query,
		token.ID,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
		token.Revoked,
	)
	if err != nil {
		return domain.RefreshToken{}, fmt.Errorf("scan error: %w", err)
	}

	return domain.NewRefreshToken(
		tokenID,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
		token.Revoked,
	), err
}
