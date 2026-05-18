package tasks_service

import (
	"context"
	"fmt"

	"github.com/daf32/golang-todoapp/internal/core/domain"
)

func (s *TasksService) GetTask(
	ctx context.Context,
	actor domain.Actor,
	id int,
) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task from repository: %w", err)
	}

	if err := authorizeTaskAccess(actor, task); err != nil {
		return domain.Task{}, err
	}

	return task, nil
}
