package tasks_service

import (
	"fmt"

	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
)

func authorizeTaskAccess(actor domain.Actor, task domain.Task) error {
	if actor.IsAdmin() || actor.UserID == task.AuthorUserID {
		return nil
	}

	return fmt.Errorf(
		"user id=%d has no access to task id=%d: %w",
		actor.UserID,
		task.ID,
		core_errors.ErrForbidden,
	)
}
