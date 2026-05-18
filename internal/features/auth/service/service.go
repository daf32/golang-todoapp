package auth_service

import (
	"context"
	"time"

	"github.com/daf32/golang-todoapp/internal/core/domain"
	users_postgres_repository "github.com/daf32/golang-todoapp/internal/features/users/repository/postgres"
	"github.com/golang-jwt/jwt/v5"
)

type RefreshTokenRepository interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	GetUserByEmail(
		ctx context.Context,
		email string,
	) (domain.User, error)

	CreateRefreshToken(
		ctx context.Context,
		userID int,
		ttl time.Duration,
	) (domain.RefreshToken, error)

	GetRefreshToken(
		ctx context.Context,
		tokenString string,
	) (domain.RefreshToken, error)

	RevokeRefreshToken(
		ctx context.Context,
		tokenString string,
	) error
}

type AuthService struct {
	refreshTokenRepository RefreshTokenRepository
	usersRepository        users_postgres_repository.UsersRepository
	jwtSecret              []byte
	accessTokenTTL         time.Duration
	refreshTokenTTL        time.Duration
}

func NewAuthService(
	refreshTokenRepository RefreshTokenRepository,
	usersRepository users_postgres_repository.UsersRepository,
	jwtSecret string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *AuthService {
	return &AuthService{
		refreshTokenRepository: refreshTokenRepository,
		usersRepository:        usersRepository,
		jwtSecret:              []byte(jwtSecret),
		accessTokenTTL:         accessTokenTTL,
		refreshTokenTTL:        refreshTokenTTL,
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
