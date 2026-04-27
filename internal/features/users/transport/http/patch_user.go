package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_logger "github.com/daf32/golang-todoapp/internal/core/logger"
	core_http_request "github.com/daf32/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/daf32/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/daf32/golang-todoapp/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name" swaggertype:"string" example:"Jebron Lames"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+71112223344"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("`FullName` can't be NULL")
		}

		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("`FullName` must be between 3 and 100 symbols")
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*r.PhoneNumber.Value))

			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("`PhoneNumber` must be between 10 and 15 symbols")
			}

			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("`PhoneNumber` must startswith '+' symbol")
			}
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

// PatchUser 	 godoc
// @Summary 	 Change user
// @Description  Change existing user information
// @Description  ### Logic update fields (Three-state logic):
// @Description  1. **The field is not transmitted**: `phone_number` ignored, the value in the database does not change
// @Description  2. **Passed value**: `"phone_number`": "+71112223344"` - set a new phone number in the database
// @Description  3. **Passed null**: `"phone_number`": "null"` - clear a field in the database (set to NULL)
// @Description  Restrictions: `full_name ` can't be null
// @Tags 		 users
// @Accept 		 json
// @Produce		 json
// @Param 		 id path int true "User id to change"
// @Param 		 request body PatchUserRequest true "PatchUser request body"
// @Success 	 200 {object} PatchUserResponse "Seccessfull changed user"
// @Failure 	 400 {object} core_http_response.ErrorResponse  "Bad request"
// @Failure 	 404 {object} core_http_response.ErrorResponse  "User not found"
// @Failure 	 409 {object} core_http_response.ErrorResponse  "Conflict"
// @Failure 	 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 	 	 /users/{id} [patch]
func (h *UsersHTTPHanlder) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathValues(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)

		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeEndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.userService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)

		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FullName.ToDomain(),
		request.PhoneNumber.ToDomain(),
	)
}
