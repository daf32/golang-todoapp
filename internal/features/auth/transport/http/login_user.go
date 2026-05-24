package auth_transport_http

import (
	"net/http"

	core_auth "github.com/daf32/golang-todoapp/internal/core/auth"
	core_logger "github.com/daf32/golang-todoapp/internal/core/logger"
	core_http_request "github.com/daf32/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/daf32/golang-todoapp/internal/core/transport/http/response"
)

type LoginUserRequest struct {
	Email    string `json:"email"     validate:"required,min=5,max=255"  example:"user@example.com"`
	Password string `json:"password"  validate:"required,min=8,max=72"   example:"user_password"`
}

type LoginUserResponse struct {
	AccessToken  string `json:"access_token"   example:"acess_token"`
	RefreshToken string `json:"refresh_token"  example:"refresh_token"`
}

// LoginUer 	 godoc
// @Summary 	 Login user
// @Description  Login user in system
// @Tags 		 auth
// @Accept 		 json
// @Produce 	 json
// @Param 		 request body LoginUserRequest true             "LoginUser request body"
// @Success 	 200 {object} LoginUserResponse                 "Seccessfull login user"
// @Failure 	 400 {object} core_http_response.ErrorResponse  "Bad request"
// @Failure 	 401 {object} core_http_response.ErrorResponse  "Invalid credentials"
// @Failure 	 404 {object} core_http_response.ErrorResponse  "User not found"
// @Failure 	 409 {object} core_http_response.ErrorResponse  "Conflict"
// @Failure 	 500 {object} core_http_response.ErrorResponse  "Internal server error"
// @Router 	 	 /auth/login [post]
func (h *AuthHTTPHandler) LoginUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request LoginUserRequest
	if err := core_http_request.DecodeEndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate request",
		)

		return
	}

	acessToken, refreshToken, err := h.authService.LoginUser(
		ctx,
		request.Email,
		core_auth.PlainPassword(request.Password),
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to login user",
		)

		return
	}

	response := LoginUserResponse{
		AccessToken:  acessToken,
		RefreshToken: refreshToken.Token,
	}
	responseHandler.JSONResponse(response, http.StatusOK)
}
