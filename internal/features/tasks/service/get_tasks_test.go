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

func TestGetTasks(t *testing.T) {
	repositoryErr := errors.New("repository unavailable")
	userID := intPtr(1)
	validLimit := intPtr(10)
	validOffset := intPtr(0)
	negativeLimit := intPtr(-1)
	negativeOffset := intPtr(-1)

	secondTask := validTask()
	secondTask.ID = 2
	secondTask.Version = 3
	secondTask.Title = "other_task"
	secondTask.AuthorUserID = 2

	repositoryTasks := []domain.Task{
		validTask(),
		secondTask,
	}

	testCases := []struct {
		name            string
		ctx             context.Context
		userID          *int
		limit           *int
		offset          *int
		repositoryTasks []domain.Task
		repositoryErr   error
		wantServiceErr  error
		wantRepoCalled  bool
	}{
		{
			name:            "successful get tasks",
			ctx:             context.Background(),
			userID:          userID,
			limit:           validLimit,
			offset:          validOffset,
			repositoryTasks: repositoryTasks,
			wantRepoCalled:  true,
		},
		{
			name:           "wrap repository error",
			ctx:            context.Background(),
			userID:         userID,
			limit:          validLimit,
			offset:         validOffset,
			repositoryErr:  repositoryErr,
			wantServiceErr: repositoryErr,
			wantRepoCalled: true,
		},
		{
			name:           "reject negative limit",
			ctx:            context.Background(),
			userID:         userID,
			limit:          negativeLimit,
			offset:         validOffset,
			wantServiceErr: core_errors.ErrInvalidArgument,
			wantRepoCalled: false,
		},
		{
			name:           "reject negative offset",
			ctx:            context.Background(),
			userID:         userID,
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

			repo := tasks_service.NewMockTasksRepository(t)
			if tc.wantRepoCalled {
				repo.On(
					"GetTasks",
					tc.ctx,
					tc.userID,
					tc.limit,
					tc.offset,
				).Return(tc.repositoryTasks, tc.repositoryErr).Once()
			}

			srvc := tasks_service.NewTasksService(repo)

			tasks, err := srvc.GetTasks(
				tc.ctx,
				tc.userID,
				tc.limit,
				tc.offset,
			)

			if tc.wantServiceErr != nil {
				assert.ErrorIs(t, err, tc.wantServiceErr)
				if !tc.wantRepoCalled {
					repo.AssertNotCalled(
						t,
						"GetTasks",
						tc.ctx,
						tc.userID,
						tc.limit,
						tc.offset,
					)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.repositoryTasks, tasks)
		})
	}
}
