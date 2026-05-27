package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	core_models "github.com/daf32/golang-todoapp/internal/core/repository/models"
	core_postgres_pool "github.com/daf32/golang-todoapp/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) GetUser(
	ctx context.Context,
	id int,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, full_name, phone_number, email, password_hash, role, email_verified, email_verified_at, created_at FROM todoapp.users
	WHERE id=$1;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		id,
	)

	var userModel core_models.UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
		&userModel.Email,
		&userModel.PasswordHash,
		&userModel.Role,
		&userModel.EmailVerified,
		&userModel.EmailVerifiedAt,
		&userModel.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id='%d': %w",
				id,
				core_errors.ErrNotFound,
			)
		}

		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	return domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
		userModel.Email,
		userModel.PasswordHash,
		userModel.Role,
		userModel.EmailVerified,
		userModel.EmailVerifiedAt,
		userModel.CreatedAt,
	), nil
}
