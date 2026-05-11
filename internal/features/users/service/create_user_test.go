package users_service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	users_service "github.com/daf32/golang-todoapp/internal/features/users/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	repositoryErr := errors.New("repository unavailable")

	testCases := []struct {
		name           string
		ctx            context.Context
		inputUser      domain.User
		repositoryUser domain.User
		repositoryErr  error
		wantServiceErr error
		wantRepoCalled bool
	}{
		{
			name: "successful create user",
			ctx:  context.Background(),
			inputUser: domain.User{
				ID:          domain.UninitializedID,
				Version:     domain.UninitializedVersion,
				FullName:    "John Doe",
				PhoneNumber: nil,
			},
			repositoryUser: domain.User{
				ID:          1,
				Version:     1,
				FullName:    "John Doe",
				PhoneNumber: nil,
			},
			wantRepoCalled: true,
		},
		{
			name: "wrap repository error",
			ctx:  context.Background(),
			inputUser: domain.User{
				ID:          domain.UninitializedID,
				Version:     domain.UninitializedVersion,
				FullName:    "John Doe",
				PhoneNumber: nil,
			},
			repositoryErr:  repositoryErr,
			wantServiceErr: repositoryErr,
			wantRepoCalled: true,
		},
		{
			name: "reject invalid user before repository call",
			ctx:  context.Background(),
			inputUser: domain.User{
				ID:          domain.UninitializedID,
				Version:     domain.UninitializedVersion,
				FullName:    "JD",
				PhoneNumber: nil,
			},
			wantServiceErr: core_errors.ErrInvalidArgument,
			wantRepoCalled: false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := users_service.NewMockUsersRepository(t)
			if tc.wantRepoCalled {
				repo.On(
					"CreateUser",
					tc.ctx,
					tc.inputUser,
				).Return(tc.repositoryUser, tc.repositoryErr).Once()
			}

			srvc := users_service.NewUsersService(repo)

			user, err := srvc.CreateUser(
				tc.ctx,
				tc.inputUser,
			)

			if tc.wantServiceErr != nil {
				assert.ErrorIs(t, err, tc.wantServiceErr)
				if !tc.wantRepoCalled {
					repo.AssertNotCalled(t, "CreateUser", tc.ctx, tc.inputUser)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.repositoryUser, user)
		})
	}
}
