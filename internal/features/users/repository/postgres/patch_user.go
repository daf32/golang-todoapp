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

func (r *UsersRepository) PatchUser(
	ctx context.Context,
	id int,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE todoapp.users
	SET
		full_name=$1,
		phone_number=$2,
		email=$3,
		version=version+1
	WHERE id=$4 AND version=$5
	RETURNING id, version, full_name, phone_number, email, password_hash, role, email_verified, email_verified_at, created_at;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.FullName,
		user.PhoneNumber,
		user.Email,
		id,
		user.Version,
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
				"user with id='%d' concurrently accessed: %w",
				id,
				core_errors.ErrConfict,
			)
		}

		if errors.Is(err, core_postgres_pool.ErrUniqueViolation) {
			return domain.User{}, fmt.Errorf(
				"user email already in use: %w",
				core_errors.ErrConfict,
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
