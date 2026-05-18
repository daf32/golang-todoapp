package auth_transport_http

import (
	"time"

	"github.com/daf32/golang-todoapp/internal/core/domain"
	"github.com/google/uuid"
)

type RefreshTokenDTOResponse struct {
	ID        uuid.UUID `json:"id" example:"1"`
	UserID    int       `json:"user_id" example:"2"`
	Token     string    `json:"token" example:"token"`
	ExpiresAt time.Time `json:"expires_at" example:"2026-02-26T10:30:00Z"`
	CreatedAt time.Time `json:"created_at" example:"2026-02-26T10:30:00Z"`
	Revoked   bool      `json:"revoked" example:"false"`
}

func refreshTokenDTOFromDomain(refreshToken domain.RefreshToken) RefreshTokenDTOResponse {
	return RefreshTokenDTOResponse{
		ID: refreshToken.ID,
		UserID: refreshToken.UserID,
		Token: refreshToken.Token,
		ExpiresAt: refreshToken.ExpiresAt,
		CreatedAt: refreshToken.CreatedAt,
		Revoked: refreshToken.Revoked,
	}
}
