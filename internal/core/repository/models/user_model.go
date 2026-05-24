package core_models

import (
	"time"

	"github.com/daf32/golang-todoapp/internal/core/domain"
)

type UserModel struct {
	ID              int
	Version         int
	FullName        string
	PhoneNumber     *string
	Email           string
	PasswordHash    string
	Role            domain.UserRole
	EmailVerified   bool
	EmailVerifiedAt *time.Time
}

func UserDomainsFromModels(users []UserModel) []domain.User {
	usersDomains := make([]domain.User, len(users))

	for i, user := range users {
		usersDomains[i] = domain.NewUser(
			user.ID,
			user.Version,
			user.FullName,
			user.PhoneNumber,
			user.Email,
			user.PasswordHash,
			user.Role,
			user.EmailVerified,
			user.EmailVerifiedAt,
		)
	}

	return usersDomains
}
