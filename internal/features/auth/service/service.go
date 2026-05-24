package auth_service

import (
	"context"
	"time"

	core_auth "github.com/daf32/golang-todoapp/internal/core/auth"
	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_logger "github.com/daf32/golang-todoapp/internal/core/logger"
	core_mailer "github.com/daf32/golang-todoapp/internal/core/mailer"
	"github.com/golang-jwt/jwt/v5"
)

type AuthRepository interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	CreateRefreshToken(
		ctx context.Context,
		userID int,
		ttl time.Duration,
	) (core_auth.RefreshToken, error)

	GetRefreshToken(
		ctx context.Context,
		tokenString string,
	) (core_auth.RefreshToken, error)

	RevokeRefreshToken(
		ctx context.Context,
		tokenString string,
	) error

	CreateEmailConfirmationToken(
		ctx context.Context,
		userID int,
		ttl time.Duration,
	) (domain.EmailConfirmationToken, error)

	GetAndConsumeEmailConfirmationToken(
		ctx context.Context,
		token string,
	) (domain.EmailConfirmationToken, error)
}

type UsersRepository interface {
	GetUser(
		ctx context.Context,
		id int,
	) (domain.User, error)

	GetUserByEmail(
		ctx context.Context,
		email string,
	) (domain.User, error)
}

type AuthService struct {
	authRepository            AuthRepository
	usersRepository           UsersRepository
	mailer                    core_mailer.Mailer
	log                       *core_logger.Logger
	jwtSecret                 []byte
	accessTokenTTL            time.Duration
	refreshTokenTTL           time.Duration
	emailConfirmationTokenTTL time.Duration
}

func NewAuthService(
	authRepository AuthRepository,
	usersRepository UsersRepository,
	mailer core_mailer.Mailer,
	log *core_logger.Logger,
	jwtSecret string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
	emailConfirmationTokenTTL time.Duration,
) *AuthService {
	return &AuthService{
		authRepository:            authRepository,
		usersRepository:           usersRepository,
		mailer:                    mailer,
		log:                       log,
		jwtSecret:                 []byte(jwtSecret),
		accessTokenTTL:            accessTokenTTL,
		refreshTokenTTL:           refreshTokenTTL,
		emailConfirmationTokenTTL: emailConfirmationTokenTTL,
	}
}

func (s *AuthService) generateAccessToken(user domain.User) (string, error) {
	expirationTime := time.Now().Add(s.accessTokenTTL)

	claims := jwt.MapClaims{
		"sub":          user.ID,
		"role":         string(user.Role),
		"full_name":    user.FullName,
		"phone_number": user.PhoneNumber,
		"email":        user.Email,
		"exp":          expirationTime.Unix(),
		"iat":          time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
