package tasks_service_test

import "github.com/daf32/golang-todoapp/internal/core/domain"

func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func validTask() domain.Task {
	return domain.Task{
		ID:           1,
		Version:      1,
		Title:        "test_task",
		Description:  nil,
		Completed:    false,
		CreatedAt:    createdAt,
		CompletedAt:  nil,
		AuthorUserID: 1,
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
