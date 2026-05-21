package users_transport_http

import (
	"fmt"
	"net/http"

	core_auth "github.com/daf32/golang-todoapp/internal/core/auth"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	core_logger "github.com/daf32/golang-todoapp/internal/core/logger"
	core_http_middleware "github.com/daf32/golang-todoapp/internal/core/transport/http/middleware"
	core_http_request "github.com/daf32/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/daf32/golang-todoapp/internal/core/transport/http/response"
)

type ChangePasswordRequest struct {
	Password        string `json:"password" validate:"required,min=8,max=72" example:"password"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,max=72" example:"confirm_password"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=72" example:"new_password"`
}

func (r *ChangePasswordRequest) Validate() error {
	if r.NewPassword != r.ConfirmPassword {
		return fmt.Errorf("`new_password` and `confirm_password` do not match")
	}
	return nil
}

// ChangeUserPassword godoc
// @Summary      Change user password
// @Description  Change the password of an existing user.
// @Description  The caller must provide the current password and the new password (confirmed twice).
// @Description  Regular users can only change their own password. Admins can change any user's password.
// @Tags         users
// @Accept       json
// @Security     BearerAuth
// @Param        id      path  int                    true  "User ID"
// @Param        request body  ChangePasswordRequest  true  "ChangeUserPassword request body"
// @Success      204
// @Failure      400  {object}  core_http_response.ErrorResponse  "Bad request (validation failed)"
// @Failure      401  {object}  core_http_response.ErrorResponse  "Unauthorized"
// @Failure      403  {object}  core_http_response.ErrorResponse  "Forbidden"
// @Failure      404  {object}  core_http_response.ErrorResponse  "User not found"
// @Failure      500  {object}  core_http_response.ErrorResponse  "Internal server error"
// @Router       /users/{id}/password [patch]
func (h *UsersHTTPHanlder) ChangeUserPassword(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	actor, ok := core_http_middleware.GetActor(r)
	if !ok {
		responseHandler.ErrorResponse(
			core_errors.ErrInvalidCredentials,
			"failed to get authenticated actor from request context",
		)

		return
	}

	userID, err := core_http_request.GetIntPathValues(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)

		return
	}

	var request ChangePasswordRequest
	if err := core_http_request.DecodeEndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	if err := h.userService.ChangeUserPassword(
		ctx,
		actor,
		userID,
		core_auth.PlainPassword(request.Password),
		core_auth.PlainPassword(request.NewPassword),
		core_auth.PlainPassword(request.ConfirmPassword),
	); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to change password",
		)

		return
	}

	responseHandler.NoContentResponse()
}
