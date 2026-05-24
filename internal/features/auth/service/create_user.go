package auth_service

import (
	"context"
	"fmt"
	"time"

	core_auth "github.com/daf32/golang-todoapp/internal/core/auth"
	"github.com/daf32/golang-todoapp/internal/core/domain"
	"go.uber.org/zap"
)

func (s *AuthService) CreateUser(
	ctx context.Context,
	user domain.User,
	plainPassword core_auth.PlainPassword,
	confirmationURL string,
) (domain.User, error) {
	if err := plainPassword.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate password: %w", err)
	}

	hashedPassword, err := plainPassword.Hash()
	if err != nil {
		return domain.User{}, fmt.Errorf("generate hash password: %w", err)
	}

	user.PasswordHash = hashedPassword

	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}

	user, err = s.authRepository.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	go s.sendConfirmationEmail(user, confirmationURL)

	return user, nil
}

func (s *AuthService) sendConfirmationEmail(
	user domain.User,
	confirmationURL string,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	token, err := s.authRepository.CreateEmailConfirmationToken(ctx, user.ID, s.emailConfirmationTokenTTL)
	if err != nil {
		s.log.Error(
			"failed to create email confirmation token",
			zap.Int("user_id", user.ID),
			zap.Error(err),
		)
	}

	link := confirmationURL + token.Token

	body := fmt.Sprintf(
		"Hi %s,\n\nPlease confirm your email by clicking the link below:\n\n%s\n\nThe link expires in 24 hours.",
		user.FullName,
		link,
	)

	if err := s.mailer.SendEmail(user.Email, "Confirm your email", body); err != nil {
		s.log.Error(
			"failed to send confirmation email",
			zap.String("email", user.Email),
			zap.Int("user_id", user.ID),
			zap.Error(err),
		)
	}
}
