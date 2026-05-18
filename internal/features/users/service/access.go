package users_service

import (
	"fmt"

	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
)

func authorizeUserAccess(actor domain.Actor, targetUserID int) error {
	if actor.IsAdmin() || actor.UserID == targetUserID {
		return nil
	}

	return fmt.Errorf(
		"user id=%d has no access to user id=%d: %w",
		actor.UserID,
		targetUserID,
		core_errors.ErrForbidden,
	)
}
