package users_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	users_service "github.com/daf32/golang-todoapp/internal/features/users/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCleanupUnverifiedUsers(t *testing.T) {
	repoErr := errors.New("db down")

	testCases := []struct {
		name           string
		minAge         time.Duration
		repoCount      int
		repoErr        error
		wantCount      int
		wantServiceErr error
		wantRepoCalled bool
	}{
		{
			name:           "successful cleanup",
			minAge:         7 * 24 * time.Hour,
			repoCount:      3,
			wantCount:      3,
			wantRepoCalled: true,
		},
		{
			name:           "minAge zero is rejected",
			minAge:         0,
			wantServiceErr: core_errors.ErrInvalidArgument,
			wantRepoCalled: false,
		},
		{
			name:           "negative minAge is rejected",
			minAge:         -time.Hour,
			wantServiceErr: core_errors.ErrInvalidArgument,
			wantRepoCalled: false,
		},
		{
			name:           "repo error propagates",
			minAge:         24 * time.Hour,
			repoErr:        repoErr,
			wantServiceErr: repoErr,
			wantRepoCalled: true,
		},
		{
			name:           "no users to clean up returns zero",
			minAge:         24 * time.Hour,
			repoCount:      0,
			wantCount:      0,
			wantRepoCalled: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			usersRepo := users_service.NewMockUsersRepository(t)

			if tc.wantRepoCalled {
				usersRepo.On(
					"DeleteUnverifiedUsersOlderThan",
					mock.Anything,
					mock.MatchedBy(func(t time.Time) bool {
						expected := time.Now().Add(-tc.minAge)
						return t.After(expected.Add(-time.Second)) && t.Before(expected.Add(time.Second))
					}),
				).Return(tc.repoCount, tc.repoErr).Once()
			}

			srvc := users_service.NewUsersService(usersRepo, nil, nil)

			count, err := srvc.CleanupUnverifiedUsers(context.Background(), tc.minAge)
			if tc.wantServiceErr != nil {
				assert.ErrorIs(t, err, tc.wantServiceErr)
				return
			}

			require.NoError(t, err)

			require.Equal(t, count, tc.wantCount)
		})
	}
}
