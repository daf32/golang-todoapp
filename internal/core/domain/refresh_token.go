package domain

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID
	UserID    int
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
	Revoked   bool
}

func NewRefreshToken(
	id uuid.UUID,
	userID int,
	token string,
	expiresAt time.Time,
	createdAt time.Time,
	revoked bool,
) RefreshToken {
	return RefreshToken{
		ID:        id,
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: createdAt,
		Revoked:   revoked,
	}
}

func NewRefreshTokenUninitialized(
	userID int,
	token string,
	expiredAt time.Time,
) RefreshToken {
	return NewRefreshToken(
		UninitializedUID,
		userID,
		token,
		expiredAt,
		time.Now(),
		false,
	)
}
