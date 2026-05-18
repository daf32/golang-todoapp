package auth_service

import (
	"context"
	"fmt"

	core_auth "github.com/daf32/golang-todoapp/internal/core/auth/password"
	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
)

func (s *AuthService) LoginUser(
	ctx context.Context,
	email string,
	password string,
) (string, domain.RefreshToken, error) {
	emailLen := len([]rune(email))
	if emailLen < 5 || emailLen > 255 {
		return "", domain.RefreshToken{}, fmt.Errorf(
			"invalid email len: %d: %w",
			emailLen,
			core_errors.ErrInvalidArgument,
		)
	}

	passwordLen := len([]rune(password))
	if passwordLen < 5 || passwordLen > 255 {
		return "", domain.RefreshToken{}, fmt.Errorf(
			"invalid password len: %d: %w",
			passwordLen,
			core_errors.ErrInvalidArgument,
		)
	}

	user, err := s.refreshTokenRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return "", domain.RefreshToken{}, fmt.Errorf(
			"get user by email: %v: %w",
			email,
			err,
		)
	}

	if err := core_auth.VerifyPassword(user.PasswordHash, password); err != nil {
		return "", domain.RefreshToken{}, fmt.Errorf(
			"verify password: %v: %w",
			password,
			core_errors.ErrInvalidCredentials,
		)
	}

	acessToken, err := s.generateAccessToken(user)
	if err != nil {
		return "", domain.RefreshToken{}, fmt.Errorf(
			"generate access token: %w",
			err,
		)
	}

	refreshToken, err := s.refreshTokenRepository.CreateRefreshToken(ctx, user.ID, s.refreshTokenTTL)
	if err != nil {
		return "", domain.RefreshToken{}, fmt.Errorf(
			"create refresh token: %w",
			err,
		)
	}

	return acessToken, refreshToken, nil
}
