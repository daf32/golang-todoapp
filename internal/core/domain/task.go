package domain

import (
	"fmt"
	"time"

	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
)

type Task struct {
	ID      int
	Version int

	Title       string
	Description *string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time

	AuthorUserID int
}

func NewTask(
	id int,
	version int,
	title string,
	description *string,
	compleated bool,
	createdAt time.Time,
	compleatedAt *time.Time,
	authorUserId int,
) Task {
	return Task{
		ID:           id,
		Version:      version,
		Title:        title,
		Description:  description,
		Completed:    compleated,
		CreatedAt:    createdAt,
		CompletedAt:  compleatedAt,
		AuthorUserID: authorUserId,
	}
}

func NewTaskUninitialized(
	title string,
	description *string,
	authorUserID int,
) Task {
	return NewTask(
		UninitializedID,
		UninitializedVersion,
		title,
		description,
		false,
		time.Now(),
		nil,
		authorUserID,
	)
}

func (t *Task) Validate() error {
	titleLen := len([]rune(t.Title))
	if titleLen < 1 || titleLen > 100 {
		return fmt.Errorf(
			"invalid title len: %d: %w",
			titleLen,
			core_errors.ErrInvalidArgument,
		)
	}

	if t.Description != nil {
		descriptionLen := len([]rune(*t.Description))
		if descriptionLen < 1 || descriptionLen > 1000 {
			return fmt.Errorf(
				"invalid description len: %d: %w",
				descriptionLen,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if t.Completed {
		if t.CompletedAt == nil {
			return fmt.Errorf(
				"`CompleatedAt` can't be `nil` if `Completed`==`true: %w`",
				core_errors.ErrInvalidArgument,
			)
		}

		if t.CompletedAt.Before(*t.CompletedAt) {
			return fmt.Errorf(
				"`CompleatedAt` can't be before `CreatedAt`: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	} else {
		if t.CompletedAt != nil {
			return fmt.Errorf(
				"`CompleateAt` must be `nil` if `CompleateAt` == `false`: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}

type TaskPatch struct {
	Title       Nullable[string]
	Description Nullable[string]
	Completed   Nullable[bool]
}

func NewTaskPatch(
	title Nullable[string],
	description Nullable[string],
	completed Nullable[bool],
) TaskPatch {
	return TaskPatch{
		Title:       title,
		Description: description,
		Completed:   completed,
	}
}

func (p *TaskPatch) Validate() error {
	if p.Title.Set && p.Title.Value == nil {
		return fmt.Errorf(
			"`Title` can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	
	if p.Completed.Set && p.Completed.Value == nil {
		return fmt.Errorf(
			"`Completed` can't be patche to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	
	return nil
}

func (t *Task) ApplyPatch(patch TaskPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate task patch: %w", err)
	}
	
	tmp := *t
	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}
	
	if patch.Description.Set {
		tmp.Description = patch.Description.Value
	}
	
	if patch.Completed.Set {
		tmp.Completed = *patch.Completed.Value
	
		if tmp.Completed {
			completedAt := time.Now()
			tmp.CompletedAt = &completedAt
		} else {
			tmp.CompletedAt = nil
		}
	}
	
	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched task: %w", err)
	}
	
	*t = tmp
	
	return nil
}
