package users_transport_http

import (
	"net/http"

	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	core_logger "github.com/daf32/golang-todoapp/internal/core/logger"
	core_dto "github.com/daf32/golang-todoapp/internal/core/transport/http/dto"
	core_http_middleware "github.com/daf32/golang-todoapp/internal/core/transport/http/middleware"
	core_http_request "github.com/daf32/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/daf32/golang-todoapp/internal/core/transport/http/response"
)

type GetUserResponse core_dto.UserDTOResponse

// GetUser 	 	 godoc
// @Summary 	 Get user
// @Description  Get a specific user by ID
// @Tags 		 users
// @Produce		 json
// @Security 	 BearerAuth
// @Param 		 id path int true "ID of the received user"
// @Success 	 200 {object} GetUserResponse "User seccessfuly found"
// @Failure 	 400 {object} core_http_response.ErrorResponse  "Bad request"
// @Failure 	 401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure 	 403 {object} core_http_response.ErrorResponse "Forbidden"
// @Failure 	 404 {object} core_http_response.ErrorResponse  "User not found"
// @Failure 	 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 	 	 /users/{id} [get]
func (h *UsersHTTPHanlder) GetUser(rw http.ResponseWriter, r *http.Request) {
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

	user, err := h.userService.GetUser(ctx, actor, userID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user",
		)

		return
	}

	response := GetUserResponse(core_dto.UserDTOFromDomain(user))

	responseHandler.JSONResponse(response, http.StatusOK)
}
