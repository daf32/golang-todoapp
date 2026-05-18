package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	core_logger "github.com/daf32/golang-todoapp/internal/core/logger"
	core_http_middleware "github.com/daf32/golang-todoapp/internal/core/transport/http/middleware"
	core_http_request "github.com/daf32/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/daf32/golang-todoapp/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

// GetTasks 	 godoc
// @Summary 	 Tasks list
// @Description  View tasks list with optional pagination
// @Tags 		 tasks
// @Produce		 json
// @Security 	 BearerAuth
// @Param 		 user_id query int false "Filter by task id"
// @Param 		 limit query int false "Tasks page size"
// @Param 		 offset query int false "Tasks page shifting"
// @Success 	 200 {object} GetTasksResponse "Seccessfull get a list of tasks"
// @Failure 	 400 {object} core_http_response.ErrorResponse  "Bad request"
// @Failure 	 401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure 	 403 {object} core_http_response.ErrorResponse "Forbidden"
// @Failure 	 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 	 	 /tasks [get]
func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
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

	userID, limit, offset, err := getUserIDLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID/limit/offset query params",
		)

		return
	}

	tasksDomains, err := h.tasksService.GetTasks(
		ctx,
		actor,
		userID,
		limit,
		offset,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get tasks",
		)

		return
	}

	response := GetTasksResponse(taskDTOsFromDomains(tasksDomains))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func getUserIDLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
	const (
		userIDQueryParamKey = "user_id"
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param key: %w", err)
	}

	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get limit query param: %w", err)
	}
	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get offset query param: %w", err)
	}

	return userID, limit, offset, nil
}
