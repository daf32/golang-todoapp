package auth_transport_http

import (
	"context"
	"net/http"

	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_http_middleware "github.com/daf32/golang-todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/daf32/golang-todoapp/internal/core/transport/http/server"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHTTPHandler struct {
	authService AuthService
}

type AuthService interface {
	CreateUser(
		ctx context.Context,
		userDomain domain.User,
		plainPassword string,
	) (domain.User, error)

	LoginUser(
		ctx context.Context,
		email string,
		password string,
	) (string, domain.RefreshToken, error)

	ValidateToken(
		tokenString string,
	) (jwt.MapClaims, error)

	RefreshAccessToken(
		ctx context.Context,
		refreshTokenString string,
	) (string, error)

	LogoutUser(
		ctx context.Context,
		userID int,
		refreshTokenString string,
	) error
}

func NewAuthHTTPHandler(
	authService AuthService,
) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authService: authService,
	}
}

func (h *AuthHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/auth/register",
			Handler: h.CreateUser,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Handler: h.LoginUser,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/refresh",
			Handler: h.RefreshToken,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/logout",
			Handler: h.LogoutUser,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.authService),
			},
		},
	}
}
