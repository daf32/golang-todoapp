package tasks_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	tasks_service "github.com/daf32/golang-todoapp/internal/features/tasks/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	createdAt = time.Date(
		2026,
		time.May,
		11,
		17,
		18,
		19,
		52,
		time.UTC,
	)
)

func TestCreateTask(t *testing.T) {
	repositoryErr := errors.New("repository unavailable")

	testCases := []struct {
		name           string
		ctx            context.Context
		inputTask      domain.Task
		repositoryTask domain.Task
		repositoryErr  error
		wantServiceErr error
		wantRepoCalled bool
	}{
		{
			name: "successful create task",
			ctx:  context.Background(),
			inputTask: domain.Task{
				ID:           domain.UninitializedID,
				Version:      domain.UninitializedVersion,
				Title:        "test_task",
				Completed:    false,
				CreatedAt:    createdAt,
				AuthorUserID: 1,
			},
			repositoryTask: domain.Task{
				ID:           1,
				Version:      1,
				Title:        "test_task",
				Completed:    false,
				CreatedAt:    createdAt,
				AuthorUserID: 1,
			},
			wantRepoCalled: true,
		},
		{
			name: "wrap repository error",
			ctx:  context.Background(),
			inputTask: domain.Task{
				ID:           domain.UninitializedID,
				Version:      domain.UninitializedVersion,
				Title:        "test_task",
				Completed:    false,
				CreatedAt:    createdAt,
				AuthorUserID: 1,
			},
			repositoryErr:  repositoryErr,
			wantServiceErr: repositoryErr,
			wantRepoCalled: true,
		},
		{
			name: "reject invalid task before repository call",
			ctx:  context.Background(),
			inputTask: domain.Task{
				ID:           domain.UninitializedID,
				Version:      domain.UninitializedVersion,
				Title:        "",
				Completed:    false,
				CreatedAt:    createdAt,
				AuthorUserID: 1,
			},
			wantServiceErr: core_errors.ErrInvalidArgument,
			wantRepoCalled: false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := tasks_service.NewMockTasksRepository(t)
			if tc.wantRepoCalled {
				repo.On(
					"CreateTask",
					tc.ctx,
					tc.inputTask,
				).Return(tc.repositoryTask, tc.repositoryErr).Once()
			}

			srvc := tasks_service.NewTasksService(repo)

			task, err := srvc.CreateTask(
				tc.ctx,
				tc.inputTask,
			)

			if tc.wantServiceErr != nil {
				assert.ErrorIs(t, err, tc.wantServiceErr)
				if !tc.wantRepoCalled {
					repo.AssertNotCalled(t, "CreateTask", tc.ctx, tc.inputTask)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.repositoryTask, task)
		})
	}
}
