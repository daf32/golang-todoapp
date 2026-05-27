package auth_postgres_repository

import (
	"context"
	"fmt"

	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	core_postgres_pool "github.com/daf32/golang-todoapp/internal/core/repository/postgres/pool"
)

func (r *AuthRepository) GetUserOAuthIdentity(
	ctx context.Context,
	provider, providerSub string,
) (domain.UserOAuthIdentity, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, user_id, provider, provider_sub, email, created_at
		FROM todoapp.user_oauth_identities
		WHERE provider=$1 AND provider_sub=$2
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		provider,
		providerSub,
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
		if err == core_postgres_pool.ErrNoRows {
			return domain.UserOAuthIdentity{}, fmt.Errorf(
				"oauth identity not found provider=%s sub%s: %w",
				provider,
				providerSub,
				core_errors.ErrNotFound,
			)
		}

		return domain.UserOAuthIdentity{}, fmt.Errorf("scan error: %w", err)
	}

	return identity, nil
}
