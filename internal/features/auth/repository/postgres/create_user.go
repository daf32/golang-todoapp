package auth_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	core_models "github.com/daf32/golang-todoapp/internal/core/repository/models"
	core_postgres_pool "github.com/daf32/golang-todoapp/internal/core/repository/postgres/pool"
)

func (r *AuthRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO todoapp.users (full_name, phone_number, email, password_hash, role, email_verified, email_verified_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, version, full_name, phone_number, email, password_hash, role, email_verified, email_verified_at, created_at;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.FullName,
		user.PhoneNumber,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.EmailVerified,
		user.EmailVerifiedAt,
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
