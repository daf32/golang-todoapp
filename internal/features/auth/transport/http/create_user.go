package auth_transport_http

import (
	"net/http"

	core_auth "github.com/daf32/golang-todoapp/internal/core/auth"
	"github.com/daf32/golang-todoapp/internal/core/domain"
	core_logger "github.com/daf32/golang-todoapp/internal/core/logger"
	core_dto "github.com/daf32/golang-todoapp/internal/core/transport/http/dto"
	core_http_request "github.com/daf32/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/daf32/golang-todoapp/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name"     validate:"required,min=3,max=100"                 example:"user_name"`
	PhoneNumber *string `json:"phone_number"  validate:"omitempty,min=10,max=15,startswith=+"   example:"+79998887766"`
	Email       string  `json:"email"         validate:"email,required,min=5,max=255"           example:"user@example.com"`
	Password    string  `json:"password"      validate:"required,min=8,max=72"                 example:"password_example"`
}

type CreateUserResponse core_dto.UserDTOResponse

// CreateUser 	 godoc
// @Summary 	 Create user
// @Description  Create new user in system
// @Tags 		 auth
// @Accept 		 json
// @Produce 	 json
// @Param 		 request body CreateUserRequest true "CreateUser request body"
// @Success 	 201 {object} CreateUserResponse "Seccessfull created user"
// @Failure 	 400 {object} core_http_response.ErrorResponse  "Bad request"
// @Failure 	 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 	 	 /auth/register [post]
func (h *AuthHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request CreateUserRequest
	if err := core_http_request.DecodeEndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	confirmationURL := h.appBaseURL + h.apiVersion.Path(confirmEmailPath) + "?token="

	userDomain := domainFromDTO(request)

	userDomain, err := h.authService.CreateUser(
		ctx,
		userDomain,
		core_auth.PlainPassword(request.Password),
		confirmationURL,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create user",
		)

		return
	}

	response := CreateUserResponse(core_dto.UserDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(
		dto.FullName,
		dto.PhoneNumber,
		dto.Email,
		domain.UserRoleUser,
	)
}
