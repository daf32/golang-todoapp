package users_postgres_repository

import (
	"context"
	"fmt"
	"time"
)

func (r *UsersRepository) DeleteUnverifiedUsersOlderThan(
	ctx context.Context,
	cutoff time.Time,
) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tag, err := r.pool.Exec(
		ctx,
		`
		DELETE FROM todoapp.users
		WHERE email_verified = false
			AND created_at < $1;
		`,
		cutoff,
	)
	if err != nil {
		return 0, fmt.Errorf("delete unverified users: %w", err)
	}

	return int(tag.RowsAffected()), nil
}
