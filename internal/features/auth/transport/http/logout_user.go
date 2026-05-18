package auth_transport_http

import (
	"net/http"

	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	core_logger "github.com/daf32/golang-todoapp/internal/core/logger"
	core_http_middleware "github.com/daf32/golang-todoapp/internal/core/transport/http/middleware"
	core_http_request "github.com/daf32/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/daf32/golang-todoapp/internal/core/transport/http/response"
)

type LogoutUserRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,min=1,max=255" example:"refresh_token"`
}

// LogoutUser  godoc
// @Summary 	 Logout user
// @Description  Revoke the provided refresh token for the authenticated user
// @Tags 		 auth
// @Accept 		 json
// @Produce 	 json
// @Security 	 BearerAuth
// @Param 		 request body LogoutUserRequest true "LogoutUser request body"
// @Success 	 204 "Successful user logout"
// @Failure 	 400 {object} core_http_response.ErrorResponse  "Bad request"
// @Failure 	 401 {object} core_http_response.ErrorResponse  "Unauthorized"
// @Failure 	 403 {object} core_http_response.ErrorResponse  "Forbidden"
// @Failure 	 404 {object} core_http_response.ErrorResponse  "Refresh token not found"
// @Failure 	 500 {object} core_http_response.ErrorResponse  "Internal server error"
// @Router 	 	 /auth/logout [post]
func (h *AuthHTTPHandler) LogoutUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request LogoutUserRequest
	if err := core_http_request.DecodeEndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	userID, ok := core_http_middleware.GetUserID(r)
	if !ok {
		responseHandler.ErrorResponse(
			core_errors.ErrInvalidCredentials,
			"failed to get user id from request context",
		)

		return
	}

	if err := h.authService.LogoutUser(ctx, userID, request.RefreshToken); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to logout user",
		)

		return
	}

	responseHandler.NoContentResponse()
}
