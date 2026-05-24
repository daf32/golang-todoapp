package auth_service_test

import (
	"context"
	"errors"
	"testing"

	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	auth_service "github.com/daf32/golang-todoapp/internal/features/auth/service"
	users_service "github.com/daf32/golang-todoapp/internal/features/users/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfirmEmail(t *testing.T) {
	repositoryErr := errors.New("repository unavalible")
	user := ValidUnverifiedUser()

	testCases := []struct {
		name           string
		token          string
		repositoryErr  error
		wantServiceErr error
		wantRepoCalled bool
	}{
		{
			name:           "successful confirmation",
			token:          "plain-token-value",
			wantRepoCalled: true,
		},
		{
			name:           "empty token rejected at service layer",
			token:          "",
			wantServiceErr: core_errors.ErrInvalidArgument,
			wantRepoCalled: false,
		},
		{
			name:           "token not found",
			token:          "missing-token",
			repositoryErr:  core_errors.ErrNotFound,
			wantServiceErr: core_errors.ErrNotFound,
			wantRepoCalled: true,
		},
		{
			name:           "wrap repositroy error",
			token:          "any-token",
			repositoryErr:  repositoryErr,
			wantServiceErr: repositoryErr,
			wantRepoCalled: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			authRepo := auth_service.NewMockAuthRepository(t)
			usersRepo := users_service.NewMockUsersRepository(t)
			mailer := newNoOpMailerMock(t)

			if tc.wantRepoCalled {
				authRepo.On(
					"GetAndConsumeEmailConfirmationToken",
					context.Background(),
					tc.token,
				).Return(validConfirmationToken(user.ID), tc.repositoryErr).Once()
			}

			srvc := auth_service.NewAuthService(
				authRepo,
				usersRepo,
				mailer,
				nil,
				"test-secret",
				0,
				0,
				0,
			)

			err := srvc.ConfirmEmail(context.Background(), tc.token)

			if tc.wantServiceErr != nil {
				assert.ErrorIs(t, err, tc.wantServiceErr)
				return
			}

			require.NoError(t, err)
		})
	}
}
