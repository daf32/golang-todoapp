package users_service_test

import "github.com/daf32/golang-todoapp/internal/core/domain"

func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func validUser() domain.User {
	return domain.User{
		ID:          1,
		Version:     1,
		FullName:    "John Doe",
		PhoneNumber: nil,
	}
}
