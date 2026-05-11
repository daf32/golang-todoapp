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

func TestGetUsers(t *testing.T) {
	repositoryErr := errors.New("repository unavailable")
	validLimit := intPtr(10)
	validOffset := intPtr(0)
	negativeLimit := intPtr(-1)
	negativeOffset := intPtr(-1)

	secondUser := validUser()
	secondUser.ID = 2
	secondUser.Version = 2
	secondUser.FullName = "Jane Smith"

	repositoryUsers := []domain.User{
		validUser(),
		secondUser,
	}

	testCases := []struct {
		name            string
		ctx             context.Context
		limit           *int
		offset          *int
		repositoryUsers []domain.User
		repositoryErr   error
		wantServiceErr  error
		wantRepoCalled  bool
	}{
		{
			name:            "successful get users",
			ctx:             context.Background(),
			limit:           validLimit,
			offset:          validOffset,
			repositoryUsers: repositoryUsers,
			wantRepoCalled:  true,
		},
		{
			name:           "wrap repository error",
			ctx:            context.Background(),
			limit:          validLimit,
			offset:         validOffset,
			repositoryErr:  repositoryErr,
			wantServiceErr: repositoryErr,
			wantRepoCalled: true,
		},
		{
			name:           "reject negative limit",
			ctx:            context.Background(),
			limit:          negativeLimit,
			offset:         validOffset,
			wantServiceErr: core_errors.ErrInvalidArgument,
			wantRepoCalled: false,
		},
		{
			name:           "reject negative offset",
			ctx:            context.Background(),
			limit:          validLimit,
			offset:         negativeOffset,
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
					"GetUsers",
					tc.ctx,
					tc.limit,
					tc.offset,
				).Return(tc.repositoryUsers, tc.repositoryErr).Once()
			}

			srvc := users_service.NewUsersService(repo)

			users, err := srvc.GetUsers(
				tc.ctx,
				tc.limit,
				tc.offset,
			)

			if tc.wantServiceErr != nil {
				assert.ErrorIs(t, err, tc.wantServiceErr)
				if !tc.wantRepoCalled {
					repo.AssertNotCalled(
						t,
						"GetUsers",
						tc.ctx,
						tc.limit,
						tc.offset,
					)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.repositoryUsers, users)
		})
	}
}
