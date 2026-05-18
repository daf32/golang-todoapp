package users_service

import (
	"context"
	"fmt"

	"github.com/daf32/golang-todoapp/internal/core/domain"
)

func (s *UsersService) PatchUser(
	ctx context.Context,
	actor domain.Actor,
	id int,
	patch domain.UserPatch,
) (domain.User, error) {
	if err := authorizeUserAccess(actor, id); err != nil {
		return domain.User{}, err
	}

	user, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user: %w", err)
	}

	if err := user.ApplyPatch(patch); err != nil {
		return domain.User{}, fmt.Errorf("apply patch: %w", err)
	}

	patchedUser, err := s.usersRepository.PatchUser(ctx, id, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("patch user: %w", err)
	}

	return patchedUser, nil
}
