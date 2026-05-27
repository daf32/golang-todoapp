package auth_service

import (
	"context"
	"fmt"

	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
)

func (s *AuthService) ConfirmEmail(
	ctx context.Context,
	token string,
) error {
	if token == "" {
		return fmt.Errorf("token is required: %w", core_errors.ErrInvalidArgument)
	}

	_, err := s.authRepository.GetAndConsumeEmailConfirmationToken(ctx, token)
	if err != nil {
		return fmt.Errorf("confirm email: %w", err)
	}

	return nil
}
