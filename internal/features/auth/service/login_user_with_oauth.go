package auth_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	core_auth "github.com/daf32/golang-todoapp/internal/core/auth"
	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	core_oauth "github.com/daf32/golang-todoapp/internal/core/oauth"
)

func (s *AuthService) LoginUserWithOAuth(
	ctx context.Context,
	providerName, code, codeVerifier string,
) (string, core_auth.RefreshToken, error) {
	provider, ok := s.oauthProviders[providerName]
	if !ok {
		return "", core_auth.RefreshToken{}, fmt.Errorf(
			"unknown oauth provider %q: %w",
			providerName,
			core_errors.ErrInvalidArgument,
		)
	}

	userInfo, err := provider.Exchange(ctx, code, codeVerifier)
	if err != nil {
		return "", core_auth.RefreshToken{}, fmt.Errorf("exchange oauth code: %w", err)
	}

	user, err := s.findOrCreateOAuthUser(ctx, providerName, userInfo)
	if err != nil {
		return "", core_auth.RefreshToken{}, err
	}

	return s.issueTokens(ctx, user)
}

func (s *AuthService) findOrCreateOAuthUser(
	ctx context.Context,
	providerName string,
	info core_oauth.UserInfo,
) (domain.User, error) {
	identity, err := s.authRepository.GetUserOAuthIdentity(ctx, providerName, info.Sub)
	if err == nil {
		user, err := s.usersRepository.GetUser(ctx, identity.UserID)
		if err != nil {
			return domain.User{}, fmt.Errorf("get user for oauth identity: %w", err)
		}

		return user, nil
	}

	if !errors.Is(err, core_errors.ErrNotFound) {
		return domain.User{}, fmt.Errorf("lookup oauth identity: %w", err)
	}

	existingUser, err := s.usersRepository.GetUserByEmail(ctx, info.Email)
	if err == nil {
		if _, err := s.authRepository.CreateUserOAuthIdentity(
			ctx,
			existingUser.ID,
			providerName,
			info.Sub,
			info.Email,
		); err != nil {
			return domain.User{}, fmt.Errorf("link oauth identity to existing user: %w", err)
		}

		return existingUser, nil
	}

	if !errors.Is(err, core_errors.ErrNotFound) {
		return domain.User{}, fmt.Errorf("lookup user by email: %w", err)
	}

	return s.createOAuthUser(ctx, providerName, info)
}

func (s *AuthService) createOAuthUser(
	ctx context.Context,
	providerName string,
	info core_oauth.UserInfo,
) (domain.User, error) {
	newUser := domain.NewUserUninitialized(
		info.Name,
		nil,
		info.Email,
		domain.UserRoleUser,
	)
	newUser.EmailVerified = true
	newUser.EmailVerifiedAt = timePtr(time.Now())

	createdUser, err := s.authRepository.CreateUser(ctx, newUser)
	if err != nil {
		return domain.User{}, fmt.Errorf("create user via oauth: %w", err)
	}

	if _, err := s.authRepository.CreateUserOAuthIdentity(
		ctx,
		createdUser.ID,
		providerName,
		info.Sub,
		info.Email,
	); err != nil {
		return domain.User{}, fmt.Errorf("create oauth identity for new user: %w", err)
	}

	return createdUser, nil
}

func timePtr(t time.Time) *time.Time { return &t }

func (s *AuthService) BuildOAuthURL(
	providerName, state, codeVerifier string,
) (string, error) {
	provider, ok := s.oauthProviders[providerName]
	if !ok {
		return "", fmt.Errorf(
			"unknown oauth provider %q: %w",
			providerName,
			core_errors.ErrInvalidArgument,
		)
	}

	return provider.AuthCodeURL(state, codeVerifier), nil
}
