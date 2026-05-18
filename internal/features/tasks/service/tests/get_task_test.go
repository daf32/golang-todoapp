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

func TestGetTask(t *testing.T) {
	repositoryErr := errors.New("repository unavailable")
	task := validTask()

	testCases := []struct {
		name           string
		ctx            context.Context
		actor          domain.Actor
		id             int
		repositoryTask domain.Task
		repositoryErr  error
		wantServiceErr error
	}{
		{
			name:           "successful get task",
			ctx:            context.Background(),
			actor:          validUserActor(),
			id:             task.ID,
			repositoryTask: task,
		},
		{
			name:           "wrap repository error",
			ctx:            context.Background(),
			actor:          validUserActor(),
			id:             task.ID,
			repositoryTask: validTask(),
			repositoryErr:  repositoryErr,
			wantServiceErr: repositoryErr,
		},
		{
			name: "reject access to another user's task for non-admin actor",
			ctx:  context.Background(),
			actor: domain.Actor{
				UserID: 2,
				Role:   domain.UserRoleUser,
			},
			id:             task.ID,
			repositoryTask: task,
			wantServiceErr: core_errors.ErrForbidden,
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
			).Return(tc.repositoryTask, tc.repositoryErr).Once()

			srvc := tasks_service.NewTasksService(repo)

			task, err := srvc.GetTask(tc.ctx, tc.actor, tc.id)

			if tc.wantServiceErr != nil {
				assert.ErrorIs(t, err, tc.wantServiceErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.repositoryTask, task)
		})
	}
}
