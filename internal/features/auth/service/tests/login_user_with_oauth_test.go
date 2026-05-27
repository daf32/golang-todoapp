package auth_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	core_auth "github.com/daf32/golang-todoapp/internal/core/auth"
	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	core_oauth "github.com/daf32/golang-todoapp/internal/core/oauth"
	auth_service "github.com/daf32/golang-todoapp/internal/features/auth/service"
	users_service "github.com/daf32/golang-todoapp/internal/features/users/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLoginUserWithOAuth(t *testing.T) {
	t.Parallel()

	info := validOAuthUserInfo()
	exchangeErr := errors.New("oauth exchange boom")

	t.Run("unknown provider — returns ErrInvalidArgument before any exchange", func(t *testing.T) {
		t.Parallel()

		authRepo := auth_service.NewMockAuthRepository(t)
		usersRepo := users_service.NewMockUsersRepository(t)
		mailer := newNoOpMailerMock(t)
		google := newMockGoogleProvider(t)

		srvc := auth_service.NewAuthService(
			authRepo, usersRepo, mailer, nil,
			"test-secret", 0, 0, 0,
			[]core_oauth.Provider{google},
		)

		_, _, err := srvc.LoginUserWithOAuth(
			context.Background(), "mastodon", "code", "verifier",
		)
		assert.ErrorIs(t, err, core_errors.ErrInvalidArgument)
	})

	t.Run("exchange error propagates", func(t *testing.T) {
		t.Parallel()

		authRepo := auth_service.NewMockAuthRepository(t)
		usersRepo := users_service.NewMockUsersRepository(t)
		mailer := newNoOpMailerMock(t)
		google := newMockGoogleProvider(t)

		google.On("Exchange", mock.Anything, "bad", "verifier").
			Return(core_oauth.UserInfo{}, exchangeErr).Once()

		srvc := auth_service.NewAuthService(
			authRepo, usersRepo, mailer, nil,
			"test-secret", 0, 0, 0,
			[]core_oauth.Provider{google},
		)

		_, _, err := srvc.LoginUserWithOAuth(
			context.Background(), core_oauth.ProviderGoogle, "bad", "verifier",
		)
		assert.ErrorIs(t, err, exchangeErr)
	})

	t.Run("identity already linked — load user and issue tokens", func(t *testing.T) {
		t.Parallel()

		existing := validVerifiedUser()
		existing.ID = 42
		existing.Email = info.Email

		authRepo := auth_service.NewMockAuthRepository(t)
		usersRepo := users_service.NewMockUsersRepository(t)
		mailer := newNoOpMailerMock(t)
		google := newMockGoogleProvider(t)

		google.On("Exchange", mock.Anything, "code", "verifier").
			Return(info, nil).Once()

		authRepo.On("GetUserOAuthIdentity",
			mock.Anything, core_oauth.ProviderGoogle, info.Sub,
		).Return(domain.UserOAuthIdentity{UserID: existing.ID}, nil).Once()

		usersRepo.On("GetUser", mock.Anything, existing.ID).
			Return(existing, nil).Once()

		authRepo.On("CreateRefreshToken",
			mock.Anything, existing.ID, mock.Anything,
		).Return(core_auth.RefreshToken{Token: "refresh-abc"}, nil).Once()

		srvc := auth_service.NewAuthService(
			authRepo, usersRepo, mailer, nil,
			"test-secret", time.Minute, time.Hour, 0,
			[]core_oauth.Provider{google},
		)

		access, refresh, err := srvc.LoginUserWithOAuth(
			context.Background(), core_oauth.ProviderGoogle, "code", "verifier",
		)
		require.NoError(t, err)
		assert.NotEmpty(t, access)
		assert.Equal(t, "refresh-abc", refresh.Token)
	})

	t.Run("identity missing, user exists by email — link identity", func(t *testing.T) {
		t.Parallel()

		existing := validVerifiedUser()
		existing.ID = 43
		existing.Email = info.Email

		authRepo := auth_service.NewMockAuthRepository(t)
		usersRepo := users_service.NewMockUsersRepository(t)
		mailer := newNoOpMailerMock(t)
		google := newMockGoogleProvider(t)

		google.On("Exchange", mock.Anything, "code", "verifier").
			Return(info, nil).Once()

		authRepo.On("GetUserOAuthIdentity",
			mock.Anything, core_oauth.ProviderGoogle, info.Sub,
		).Return(domain.UserOAuthIdentity{}, core_errors.ErrNotFound).Once()

		usersRepo.On("GetUserByEmail", mock.Anything, info.Email).
			Return(existing, nil).Once()

		authRepo.On("CreateUserOAuthIdentity",
			mock.Anything, existing.ID, core_oauth.ProviderGoogle, info.Sub, info.Email,
		).Return(domain.UserOAuthIdentity{}, nil).Once()

		authRepo.On("CreateRefreshToken",
			mock.Anything, existing.ID, mock.Anything,
		).Return(core_auth.RefreshToken{Token: "refresh-link"}, nil).Once()

		srvc := auth_service.NewAuthService(
			authRepo, usersRepo, mailer, nil,
			"test-secret", time.Minute, time.Hour, 0,
			[]core_oauth.Provider{google},
		)

		_, refresh, err := srvc.LoginUserWithOAuth(
			context.Background(), core_oauth.ProviderGoogle, "code", "verifier",
		)
		require.NoError(t, err)
		assert.Equal(t, "refresh-link", refresh.Token)
	})

	t.Run("brand new user — create user and identity", func(t *testing.T) {
		t.Parallel()

		authRepo := auth_service.NewMockAuthRepository(t)
		usersRepo := users_service.NewMockUsersRepository(t)
		mailer := newNoOpMailerMock(t)
		google := newMockGoogleProvider(t)

		google.On("Exchange", mock.Anything, "code", "verifier").
			Return(info, nil).Once()

		authRepo.On("GetUserOAuthIdentity",
			mock.Anything, core_oauth.ProviderGoogle, info.Sub,
		).Return(domain.UserOAuthIdentity{}, core_errors.ErrNotFound).Once()

		usersRepo.On("GetUserByEmail", mock.Anything, info.Email).
			Return(domain.User{}, core_errors.ErrNotFound).Once()

		authRepo.On("CreateUser",
			mock.Anything,
			mock.MatchedBy(func(u domain.User) bool {
				return u.Email == info.Email &&
					u.EmailVerified &&
					u.EmailVerifiedAt != nil &&
					u.PasswordHash == ""
			}),
		).Return(func(_ context.Context, u domain.User) (domain.User, error) {
			u.ID = 99
			u.Version = 1
			return u, nil
		}).Once()

		authRepo.On("CreateUserOAuthIdentity",
			mock.Anything, 99, core_oauth.ProviderGoogle, info.Sub, info.Email,
		).Return(domain.UserOAuthIdentity{}, nil).Once()

		authRepo.On("CreateRefreshToken",
			mock.Anything, 99, mock.Anything,
		).Return(core_auth.RefreshToken{Token: "refresh-new"}, nil).Once()

		srvc := auth_service.NewAuthService(
			authRepo, usersRepo, mailer, nil,
			"test-secret", time.Minute, time.Hour, 0,
			[]core_oauth.Provider{google},
		)

		_, refresh, err := srvc.LoginUserWithOAuth(
			context.Background(), core_oauth.ProviderGoogle, "code", "verifier",
		)
		require.NoError(t, err)
		assert.Equal(t, "refresh-new", refresh.Token)
	})
}
