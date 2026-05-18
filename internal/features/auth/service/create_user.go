package auth_service

import (
	"context"
	"fmt"

	core_auth "github.com/daf32/golang-todoapp/internal/core/auth/password"
	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
)

func (s *AuthService) CreateUser(
	ctx context.Context,
	user domain.User,
	plainPassword string,
) (domain.User, error) {
	if plainPassword == "" {
		return domain.User{}, fmt.Errorf("password required: %w", core_errors.ErrInvalidArgument)
	}

	hashedPassword, err := core_auth.HashPassword(plainPassword)
	if err != nil {
		return domain.User{}, fmt.Errorf("generate hash password: %w", err)
	}

	user.PasswordHash = hashedPassword

	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}

	user, err = s.refreshTokenRepository.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}
