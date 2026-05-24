package auth_service_test

import (
	"context"
	"errors"
	"testing"

	core_auth "github.com/daf32/golang-todoapp/internal/core/auth"
	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	auth_service "github.com/daf32/golang-todoapp/internal/features/auth/service"
	users_service "github.com/daf32/golang-todoapp/internal/features/users/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser(t *testing.T) {
	const plainPassword = core_auth.PlainPassword("Password1")
	const confirmationURL = "https://test.local/confirm?token="

	repoErr := errors.New("repository unavailable")

	testCases := []struct {
		name           string
		input          domain.User
		password       core_auth.PlainPassword
		repoCreateErr  error
		wantRepoCreate bool
		wantServiceErr error
	}{
		{
			name: "successful creation",
			input: domain.NewUserUninitialized("John Doe", nil, "john@example.com",
				domain.UserRoleUser),
			password:       plainPassword,
			wantRepoCreate: true,
		},
		{
			name: "reject invalid password",
			input: domain.NewUserUninitialized("John Doe", nil, "john@example.com",
				domain.UserRoleUser),
			password:       core_auth.PlainPassword("nodigits"),
			wantServiceErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "reject invalid user domain",
			input: domain.NewUserUninitialized("Jo", nil, "john@example.com",
				domain.UserRoleUser),
			password:       plainPassword,
			wantServiceErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "wrap repository error",
			input: domain.NewUserUninitialized("John Doe", nil, "john@example.com",
				domain.UserRoleUser),
			password:       plainPassword,
			repoCreateErr:  repoErr,
			wantRepoCreate: true,
			wantServiceErr: repoErr,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			authRepo := auth_service.NewMockAuthRepository(t)
			usersRepo := users_service.NewMockUsersRepository(t)
			mailer := newNoOpMailerMock(t)

			ctx := context.Background()

			if tc.wantRepoCreate {
				authRepo.On("CreateUser",
					ctx,
					mock.MatchedBy(func(u domain.User) bool {
						if u.FullName != tc.input.FullName ||
							u.Email != tc.input.Email ||
							u.Role != tc.input.Role {
							return false
						}
						return bcrypt.CompareHashAndPassword(
							[]byte(u.PasswordHash),
							[]byte(string(tc.password)),
						) == nil
					}),
				).Return(func(_ context.Context, u domain.User) (domain.User, error) {
					if tc.repoCreateErr != nil {
						return domain.User{}, tc.repoCreateErr
					}
					u.ID = 42
					u.Version = 1
					return u, nil
				}).Once()

				if tc.repoCreateErr == nil {
					authRepo.On("CreateEmailConfirmationToken",
						mock.Anything, 42, mock.Anything,
					).Return(domain.EmailConfirmationToken{Token: "tok"}, nil).Maybe()

					mailer.On("SendEmail",
						mock.Anything, mock.Anything, mock.Anything,
					).Return(nil).Maybe()
				}
			}

			svc := auth_service.NewAuthService(
				authRepo,
				usersRepo,
				mailer,
				nil,
				"test-secret",
				0, 0, 0,
			)

			got, err := svc.CreateUser(ctx, tc.input, tc.password, confirmationURL)

			if tc.wantServiceErr != nil {
				assert.ErrorIs(t, err, tc.wantServiceErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, 42, got.ID)
			assert.Equal(t, tc.input.Email, got.Email)
			assert.Equal(t, tc.input.FullName, got.FullName)
			assert.NoError(t,
				bcrypt.CompareHashAndPassword([]byte(got.PasswordHash),
					[]byte(string(tc.password))),
				"returned hash should validate against plaintext password",
			)
		})
	}
}
