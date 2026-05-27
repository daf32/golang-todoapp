package users_service_test

import (
	"context"
	"errors"
	"testing"

	core_auth "github.com/daf32/golang-todoapp/internal/core/auth"
	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	users_service "github.com/daf32/golang-todoapp/internal/features/users/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestChangeUserPassword(t *testing.T) {
	repositoryErr := errors.New("repository unavailable")

	password := core_auth.PlainPassword("Test_password123")
	passwordHash, _ := password.Hash()

	userWithHash := validUser()
	userWithHash.PasswordHash = passwordHash

	testCases := []struct {
		name            string
		ctx             context.Context
		actor           domain.Actor
		user            domain.User
		password        core_auth.PlainPassword
		newPassword     core_auth.PlainPassword
		confirmPassword core_auth.PlainPassword
		expectRepoCall  bool
		repositoryErr   error
		wantServiceErr  error
	}{
		{
			name:            "successful change user password",
			ctx:             context.Background(),
			actor:           validUserActor(),
			user:            userWithHash,
			password:        password,
			newPassword:     core_auth.PlainPassword("Test_password456"),
			confirmPassword: core_auth.PlainPassword("Test_password456"),
			expectRepoCall:  true,
		},
		{
			name:            "wrap repository error",
			ctx:             context.Background(),
			actor:           validUserActor(),
			user:            userWithHash,
			password:        password,
			newPassword:     core_auth.PlainPassword("Test_password456"),
			confirmPassword: core_auth.PlainPassword("Test_password456"),
			expectRepoCall:  true,
			repositoryErr:   repositoryErr,
			wantServiceErr:  repositoryErr,
		},
		{
			name:            "passwords not match",
			ctx:             context.Background(),
			actor:           validAdminActor(),
			user:            userWithHash,
			password:        password,
			newPassword:     core_auth.PlainPassword("Test_password456"),
			confirmPassword: core_auth.PlainPassword("Test_password789"),
			expectRepoCall:  false,
			wantServiceErr:  core_errors.ErrInvalidArgument,
		},
		{
			name:            "invalid password",
			ctx:             context.Background(),
			actor:           validAdminActor(),
			user:            userWithHash,
			password:        core_auth.PlainPassword("SomePass1"),
			newPassword:     core_auth.PlainPassword("Test_password789"),
			confirmPassword: core_auth.PlainPassword("Test_password789"),
			expectRepoCall:  false,
			wantServiceErr:  bcrypt.ErrMismatchedHashAndPassword,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := users_service.NewMockUsersRepository(t)
			repo.On("GetUser", tc.ctx, tc.user.ID).Return(tc.user, nil).Once()
			if tc.expectRepoCall {
				repo.On(
					"ChangeUserPassword",
					tc.ctx,
					tc.user,
					mock.AnythingOfType("string"),
				).Return(tc.repositoryErr).Once()
			}

			authRepo := users_service.NewMockAuthRepository(t)
			if tc.expectRepoCall && tc.repositoryErr == nil {
				authRepo.On(
					"RevokeAllRefreshTokensForUser",
					tc.ctx,
					tc.user.ID,
				).Return(nil).Once()
			}

			svrc := users_service.NewUsersService(repo, nil, authRepo)

			err := svrc.ChangeUserPassword(
				tc.ctx,
				tc.actor,
				tc.user.ID,
				tc.password,
				tc.newPassword,
				tc.confirmPassword,
			)
			if tc.wantServiceErr != nil {
				assert.ErrorIs(t, err, tc.wantServiceErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
