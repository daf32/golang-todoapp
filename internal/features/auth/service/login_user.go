package auth_service

import (
	"context"
	"fmt"

	core_auth "github.com/daf32/golang-todoapp/internal/core/auth"
	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
)

func (s *AuthService) LoginUser(
	ctx context.Context,
	email string,
	password core_auth.PlainPassword,
) (string, core_auth.RefreshToken, error) {
	emailLen := len([]rune(email))
	if emailLen < 5 || emailLen > 255 {
		return "", core_auth.RefreshToken{}, fmt.Errorf(
			"invalid email len: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if err := password.Validate(); err != nil {
		return "", core_auth.RefreshToken{}, fmt.Errorf("validate password: %w", err)
	}

	user, err := s.usersRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return "", core_auth.RefreshToken{}, fmt.Errorf(
			"get user by email: %v: %w",
			email,
			err,
		)
	}

	if err := core_auth.VerifyPassword(user.PasswordHash, password); err != nil {
		return "", core_auth.RefreshToken{}, fmt.Errorf(
			"verify password: %w",
			core_errors.ErrInvalidCredentials,
		)
	}

	if !user.EmailVerified {
		return "", core_auth.RefreshToken{}, fmt.Errorf(
			"login rejected for unverified email %s: %w",
			user.Email,
			core_errors.ErrEmailNotVerified,
		)
	}

	return s.issueTokens(ctx, user)
}

func (s *AuthService) issueTokens(
	ctx context.Context,
	user domain.User,
) (string, core_auth.RefreshToken, error) {
	acessToken, err := s.generateAccessToken(user)
	if err != nil {
		return "", core_auth.RefreshToken{}, fmt.Errorf(
			"generate access token: %w",
			err,
		)
	}

	refreshToken, err := s.authRepository.CreateRefreshToken(ctx, user.ID, s.refreshTokenTTL)
	if err != nil {
		return "", core_auth.RefreshToken{}, fmt.Errorf(
			"create refresh token: %w",
			err,
		)
	}

	return acessToken, refreshToken, nil
}
