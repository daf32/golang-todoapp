package auth_service_test

import (
	"context"
	"testing"
	"time"

	core_auth "github.com/daf32/golang-todoapp/internal/core/auth"
	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	auth_service "github.com/daf32/golang-todoapp/internal/features/auth/service"
	users_service "github.com/daf32/golang-todoapp/internal/features/users/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginUser(t *testing.T) {
	const validPassword = core_auth.PlainPassword("Password1")

	hashedPassword, err := validPassword.Hash()
	require.NoError(t, err)

	verifiedUser := validVerifiedUser()
	verifiedUser.PasswordHash = hashedPassword

	unverifiedUser := ValidUnverifiedUser()
	unverifiedUser.PasswordHash = hashedPassword

	testCases := []struct {
		name            string
		email           string
		password        core_auth.PlainPassword
		repoUser        domain.User
		repoUserErr     error
		wantUserLookup  bool
		wantRefreshCall bool
		wantServiceErr  error
	}{
		{
			name:            "successful login for verified user",
			email:           verifiedUser.Email,
			password:        validPassword,
			repoUser:        verifiedUser,
			wantUserLookup:  true,
			wantRefreshCall: true,
		},
		{
			name:           "reject login for unverified user",
			email:          unverifiedUser.Email,
			password:       validPassword,
			repoUser:       unverifiedUser,
			wantUserLookup: true,
			wantServiceErr: core_errors.ErrEmailNotVerified,
		},
		{
			name:           "reject wrong password",
			email:          verifiedUser.Email,
			password:       core_auth.PlainPassword("WrongPass1"),
			repoUser:       verifiedUser,
			wantUserLookup: true,
			wantServiceErr: core_errors.ErrInvalidCredentials,
		},
		{
			name:           "wrap repository error when user not found",
			email:          "missing@example.com",
			password:       validPassword,
			repoUserErr:    core_errors.ErrNotFound,
			wantUserLookup: true,
			wantServiceErr: core_errors.ErrNotFound,
		},
		{
			name:           "reject invalid password format before lookup",
			email:          verifiedUser.Email,
			password:       core_auth.PlainPassword("nodigits"),
			wantServiceErr: core_errors.ErrInvalidArgument,
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

			if tc.wantUserLookup {
				usersRepo.On(
					"GetUserByEmail",
					ctx,
					tc.email,
				).Return(tc.repoUser, tc.repoUserErr).Once()
			}

			if tc.wantRefreshCall {
				authRepo.On(
					"CreateRefreshToken",
					ctx,
					tc.repoUser.ID,
					time.Duration(0),
				).Return(core_auth.RefreshToken{
					ID:    uuid.New(),
					Token: "refresh-token-value",
				}, nil).Once()
			}

			svrc := auth_service.NewAuthService(
				authRepo,
				usersRepo,
				mailer,
				nil,
				"test-secret",
				0,
				0,
				0,
				nil,
			)

			accessToken, refreshToken, err := svrc.LoginUser(ctx, tc.email, tc.password)
			if tc.wantServiceErr != nil {
				assert.ErrorIs(t, err, tc.wantServiceErr)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, accessToken)
			assert.NotEmpty(t, refreshToken)
		})
	}
}
