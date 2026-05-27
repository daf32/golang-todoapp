package users_service

import (
	"context"
	"time"

	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_logger "github.com/daf32/golang-todoapp/internal/core/logger"
)

type AuthRepository interface {
	RevokeAllRefreshTokensForUser(
		ctx context.Context,
		userID int,
	) error
}

type UsersService struct {
	usersRepository UsersRepository
	log             *core_logger.Logger
	authRepository  AuthRepository
}

func NewUsersService(
	usersRepository UsersRepository,
	log *core_logger.Logger,
	authRepository AuthRepository,
) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
		log:             log,
		authRepository:  authRepository,
	}
}

type UsersRepository interface {
	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
		emailVerified *bool,
	) ([]domain.User, error)

	GetUser(
		ctx context.Context,
		id int,
	) (domain.User, error)

	DeleteUser(
		ctx context.Context,
		id int,
	) error

	PatchUser(
		ctx context.Context,
		id int,
		user domain.User,
	) (domain.User, error)

	ChangeUserPassword(
		ctx context.Context,
		user domain.User,
		newPasswordHash string,
	) error

	GetUserByEmail(
		ctx context.Context,
		email string,
	) (domain.User, error)

	DeleteUnverifiedUsersOlderThan(
		ctx context.Context,
		cutoff time.Time,
	) (int, error)
}
