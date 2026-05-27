package users_service

import (
	"context"
	"fmt"
	"time"

	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
)

func (s *UsersService) CleanupUnverifiedUsers(
	ctx context.Context,
	minAge time.Duration,
) (int, error) {
	if minAge <= 0 {
		return 0, fmt.Errorf(
			"minAge must be positive, got %s: %w",
			minAge,
			core_errors.ErrInvalidArgument,
		)
	}

	cutoff := time.Now().Add(-minAge)

	deleted, err := s.usersRepository.DeleteUnverifiedUsersOlderThan(
		ctx,
		cutoff,
	)
	if err != nil {
		return 0, fmt.Errorf("delete unverified users: %w", err)
	}

	return deleted, nil
}
