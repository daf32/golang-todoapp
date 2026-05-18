package tasks_transport_http

import (
	"net/http"

	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	core_logger "github.com/daf32/golang-todoapp/internal/core/logger"
	core_http_middleware "github.com/daf32/golang-todoapp/internal/core/transport/http/middleware"
	core_http_request "github.com/daf32/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/daf32/golang-todoapp/internal/core/transport/http/response"
)

type GetTaskResponse TaskDTOResponse

// GetTask 	 	 godoc
// @Summary 	 Get task
// @Description  Get a specific task by ID
// @Tags 		 tasks
// @Produce		 json
// @Security 	 BearerAuth
// @Param 		 id path int true "ID of the received task"
// @Success 	 200 {object} GetTaskResponse "Task seccessfuly found"
// @Failure 	 400 {object} core_http_response.ErrorResponse  "Bad request"
// @Failure 	 401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure 	 403 {object} core_http_response.ErrorResponse "Forbidden"
// @Failure 	 404 {object} core_http_response.ErrorResponse  "Task not found"
// @Failure 	 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 	 	 /tasks/{id} [get]
func (h *TasksHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
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

	taskID, err := core_http_request.GetIntPathValues(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID path value",
		)

		return
	}

	taskDomain, err := h.tasksService.GetTask(ctx, actor, taskID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get task",
		)

		return
	}

	response := GetTaskResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}
