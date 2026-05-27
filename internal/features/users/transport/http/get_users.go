package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/daf32/golang-todoapp/internal/core/logger"
	core_dto "github.com/daf32/golang-todoapp/internal/core/transport/http/dto"
	core_http_request "github.com/daf32/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/daf32/golang-todoapp/internal/core/transport/http/response"
)

type GetUsersResponse []core_dto.UserDTOResponse

// GetUsers 	 godoc
// @Summary 	 Users list
// @Description  View users list with optional pagination
// @Tags 		 users
// @Produce		 json
// @Security 	 BearerAuth
// @Param 		 limit query int false "Users page size"
// @Param 		 offset query int false "Users page shifting"
// @Param        email_verified query bool false "Filter by email verification status"
// @Success 	 200 {object} GetUsersResponse "Seccessfull get a list of users"
// @Failure 	 400 {object} core_http_response.ErrorResponse  "Bad request"
// @Failure 	 401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure 	 403 {object} core_http_response.ErrorResponse "Forbidden"
// @Failure 	 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 	 	 /users [get]
func (h *UsersHTTPHanlder) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, emailVerified, err := getLimitOffsetEmailVerifiedQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get `limit`/`offset`/`email_verified` query param",
		)

		return
	}

	userDomains, err := h.userService.GetUsers(ctx, limit, offset, emailVerified)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get users",
		)

		return
	}

	response := GetUsersResponse(core_dto.UsersDTOFromDomains(userDomains))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func getLimitOffsetEmailVerifiedQueryParams(r *http.Request) (*int, *int, *bool, error) {
	const (
		limitQueryParamKey         = "limit"
		offsetQueryParamKey        = "offset"
		emailVerifiedQueryParamKey = "email_verified"
	)

	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get limit query param: %w", err)
	}
	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get offset query param: %w", err)
	}
	emailVerified, err := core_http_request.GetBoolQueryParam(r, emailVerifiedQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get email_verified query param: %w", err)
	}

	return limit, offset, emailVerified, nil
}
