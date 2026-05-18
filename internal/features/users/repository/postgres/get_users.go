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
) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, full_name, phone_number, email, password_hash, role FROM todoapp.users
	ORDER BY id ASC
	LIMIT $1
	OFFSET $2;
	`

	rows, err := r.pool.Query(
		ctx,
		query,
		limit,
		offset,
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
