package statistics_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	statistics_service "github.com/daf32/golang-todoapp/internal/features/statistics/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	dateFrom = time.Date(
		2026,
		time.May,
		11,
		17,
		18,
		19,
		52,
		time.UTC,
	)

	dateToBefore = dateFrom.Add(-time.Second)
)

func TestGetStatistics(t *testing.T) {
	completedAtFast := dateFrom.Add(30 * time.Minute)
	completedAtSlow := dateFrom.Add(90 * time.Minute)

	tasksCompletedRate := float64(2) / float64(3) * 100
	tasksAverageCompletionTime := time.Hour
	repositoryErr := errors.New("repository unavailable")

	testCases := []struct {
		name            string
		ctx             context.Context
		userID          *int
		from            *time.Time
		to              *time.Time
		repositoryTasks []domain.Task
		repositoryErr   error
		wantStatistics  domain.Statistics
		wantServiceErr  error
		wantRepoCalled  bool
	}{
		{
			name: "successful get statistics",
			ctx:  context.Background(),
			repositoryTasks: []domain.Task{
				{
					Title:        "todo 1",
					Completed:    true,
					CreatedAt:    dateFrom,
					CompletedAt:  &completedAtFast,
					AuthorUserID: 1,
				},
				{
					Title:        "todo 2",
					Completed:    true,
					CreatedAt:    dateFrom,
					CompletedAt:  &completedAtSlow,
					AuthorUserID: 1,
				},
				{
					Title:        "todo 3",
					Completed:    false,
					CreatedAt:    dateFrom,
					CompletedAt:  nil,
					AuthorUserID: 1,
				},
			},
			wantStatistics: domain.Statistics{
				TasksCreated:               3,
				TasksCompleted:             2,
				TasksCompletedRate:         &tasksCompletedRate,
				TasksAverageCompletionTime: &tasksAverageCompletionTime,
			},
			wantRepoCalled: true,
		},
		{
			name:           "wrap repository error",
			ctx:            context.Background(),
			repositoryErr:  repositoryErr,
			wantServiceErr: repositoryErr,
			wantRepoCalled: true,
		},
		{
			name:           "reject when to is before from",
			ctx:            context.Background(),
			from:           &dateFrom,
			to:             &dateToBefore,
			wantServiceErr: core_errors.ErrInvalidArgument,
			wantRepoCalled: false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := statistics_service.NewMockStatisticsRepository(t)
			if tc.wantRepoCalled {
				repo.On(
					"GetTasks",
					tc.ctx,
					tc.userID,
					tc.from,
					tc.to,
				).Return(tc.repositoryTasks, tc.repositoryErr).Once()
			}

			srvc := statistics_service.NewStatisticsService(repo)

			statistics, err := srvc.GetStatistics(
				tc.ctx,
				tc.userID,
				tc.from,
				tc.to,
			)

			if tc.wantServiceErr != nil {
				assert.ErrorIs(t, err, tc.wantServiceErr)
				if !tc.wantRepoCalled {
					repo.AssertNotCalled(t, "GetTasks", tc.ctx, tc.userID, tc.from, tc.to)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantStatistics, statistics)
		})
	}
}
