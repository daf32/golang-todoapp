package tasks_transport_http

import (
	"net/http"
	"time"

	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	core_logger "github.com/daf32/golang-todoapp/internal/core/logger"
	core_http_middleware "github.com/daf32/golang-todoapp/internal/core/transport/http/middleware"
	core_http_request "github.com/daf32/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/daf32/golang-todoapp/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title       string    `json:"title"       validate:"required,min=1,max=100"   example:"basketball"`
	Description *string   `json:"description" validate:"omitempty,min=1,max=1000" example:"play_basketball"`
	Date        time.Time `json:"date"        validate:"required"                 example:"2026-05-26T10:55:34Z"`
}

type CreateTaskResponse TaskDTOResponse

// CreateTask 	 godoc
// @Summary 	 Create task
// @Description  Create new task in system
// @Tags 		 tasks
// @Accept 		 json
// @Produce 	 json
// @Security 	 BearerAuth
// @Param 		 request body CreateTaskRequest true "CreateTask request body"
// @Success 	 201 {object} CreateTaskResponse "Seccessfull created task"
// @Failure 	 400 {obj ect} core_http_response.ErrorResponse  "Bad request"
// @Failure 	 401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure 	 403 {object} core_http_response.ErrorResponse "Forbidden"
// @Failure 	 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 	 	 /tasks [post]
func (h *TasksHTTPHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
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

	var request CreateTaskRequest
	if err := core_http_request.DecodeEndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	taskDomain := domain.NewTaskUninitialized(
		request.Title,
		request.Description,
		actor.UserID,
		request.Date,
	)

	taskDomain, err := h.tasksService.CreateTask(ctx, actor, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create task",
		)

		return
	}

	response := CreateTaskResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusCreated)
}
