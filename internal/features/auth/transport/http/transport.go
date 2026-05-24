package auth_transport_http

import (
	"context"
	"net/http"

	core_auth "github.com/daf32/golang-todoapp/internal/core/auth"
	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_http_middleware "github.com/daf32/golang-todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/daf32/golang-todoapp/internal/core/transport/http/server"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHTTPHandler struct {
	authService AuthService
	apiVersion  core_http_server.ApiVersion
	appBaseURL  string
}

type AuthService interface {
	CreateUser(
		ctx context.Context,
		userDomain domain.User,
		plainPassword core_auth.PlainPassword,
		confirmationURL string,
	) (domain.User, error)

	LoginUser(
		ctx context.Context,
		email string,
		password core_auth.PlainPassword,
	) (string, core_auth.RefreshToken, error)

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

	ConfirmEmail(
		ctx context.Context,
		token string,
	) error

	ResendConfirmationEmail(
		ctx context.Context,
		email string,
		confirmationURL string,
	) error
}

func NewAuthHTTPHandler(
	authService AuthService,
	apiVersion core_http_server.ApiVersion,
	appBaseURL string,
) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authService: authService,
		apiVersion:  apiVersion,
		appBaseURL:  appBaseURL,
	}
}

const confirmEmailPath = "/auth/confirm-email"

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
		{
			Method:  http.MethodGet,
			Path:    confirmEmailPath,
			Handler: h.ConfirmEmail,
		},
		{
			Method: http.MethodPost,
			Path: "auth/resend-confirmation",
			Handler: h.ResendConfirmationEmail,
		},
	}
}
