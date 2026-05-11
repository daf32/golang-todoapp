package tasks_service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	tasks_service "github.com/daf32/golang-todoapp/internal/features/tasks/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPatchTask(t *testing.T) {
	repositoryGetErr := errors.New("repository get unavailable")
	repositoryPatchErr := errors.New("repository patch unavailable")
	newTitle := "new_title"

	originalTask := validTask()
	patchedTask := validTask()
	patchedTask.Title = newTitle
	repositoryPatchedTask := patchedTask
	repositoryPatchedTask.Version = 2

	testCases := []struct {
		name                string
		ctx                 context.Context
		id                  int
		patch               domain.TaskPatch
		repositoryTask      domain.Task
		repositoryGetErr    error
		repositoryPatchTask domain.Task
		repositoryPatchErr  error
		wantServiceErr      error
		wantPatchCalled     bool
	}{
		{
			name: "successful patch task",
			ctx:  context.Background(),
			id:   originalTask.ID,
			patch: domain.TaskPatch{
				Title: domain.Nullable[string]{
					Value: &newTitle,
					Set:   true,
				},
			},
			repositoryTask:      originalTask,
			repositoryPatchTask: repositoryPatchedTask,
			wantPatchCalled:     true,
		},
		{
			name:             "wrap get task repository error",
			ctx:              context.Background(),
			id:               originalTask.ID,
			patch:            domain.TaskPatch{},
			repositoryTask:   validTask(),
			repositoryGetErr: repositoryGetErr,
			wantServiceErr:   repositoryGetErr,
			wantPatchCalled:  false,
		},
		{
			name: "reject invalid patch before patch repository call",
			ctx:  context.Background(),
			id:   originalTask.ID,
			patch: domain.TaskPatch{
				Title: domain.Nullable[string]{
					Value: nil,
					Set:   true,
				},
			},
			repositoryTask:  originalTask,
			wantServiceErr:  core_errors.ErrInvalidArgument,
			wantPatchCalled: false,
		},
		{
			name: "wrap patch repository error",
			ctx:  context.Background(),
			id:   originalTask.ID,
			patch: domain.TaskPatch{
				Title: domain.Nullable[string]{
					Value: stringPtr(newTitle),
					Set:   true,
				},
			},
			repositoryTask:      originalTask,
			repositoryPatchErr:  repositoryPatchErr,
			repositoryPatchTask: validTask(),
			wantServiceErr:      repositoryPatchErr,
			wantPatchCalled:     true,
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

			if tc.wantPatchCalled {
				expectedPatchedTask := originalTask
				expectedPatchedTask.Title = newTitle

				repo.On(
					"PatchTask",
					tc.ctx,
					tc.id,
					expectedPatchedTask,
				).Return(tc.repositoryPatchTask, tc.repositoryPatchErr).Once()
			}

			srvc := tasks_service.NewTasksService(repo)

			task, err := srvc.PatchTask(tc.ctx, tc.id, tc.patch)

			if tc.wantServiceErr != nil {
				assert.ErrorIs(t, err, tc.wantServiceErr)
				if !tc.wantPatchCalled {
					repo.AssertNotCalled(t, "PatchTask", tc.ctx, tc.id, mock.Anything)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.repositoryPatchTask, task)
		})
	}
}
