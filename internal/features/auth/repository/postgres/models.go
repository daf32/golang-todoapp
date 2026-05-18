package auth_postgres_repository

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
