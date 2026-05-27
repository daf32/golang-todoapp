package users_service

import (
	"context"
	"fmt"

	core_auth "github.com/daf32/golang-todoapp/internal/core/auth"
	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	"go.uber.org/zap"
)

func (s *UsersService) ChangeUserPassword(
	ctx context.Context,
	actor domain.Actor,
	userID int,
	password core_auth.PlainPassword,
	newPassword core_auth.PlainPassword,
	confirmPassword core_auth.PlainPassword,
) error {
	if err := authorizeUserAccess(actor, userID); err != nil {
		return err
	}

	user, err := s.usersRepository.GetUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user id=%d: %w", userID, err)
	}

	if err := password.Validate(); err != nil {
		return fmt.Errorf("validate password: %w", err)
	}

	if err := confirmPassword.Validate(); err != nil {
		return fmt.Errorf("validate confirm password: %w", err)
	}

	if err := newPassword.Validate(); err != nil {
		return fmt.Errorf("validate new password: %w", err)
	}

	if err := core_auth.VerifyPassword(user.PasswordHash, password); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	if string(confirmPassword) != string(newPassword) {
		return fmt.Errorf(
			"`new_password` and `confirm_password` do not match: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	newPasswordHash, err := newPassword.Hash()
	if err != nil {
		return fmt.Errorf("generate new password hash: %w", err)
	}

	if err := s.usersRepository.ChangeUserPassword(
		ctx,
		user,
		newPasswordHash,
	); err != nil {
		return fmt.Errorf("change user password: %w", err)
	}

	if err := s.authRepository.RevokeAllRefreshTokensForUser(ctx, userID); err != nil {
		s.log.Error("revoke sessions after password change failed",
			zap.Int("user_id", userID), zap.Error(err))
	}

	return nil
}
