package auth_service

import (
	"context"
	"fmt"
	"time"

	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
)

func (s *AuthService) RefreshAccessToken(
	ctx context.Context,
	refreshTokenString string,
) (string, error) {
	token, err := s.authRepository.GetRefreshToken(ctx, refreshTokenString)
	if err != nil {
		return "", fmt.Errorf("get refresh token='%v': %w", token, err)
	}

	if token.Revoked {
		return "", fmt.Errorf(
			"token='%v' revoked: %w",
			token,
			core_errors.ErrInvalidCredentials,
		)
	}

	if time.Now().After(token.ExpiresAt) {
		return "", fmt.Errorf(
			"token='%v' expired: %w",
			token,
			core_errors.ErrInvalidCredentials,
		)
	}

	user, err := s.usersRepository.GetUser(ctx, token.UserID)
	if err != nil {
		return "", fmt.Errorf(
			"get user with id='%d': %w",
			token.UserID,
			err,
		)
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return "", fmt.Errorf(
			"generate access token: %w",
			err,
		)
	}

	return accessToken, nil
}
