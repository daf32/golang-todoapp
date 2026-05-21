package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
)

func (r *UsersRepository) ChangeUserPassword(
	ctx context.Context,
	user domain.User,
	newPasswordHadsh string,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE todoapp.users
	SET
		password_hash=$1
	WHERE id=$2 AND version=$3
	`

	cmdTag, err := r.pool.Exec(
		ctx,
		query,
		newPasswordHadsh,
		user.ID,
		user.Version,
	)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id='%d': %w", user.ID, core_errors.ErrNotFound)
	}

	return nil
}
