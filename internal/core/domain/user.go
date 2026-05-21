package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
)

var (
	emailFormatRE = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneNumberRE = regexp.MustCompile(`^\+[0-9]+$`)
)

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

type User struct {
	ID      int
	Version int

	FullName      string
	PhoneNumber   *string
	Email         string
	PasswordHash  string
	Role          UserRole
	EmailVerified bool
}

func NewUser(
	id int,
	version int,
	fullName string,
	phoneNumber *string,
	email string,
	passwordHash string,
	role UserRole,
	emailVerified bool,
) User {
	return User{
		ID:           id,
		Version:      version,
		FullName:     fullName,
		PhoneNumber:  phoneNumber,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		EmailVerified: emailVerified,
	}
}

func NewUserUninitialized(
	fullName string,
	phoneNumber *string,
	email string,
	role UserRole,
) User {
	return NewUser(
		UninitializedID,
		UninitializedVersion,
		fullName,
		phoneNumber,
		email,
		UninitializedPassowrd,
		role,
		UninitializedEmailVerified,
	)
}

func (u *User) Validate() error {
	fullNameLen := len([]rune(u.FullName))
	if fullNameLen < 3 || fullNameLen > 100 {
		return fmt.Errorf(
			"invalid `FullName` len: %d: %w",
			fullNameLen,
			core_errors.ErrInvalidArgument,
		)
	}

	if u.PhoneNumber != nil {
		phoneNumberLen := len([]rune(*u.PhoneNumber))
		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf(
				"invalid `PhoneNumber` len: %d: %w",
				phoneNumberLen,
				core_errors.ErrInvalidArgument,
			)
		}

		if !phoneNumberRE.MatchString(*u.PhoneNumber) {
			return fmt.Errorf(
				"invalid `PhoneNumber` format: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	emailLen := len([]rune(u.Email))
	if emailLen < 5 || emailLen > 255 {
		return fmt.Errorf(
			"invalid `Email` len: %d: %w",
			emailLen,
			core_errors.ErrInvalidArgument,
		)
	}

	if !emailFormatRE.MatchString(u.Email) {
		return fmt.Errorf(
			"invalid `Email` format: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if len(u.PasswordHash) == 0 {
		return fmt.Errorf(
			"empty `PasswordHash`: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if u.Role != UserRoleUser && u.Role != UserRoleAdmin {
		return fmt.Errorf(
			"invalid `Role`: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
	Email       Nullable[string]
}

func NewUserPatch(
	fullName Nullable[string],
	phoneNumber Nullable[string],
	email Nullable[string],
) UserPatch {
	return UserPatch{
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		Email:       email,
	}
}

func (p *UserPatch) Validate() error {
	if p.FullName.Set && p.FullName.Value == nil {
		return fmt.Errorf(
			"`FullName` can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	if p.Email.Set && p.Email.Value == nil {
		return fmt.Errorf(
			"`Email` can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.Email.Set && p.Email.Value != nil {
		emailLen := len([]rune(*p.Email.Value))
		if emailLen < 5 || emailLen > 255 {
			return fmt.Errorf(
				"invalid patched `Email` len: %d: %w",
				emailLen,
				core_errors.ErrInvalidArgument,
			)
		}

		if !emailFormatRE.MatchString(*p.Email.Value) {
			return fmt.Errorf(
				"invalid patched `Email` format: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf(
			"validate user patch: %w",
			err,
		)
	}

	tmp := *u

	if patch.FullName.Set {
		tmp.FullName = *patch.FullName.Value
	}

	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}

	if patch.Email.Set {
		tmp.Email = *patch.Email.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf(
			"validate patched user: %w",
			err,
		)
	}

	*u = tmp

	return nil
}
