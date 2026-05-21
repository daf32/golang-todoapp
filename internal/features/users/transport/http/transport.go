package users_transport_http

import (
	"context"
	"net/http"

	core_auth "github.com/daf32/golang-todoapp/internal/core/auth"
	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_http_middleware "github.com/daf32/golang-todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/daf32/golang-todoapp/internal/core/transport/http/server"
)

type UsersHTTPHanlder struct {
	userService UsersService
}

type UsersService interface {
	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)

	GetUser(
		ctx context.Context,
		actor domain.Actor,
		id int,
	) (domain.User, error)

	DeleteUser(
		ctx context.Context,
		id int,
	) error

	PatchUser(
		ctx context.Context,
		actor domain.Actor,
		id int,
		patch domain.UserPatch,
	) (domain.User, error)

	ChangeUserPassword(
		ctx context.Context,
		actor domain.Actor,
		userID int,
		password core_auth.PlainPassword,
		newPassword core_auth.PlainPassword,
		confirmPassword core_auth.PlainPassword,
	) error
}

func NewUsersHTTPHanlder(
	usersService UsersService,
) *UsersHTTPHanlder {
	return &UsersHTTPHanlder{
		userService: usersService,
	}
}

func (h *UsersHTTPHanlder) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: h.GetUsers,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.RequireRole(domain.UserRoleAdmin),
			},
			/*
				 	Example of usage Middleware on separate Route

					Middleware: []core_http_middleware.Middleware {
						core_http_middleware.Dummy("get users middleware"),
					},
			*/
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{id}",
			Handler: h.GetUser,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/{id}",
			Handler: h.DeleteUser,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.RequireRole(domain.UserRoleAdmin),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/{id}",
			Handler: h.PatchUser,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/{id}/password",
			Handler: h.ChangeUserPassword,
		},
	}
}
