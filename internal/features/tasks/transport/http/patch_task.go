package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_logger "github.com/daf32/golang-todoapp/internal/core/logger"
	core_http_request "github.com/daf32/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/daf32/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/daf32/golang-todoapp/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title" swaggertype:"string" example:"Play basketball"`
	Description core_http_types.Nullable[string] `json:"description" swaggertype:"string" example:"null"`
	Completed   core_http_types.Nullable[bool]   `json:"completed" swaggertype:"boolean"`
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("`Title` can't be null")
		}

		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("`Title` must be between 1 and 100 symbols")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLen := len([]rune(*r.Description.Value))
			if descriptionLen < 1 || descriptionLen > 1000 {
				return fmt.Errorf("`Description` must be betweeen 1 and 1000 symbols")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("`Completed` can't be NULL")
		}
	}

	return nil
}

// PatchTask 	 godoc
// @Summary 	 Change task
// @Description  Change existing task information
// @Description  ### Logic update fields (Three-state logic):
// @Description  1. **The field is not transmitted**: `description` ignored, the value in the database does not change
// @Description  2. **Passed value**: `"description`": "Play basketball at 6 pm"` - set new description value
// @Description  3. **Passed null**: `"description`": "null"` - clear a field in the database (set to NULL)
// @Description  Restrictions: `title ` and `completed` can't be null
// @Tags 		 tasks
// @Accept 		 json
// @Produce		 json
// @Param 		 id path int true "Task id to change"
// @Param 		 request body PatchTaskRequest true "PatchTask request body"
// @Success 	 200 {object} PatchUserResponse "Seccessfull changed task"
// @Failure 	 400 {object} core_http_response.ErrorResponse  "Bad request"
// @Failure 	 404 {object} core_http_response.ErrorResponse  "Task not found"
// @Failure 	 409 {object} core_http_response.ErrorResponse  "Conflict"
// @Failure 	 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 	 	 /tasks/{id} [patch]
func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetIntPathValues(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID path value",
		)

		return
	}

	var request PatchTaskRequest
	if err := core_http_request.DecodeEndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"decode and validate HTTP request",
		)

		return
	}

	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch task",
		)

		return
	}

	response := PatchUserResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}

type PatchUserResponse TaskDTOResponse

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
