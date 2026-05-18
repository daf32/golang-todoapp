package tasks_service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	tasks_service "github.com/daf32/golang-todoapp/internal/features/tasks/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteTask(t *testing.T) {
	repositoryErr := errors.New("repository unavailable")

	testCases := []struct {
		name             string
		ctx              context.Context
		actor            domain.Actor
		id               int
		repositoryTask   domain.Task
		repositoryGetErr error
		repositoryErr    error
		wantServiceErr   error
		wantDeleteCalled bool
	}{
		{
			name:             "successful delete task",
			ctx:              context.Background(),
			actor:            validUserActor(),
			id:               1,
			repositoryTask:   validTask(),
			wantDeleteCalled: true,
		},
		{
			name:             "wrap repository error",
			ctx:              context.Background(),
			actor:            validUserActor(),
			id:               1,
			repositoryTask:   validTask(),
			repositoryErr:    repositoryErr,
			wantServiceErr:   repositoryErr,
			wantDeleteCalled: true,
		},
		{
			name: "reject delete access to another user's task for non-admin actor",
			ctx:  context.Background(),
			actor: domain.Actor{
				UserID: 2,
				Role:   domain.UserRoleUser,
			},
			id:               1,
			repositoryTask:   validTask(),
			wantServiceErr:   core_errors.ErrForbidden,
			wantDeleteCalled: false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := tasks_service.NewMockTasksRepository(t)
			repo.On(
				"GetTask",
				tc.ctx,
				tc.id,
			).Return(tc.repositoryTask, tc.repositoryGetErr).Once()

			if tc.wantDeleteCalled {
				repo.On(
					"DeleteTask",
					tc.ctx,
					tc.id,
				).Return(tc.repositoryErr).Once()
			}

			srvc := tasks_service.NewTasksService(repo)

			err := srvc.DeleteTask(tc.ctx, tc.actor, tc.id)

			if tc.wantServiceErr != nil {
				assert.ErrorIs(t, err, tc.wantServiceErr)
				if !tc.wantDeleteCalled {
					repo.AssertNotCalled(t, "DeleteTask", tc.ctx, tc.id)
				}
				return
			}

			require.NoError(t, err)
		})
	}
}
