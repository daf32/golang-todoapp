package auth_service

import (
	"context"
	"fmt"

	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
)

func (s *AuthService) LogoutUser(
	ctx context.Context,
	userID int,
	refreshTokenString string,
) error {
	refreshToken, err := s.refreshTokenRepository.GetRefreshToken(ctx, refreshTokenString)
	if err != nil {
		return fmt.Errorf("get refresh token: %w", err)
	}

	if refreshToken.UserID != userID {
		return fmt.Errorf(
			"user id=%d has no access to refresh token: %w",
			userID,
			core_errors.ErrForbidden,
		)
	}

	if err := s.refreshTokenRepository.RevokeRefreshToken(ctx, refreshTokenString); err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}

	return nil
}
