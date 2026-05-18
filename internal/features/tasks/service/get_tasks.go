package tasks_service

import (
	"context"
	"fmt"

	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
)

func (s *TasksService) GetTasks(
	ctx context.Context,
	actor domain.Actor,
	userID *int,
	limit *int,
	offset *int,
) ([]domain.Task, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf(
			"limit must be non-negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf(
			"offset must be non negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if !actor.IsAdmin() {
		if userID != nil && *userID != actor.UserID {
			return nil, fmt.Errorf(
				"user id=%d has no access to tasks of user id=%d: %w",
				actor.UserID,
				*userID,
				core_errors.ErrForbidden,
			)
		}

		userID = &actor.UserID
	}

	tasks, err := s.tasksRepository.GetTasks(
		ctx,
		userID,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("get tasks from repository: %w", err)
	}

	return tasks, nil
}
