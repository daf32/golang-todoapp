package domain

import "time"

type EmailConfirmationToken struct {
	Token     string
	UserID    int
	ExpiresAt time.Time
}

func NewEmailConfirmationToken(
	token string,
	userID int,
	expiresAt time.Time,
) EmailConfirmationToken {
	return EmailConfirmationToken{
		Token: token,
		UserID: userID,
		ExpiresAt: expiresAt,
	}
}
