package auth_postgres_repository

import (
	"context"
	"fmt"

	"github.com/daf32/golang-todoapp/internal/core/domain"
)

func (r *AuthRepository) CreateUserOAuthIdentity(
	ctx context.Context,
	userID int,
	provider, providerSub, email string,
) (domain.UserOAuthIdentity, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO todoapp.user_oauth_identities (user_id, provider, provider_sub, email)
			VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, provider, provider_sub, email, created_at
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		userID,
		provider,
		providerSub,
		email,
	)

	var identity domain.UserOAuthIdentity

	if err := row.Scan(
		&identity.ID,
		&identity.UserID,
		&identity.Provider,
		&identity.ProviderSub,
		&identity.Email,
		&identity.CreatedAt,
	); err != nil {
		return domain.UserOAuthIdentity{}, fmt.Errorf("scan error: %w", err)
	}

	return identity, nil
}
