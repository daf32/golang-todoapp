package users_service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/daf32/golang-todoapp/internal/core/domain"
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
		id             int
		repositoryUser domain.User
		repositoryErr  error
		wantServiceErr error
	}{
		{
			name:           "successful get user",
			ctx:            context.Background(),
			id:             user.ID,
			repositoryUser: user,
		},
		{
			name:           "wrap repository error",
			ctx:            context.Background(),
			id:             user.ID,
			repositoryUser: validUser(),
			repositoryErr:  repositoryErr,
			wantServiceErr: repositoryErr,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := users_service.NewMockUsersRepository(t)
			repo.On(
				"GetUser",
				tc.ctx,
				tc.id,
			).Return(tc.repositoryUser, tc.repositoryErr).Once()

			srvc := users_service.NewUsersService(repo)

			got, err := srvc.GetUser(tc.ctx, tc.id)

			if tc.wantServiceErr != nil {
				assert.ErrorIs(t, err, tc.wantServiceErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.repositoryUser, got)
		})
	}
}
