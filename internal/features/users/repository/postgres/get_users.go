package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_models "github.com/daf32/golang-todoapp/internal/core/repository/models"
)

func (r *UsersRepository) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
	emailVerified *bool,
) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, full_name, phone_number, email, password_hash, role, email_verified, email_verified_at, created_at FROM todoapp.users
	`

	args := []any{}

	if emailVerified != nil {
		args = append(args, *emailVerified)
		query += fmt.Sprintf(" WHERE email_verified = $%d", len(args))
	}

	query += " ORDER BY id"

	if limit != nil {
		args = append(args, *limit)
		query += fmt.Sprintf(" LIMIT $%d", len(args))
	}

	if offset != nil {
		args = append(args, *offset)
		query += fmt.Sprintf(" OFFSET $%d", len(args))
	}

	rows, err := r.pool.Query(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}
	defer rows.Close()

	var userModels []core_models.UserModel
	for rows.Next() {
		var userModel core_models.UserModel

		err := rows.Scan(
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
			return nil, fmt.Errorf("scan users: %w", err)
		}

		userModels = append(userModels, userModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	return core_models.UserDomainsFromModels(userModels), nil
}
