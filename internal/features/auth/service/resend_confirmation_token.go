package auth_service

import (
	"context"
)

func (s *AuthService) ResendConfirmationEmail(
	ctx context.Context,
	email string,
	confirmationURL string,
) error {
	user, err := s.usersRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil
	}

	if user.EmailVerified {
		return nil
	}

	go s.sendConfirmationEmail(user, confirmationURL)

	return nil
}
