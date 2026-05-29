package core_dto

import (
	"time"

	"github.com/daf32/golang-todoapp/internal/core/domain"
)

type UserDTOResponse struct {
	ID              int        `json:"id" example:"10"`
	Version         int        `json:"version" example:"3"`
	FullName        string     `json:"full_name" example:"Ivan Ivanov"`
	PhoneNumber     *string    `json:"phone_number" example:"+79998887766"`
	Email           string     `json:"email" example:"user@example.com"`
	Role            string     `json:"role" example:"user"`
	EmailVerified   bool       `json:"email_verified" example:"true"`
	EmailVerifiedAt *time.Time `json:"email_verified_at" example:"2026-01-15T10:30:00Z"`
	CreatedAt       time.Time  `json:"created_at" example:"2026-01-10T09:00:00Z"`
}

func UserDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:              user.ID,
		Version:         user.Version,
		FullName:        user.FullName,
		PhoneNumber:     user.PhoneNumber,
		Email:           user.Email,
		Role:            string(user.Role),
		EmailVerified:   user.EmailVerified,
		EmailVerifiedAt: user.EmailVerifiedAt,
		CreatedAt:       user.CreatedAt,
	}
}

func UsersDTOFromDomains(users []domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(users))

	for i, user := range users {
		usersDTO[i] = UserDTOFromDomain(user)
	}

	return usersDTO
}
