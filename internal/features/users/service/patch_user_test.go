package users_service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	users_service "github.com/daf32/golang-todoapp/internal/features/users/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPatchUser(t *testing.T) {
	repositoryGetErr := errors.New("repository get unavailable")
	repositoryPatchErr := errors.New("repository patch unavailable")
	newFullName := "Jane Smith"

	originalUser := validUser()
	patchedUser := validUser()
	patchedUser.FullName = newFullName
	repositoryPatchedUser := patchedUser
	repositoryPatchedUser.Version = 2

	testCases := []struct {
		name                string
		ctx                 context.Context
		id                  int
		patch               domain.UserPatch
		repositoryUser      domain.User
		repositoryGetErr    error
		repositoryPatchUser domain.User
		repositoryPatchErr  error
		wantServiceErr      error
		wantPatchCalled     bool
	}{
		{
			name: "successful patch user",
			ctx:  context.Background(),
			id:   originalUser.ID,
			patch: domain.UserPatch{
				FullName: domain.Nullable[string]{
					Value: &newFullName,
					Set:   true,
				},
			},
			repositoryUser:      originalUser,
			repositoryPatchUser: repositoryPatchedUser,
			wantPatchCalled:     true,
		},
		{
			name:             "wrap get user repository error",
			ctx:              context.Background(),
			id:               originalUser.ID,
			patch:            domain.UserPatch{},
			repositoryUser:   validUser(),
			repositoryGetErr: repositoryGetErr,
			wantServiceErr:   repositoryGetErr,
			wantPatchCalled:  false,
		},
		{
			name: "reject invalid patch before patch repository call",
			ctx:  context.Background(),
			id:   originalUser.ID,
			patch: domain.UserPatch{
				FullName: domain.Nullable[string]{
					Value: nil,
					Set:   true,
				},
			},
			repositoryUser:  originalUser,
			wantServiceErr:  core_errors.ErrInvalidArgument,
			wantPatchCalled: false,
		},
		{
			name: "reject invalid patched user before patch repository call",
			ctx:  context.Background(),
			id:   originalUser.ID,
			patch: domain.UserPatch{
				PhoneNumber: domain.Nullable[string]{
					Value: stringPtr("not-a-phone"),
					Set:   true,
				},
			},
			repositoryUser:  originalUser,
			wantServiceErr:  core_errors.ErrInvalidArgument,
			wantPatchCalled: false,
		},
		{
			name: "wrap patch repository error",
			ctx:  context.Background(),
			id:   originalUser.ID,
			patch: domain.UserPatch{
				FullName: domain.Nullable[string]{
					Value: stringPtr(newFullName),
					Set:   true,
				},
			},
			repositoryUser:      originalUser,
			repositoryPatchErr:  repositoryPatchErr,
			repositoryPatchUser: validUser(),
			wantServiceErr:      repositoryPatchErr,
			wantPatchCalled:     true,
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
			).Return(tc.repositoryUser, tc.repositoryGetErr).Once()

			if tc.wantPatchCalled {
				expectedPatchedUser := originalUser
				expectedPatchedUser.FullName = newFullName

				repo.On(
					"PatchUser",
					tc.ctx,
					tc.id,
					expectedPatchedUser,
				).Return(tc.repositoryPatchUser, tc.repositoryPatchErr).Once()
			}

			srvc := users_service.NewUsersService(repo)

			user, err := srvc.PatchUser(tc.ctx, tc.id, tc.patch)

			if tc.wantServiceErr != nil {
				assert.ErrorIs(t, err, tc.wantServiceErr)
				if !tc.wantPatchCalled {
					repo.AssertNotCalled(t, "PatchUser", tc.ctx, tc.id, mock.Anything)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.repositoryPatchUser, user)
		})
	}
}
