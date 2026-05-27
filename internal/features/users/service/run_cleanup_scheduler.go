package users_service

import (
	"context"
	"time"

	"go.uber.org/zap"
)

func (s *UsersService) RunUnverifiedCleanupLoop(
	ctx context.Context,
	interval, minAge time.Duration,
) {
	s.log.Info(
		"starting unverified user cleanup loop",
		zap.Duration("interval", interval),
		zap.Duration("min_age", minAge),
	)

	s.runUnverifiedCleanup(ctx, minAge)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
			case <- ctx.Done():
				s.log.Info("unverified user cleanup loop stopped")
				return
			case <-ticker.C:
				s.runUnverifiedCleanup(ctx, minAge)
		}
	}
}

func (s *UsersService) runUnverifiedCleanup(ctx context.Context, minAge time.Duration) {
	deleted, err := s.CleanupUnverifiedUsers(ctx, minAge)
	if err != nil {
		s.log.Error("unverified user cleanup failed", zap.Error(err))
		return
	}
	s.log.Info(
		"unverified user cleanup ran",
		zap.Int("deleted", deleted),
		zap.Duration("min_age", minAge),
	)
}
