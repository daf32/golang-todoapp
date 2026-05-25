package auth_service_test

import (
	"time"

	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_mailer "github.com/daf32/golang-todoapp/internal/core/mailer"
	core_oauth "github.com/daf32/golang-todoapp/internal/core/oauth"
	"github.com/stretchr/testify/mock"
)

func validVerifiedUser() domain.User {
	return domain.User{
		ID:            1,
		Version:       1,
		FullName:      "John Doe",
		Email:         "john.doe@example.com",
		PasswordHash:  "bcrypt-hash-placeholder",
		Role:          domain.UserRoleUser,
		EmailVerified: true,
	}
}

func ValidUnverifiedUser() domain.User {
	return domain.User{
		ID:            2,
		Version:       1,
		FullName:      "Jane Doe",
		Email:         "jane.doe@example.com",
		PasswordHash:  "bcrypt-hash-placeholder",
		Role:          domain.UserRoleUser,
		EmailVerified: false,
	}
}

func validConfirmationToken(userID int) domain.EmailConfirmationToken {
	return domain.EmailConfirmationToken{
		Token:     "plain-token-value",
		UserID:    userID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
}

func newNoOpMailerMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *core_mailer.MockMailer {
	return core_mailer.NewMockMailer(t)
}

func validOAuthUserInfo() core_oauth.UserInfo {
	return core_oauth.UserInfo{
		Sub:           "google-sub-12345",
		Email:         "oauth.user@example.com",
		EmailVerified: true,
		Name:          "OAuth User",
	}
}

func newMockGoogleProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *core_oauth.MockProvider {
	p := core_oauth.NewMockProvider(t)
	p.On("Name").Return(core_oauth.ProviderGoogle).Maybe()
	return p
}
