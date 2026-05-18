package auth_service

import (
	"errors"
	"fmt"

	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	"github.com/golang-jwt/jwt/v5"
)

func (s *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token: %w", core_errors.ErrInvalidCredentials)
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token expired: %w", core_errors.ErrInvalidCredentials)
		}
		return nil, fmt.Errorf("invalid token: %w", core_errors.ErrInvalidCredentials)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token: %w", core_errors.ErrInvalidCredentials)
}
