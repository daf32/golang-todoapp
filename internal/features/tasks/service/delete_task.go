package tasks_service

import (
	"context"
	"fmt"

	"github.com/daf32/golang-todoapp/internal/core/domain"
)

func (s *TasksService) DeleteTask(
	ctx context.Context,
	actor domain.Actor,
	id int,
) error {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return fmt.Errorf("get task from repository: %w", err)
	}

	if err := authorizeTaskAccess(actor, task); err != nil {
		return err
	}

	if err := s.tasksRepository.DeleteTask(ctx, id); err != nil {
		return fmt.Errorf("delete task from repository: %w", err)
	}

	return nil
}
