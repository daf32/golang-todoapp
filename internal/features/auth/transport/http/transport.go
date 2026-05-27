package auth_transport_http

import (
	"context"
	"net/http"

	core_auth "github.com/daf32/golang-todoapp/internal/core/auth"
	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_http_cookie "github.com/daf32/golang-todoapp/internal/core/transport/http/cookie"
	core_http_middleware "github.com/daf32/golang-todoapp/internal/core/transport/http/middleware"
	core_ratelimit "github.com/daf32/golang-todoapp/internal/core/transport/http/middleware/ratelimit"
	core_http_server "github.com/daf32/golang-todoapp/internal/core/transport/http/server"
	"github.com/golang-jwt/jwt/v5"
)

type RateLimiters struct {
	Login    core_ratelimit.Limiter
	Register core_ratelimit.Limiter
	Resend   core_ratelimit.Limiter
	Refresh  core_ratelimit.Limiter
}

type AuthHTTPHandler struct {
	authService   AuthService
	apiVersion    core_http_server.ApiVersion
	appBaseURL    string
	cookieManager *core_http_cookie.Manager
	rateLimits    RateLimiters
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

	LoginUserWithOAuth(
		ctx context.Context,
		providerName, code, codeVerifier string,
	) (string, core_auth.RefreshToken, error)

	BuildOAuthURL(
		providerName, state, codeVerifier string,
	) (string, error)

	LogoutAllUserSessions(
		ctx context.Context,
		userID int,
	) error
}

func NewAuthHTTPHandler(
	authService AuthService,
	apiVersion core_http_server.ApiVersion,
	appBaseURL string,
	cookieManager *core_http_cookie.Manager,
	rateLimits RateLimiters,
) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authService:   authService,
		apiVersion:    apiVersion,
		appBaseURL:    appBaseURL,
		cookieManager: cookieManager,
		rateLimits:    rateLimits,
	}
}

const confirmEmailPath = "/auth/confirm-email"

func (h *AuthHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/auth/register",
			Handler: h.CreateUser,
			Middleware: []core_http_middleware.Middleware{
				core_ratelimit.Middleware(h.rateLimits.Register, nil),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Handler: h.LoginUser,
			Middleware: []core_http_middleware.Middleware{
				core_ratelimit.Middleware(h.rateLimits.Login, nil),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/refresh",
			Handler: h.RefreshToken,
			Middleware: []core_http_middleware.Middleware{
				core_ratelimit.Middleware(h.rateLimits.Refresh, nil),
			},
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
			Method:  http.MethodPost,
			Path:    "auth/resend-confirmation",
			Handler: h.ResendConfirmationEmail,
			Middleware: []core_http_middleware.Middleware{
				core_ratelimit.Middleware(h.rateLimits.Resend, nil),
			},
		},

		{
			Method:  http.MethodGet,
			Path:    "/auth/oauth/{provider}",
			Handler: h.StartOAuth,
		},
		{
			Method:  http.MethodGet,
			Path:    "/auth/oauth/{provider}/callback",
			Handler: h.OAuthCallback,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/logout-all",
			Handler: h.LogoutAllUserSessions,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.authService),
			},
		},
	}
}
