package auth_transport_http

import (
	"net/http"

	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	core_logger "github.com/daf32/golang-todoapp/internal/core/logger"
	core_http_middleware "github.com/daf32/golang-todoapp/internal/core/transport/http/middleware"
	core_http_response "github.com/daf32/golang-todoapp/internal/core/transport/http/response"
)

// LogoutAllUserSessions godoc
// @Summary      Logout all sessions
// @Description  Revoke every refresh token for the authenticated user
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      204 "All sessions revoked"
// @Failure      401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /auth/logout-all [post]
func (h *AuthHTTPHandler) LogoutAllUserSessions(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, ok := core_http_middleware.GetUserID(r)
	if !ok {
		responseHandler.ErrorResponse(
			core_errors.ErrInvalidCredentials,
			"failed to get user id from requesr context",
		)

		return
	}

	if err := h.authService.LogoutAllUserSessions(ctx, userID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to logout all user sessions",
		)

		return
	}

	responseHandler.NoContentResponse()
}
