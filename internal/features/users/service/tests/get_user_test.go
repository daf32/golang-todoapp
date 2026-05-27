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

func TestGetUser(t *testing.T) {
	repositoryErr := errors.New("repository unavailable")
	user := validUser()

	testCases := []struct {
		name           string
		ctx            context.Context
		actor          domain.Actor
		id             int
		repositoryUser domain.User
		repositoryErr  error
		wantServiceErr error
		wantRepoCalled bool
	}{
		{
			name:           "successful get user",
			ctx:            context.Background(),
			actor:          validUserActor(),
			id:             user.ID,
			repositoryUser: user,
			wantRepoCalled: true,
		},
		{
			name:           "wrap repository error",
			ctx:            context.Background(),
			actor:          validUserActor(),
			id:             user.ID,
			repositoryUser: validUser(),
			repositoryErr:  repositoryErr,
			wantServiceErr: repositoryErr,
			wantRepoCalled: true,
		},
		{
			name: "reject access to another user for non-admin actor",
			ctx:  context.Background(),
			actor: domain.Actor{
				UserID: 2,
				Role:   domain.UserRoleUser,
			},
			id:             user.ID,
			wantServiceErr: core_errors.ErrForbidden,
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
					"GetUser",
					tc.ctx,
					tc.id,
				).Return(tc.repositoryUser, tc.repositoryErr).Once()
			}

			srvc := users_service.NewUsersService(repo, nil, nil)

			got, err := srvc.GetUser(tc.ctx, tc.actor, tc.id)

			if tc.wantServiceErr != nil {
				assert.ErrorIs(t, err, tc.wantServiceErr)
				if !tc.wantRepoCalled {
					repo.AssertNotCalled(t, "GetUser", tc.ctx, tc.id)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.repositoryUser, got)
		})
	}
}
