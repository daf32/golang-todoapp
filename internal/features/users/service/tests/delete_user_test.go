package users_service_test

import (
	"context"
	"errors"
	"testing"

	users_service "github.com/daf32/golang-todoapp/internal/features/users/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteUser(t *testing.T) {
	repositoryErr := errors.New("repository unavailable")

	testCases := []struct {
		name           string
		ctx            context.Context
		id             int
		repositoryErr  error
		wantServiceErr error
	}{
		{
			name: "successful delete user",
			ctx:  context.Background(),
			id:   1,
		},
		{
			name:           "wrap repository error",
			ctx:            context.Background(),
			id:             1,
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
				"DeleteUser",
				tc.ctx,
				tc.id,
			).Return(tc.repositoryErr).Once()

			srvc := users_service.NewUsersService(repo, nil, nil)

			err := srvc.DeleteUser(tc.ctx, tc.id)

			if tc.wantServiceErr != nil {
				assert.ErrorIs(t, err, tc.wantServiceErr)
				return
			}

			require.NoError(t, err)
		})
	}
}
