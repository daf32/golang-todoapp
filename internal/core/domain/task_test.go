package domain

import (
	"testing"
	"time"

	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
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

	comletedAtBefore = createdAt.Add(-time.Second)
	comletedAtAfter  = createdAt.Add(time.Second)
)

func TestTaskValidate(t *testing.T) {
	testCases := []struct {
		name    string
		task    Task
		wantErr error
	}{
		{
			name: "success validation",
			task: Task{
				Title:        "test_task",
				Description:  nil,
				Completed:    false,
				CreatedAt:    createdAt,
				CompletedAt:  nil,
				AuthorUserID: 1,
			},
		},
		{
			name: "reject empty title",
			task: Task{
				Title:        "",
				Description:  nil,
				Completed:    false,
				CreatedAt:    createdAt,
				CompletedAt:  nil,
				AuthorUserID: 1,
			},
			wantErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "reject more than 100 title",
			task: Task{
				Title:        "ttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttt",
				Description:  nil,
				Completed:    false,
				CreatedAt:    createdAt,
				CompletedAt:  nil,
				AuthorUserID: 1,
			},
			wantErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "reject when CompletedAt is nil if Completed is true",
			task: Task{
				Title:        "test_task",
				Description:  nil,
				Completed:    true,
				CreatedAt:    createdAt,
				CompletedAt:  nil,
				AuthorUserID: 1,
			},
			wantErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "reject when CompletedAt before CreatedAt",
			task: Task{
				Title:        "test_task",
				Description:  nil,
				Completed:    true,
				CreatedAt:    createdAt,
				CompletedAt:  &comletedAtBefore,
				AuthorUserID: 1,
			},
			wantErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "reject when Completed is false, but CompletedAt is not nil",
			task: Task{
				Title:        "test_task",
				Description:  nil,
				Completed:    false,
				CreatedAt:    createdAt,
				CompletedAt:  &comletedAtAfter,
				AuthorUserID: 1,
			},
			wantErr: core_errors.ErrInvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.task.Validate()

			if tc.wantErr == nil {
				require.NoError(t, err)
				return
			}

			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestTaskApplyPatch(t *testing.T) {
	task := Task{
		Title:       "test_task",
		Description: nil,
		Completed:   false,
		CreatedAt: time.Date(
			2026,
			time.May,
			11,
			17,
			18,
			19,
			52,
			time.UTC,
		),
		CompletedAt:  nil,
		AuthorUserID: 1,
	}

	newTitle := "new_title"

	testCases := []struct {
		name              string
		patch             TaskPatch
		wantValidateErr   error
		wantApplyPatchErr error
	}{
		{
			name: "success apply patch",
			patch: TaskPatch{
				Title: Nullable[string]{
					Value: &newTitle,
					Set:   true,
				},
			},
		},
		{
			name: "reject validate patch: completed is set and value is nil",
			patch: TaskPatch{
				Completed: Nullable[bool]{
					Value: nil,
					Set:   true,
				},
			},
			wantValidateErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "reject validate patch: title is set and value is nil",
			patch: TaskPatch{
				Title: Nullable[string]{
					Value: nil,
					Set:   true,
				},
			},
			wantValidateErr: core_errors.ErrInvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.patch.Validate()

			if tc.wantValidateErr != nil {
				assert.ErrorIs(t, err, tc.wantValidateErr)
				return
			}

			require.NoError(t, err)

			err = task.ApplyPatch(tc.patch)
			if tc.wantApplyPatchErr == nil {
				require.NoError(t, err)
			}

			assert.ErrorIs(t, err, tc.wantApplyPatchErr)
		})
	}
}
