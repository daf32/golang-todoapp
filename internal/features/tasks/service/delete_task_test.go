package tasks_service_test

import (
	"context"
	"errors"
	"testing"

	tasks_service "github.com/daf32/golang-todoapp/internal/features/tasks/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteTask(t *testing.T) {
	repositoryErr := errors.New("repository unavailable")

	testCases := []struct {
		name           string
		ctx            context.Context
		id             int
		repositoryErr  error
		wantServiceErr error
	}{
		{
			name: "successful delete task",
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

			repo := tasks_service.NewMockTasksRepository(t)
			repo.On(
				"DeleteTask",
				tc.ctx,
				tc.id,
			).Return(tc.repositoryErr).Once()

			srvc := tasks_service.NewTasksService(repo)

			err := srvc.DeleteTask(tc.ctx, tc.id)

			if tc.wantServiceErr != nil {
				assert.ErrorIs(t, err, tc.wantServiceErr)
				return
			}

			require.NoError(t, err)
		})
	}
}
