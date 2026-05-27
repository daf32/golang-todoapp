package core_auth

import (
	"fmt"
	"unicode"

	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	"golang.org/x/crypto/bcrypt"
)

type PlainPassword string

func (p PlainPassword) Hash() (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(string(p)),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

func (p PlainPassword) Validate() error {
	const minLen, maxLen = 8, 72

	l := len([]byte(string(p)))
	if l < minLen && l > maxLen {
		return fmt.Errorf(
			"password length must be %d–%d chars, got %d: %w",
			minLen, maxLen, l,
			core_errors.ErrInvalidArgument,
		)
	}

	var hasUpper, hasLower, hasDigit bool
	for _, r := range string(p) {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter: %w",
			core_errors.ErrInvalidArgument)
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter: %w",
			core_errors.ErrInvalidArgument)
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit: %w",
			core_errors.ErrInvalidArgument)
	}

	return nil
}

func VerifyPassword(hashedPassword string, providedPassword PlainPassword) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(string(providedPassword)))
}
