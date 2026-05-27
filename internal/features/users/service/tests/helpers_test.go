package users_service_test

import "github.com/daf32/golang-todoapp/internal/core/domain"

func typePtr[T any](v T) *T {
	return &v
}

func validUser() domain.User {
	return domain.User{
		ID:           1,
		Version:      1,
		FullName:     "John Doe",
		PhoneNumber:  nil,
		Email:        "john.doe@example.com",
		PasswordHash: "bcrypt-hash-placeholder",
		Role:         domain.UserRoleUser,
	}
}

func validUserActor() domain.Actor {
	return domain.Actor{
		UserID: 1,
		Role:   domain.UserRoleUser,
	}
}

func validAdminActor() domain.Actor {
	return domain.Actor{
		UserID: 999,
		Role:   domain.UserRoleAdmin,
	}
}
