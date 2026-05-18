package users_service

import (
	"context"
	"fmt"

	"github.com/daf32/golang-todoapp/internal/core/domain"
)

func (s *UsersService) GetUser(
	ctx context.Context,
	actor domain.Actor,
	id int,
) (domain.User, error) {
	if err := authorizeUserAccess(actor, id); err != nil {
		return domain.User{}, err
	}

	user, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf(
			"get user from repository: %w",
			err,
		)
	}

	return user, nil
}
