package auth_transport_http

import (
	"net/http"

	core_logger "github.com/daf32/golang-todoapp/internal/core/logger"
	core_http_request "github.com/daf32/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/daf32/golang-todoapp/internal/core/transport/http/response"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,min=1,max=255" example:"refresh_token"`
}

type RefreshResponse struct {
	Token string `json:"access_token" example:"access_token"`
}

// RefreshToken  godoc
// @Summary 	 Refresh tokens
// @Description  Obtain a new access token
// @Tags 		 auth
// @Accept 		 json
// @Produce 	 json
// @Param 		 request body RefreshRequest   true             "RefreshToken request body"
// @Success 	 200 {object} RefreshResponse                   "Seccessfull refresh token"
// @Failure 	 400 {object} core_http_response.ErrorResponse  "Bad request"
// @Failure 	 401 {object} core_http_response.ErrorResponse  "Invalid credentials"
// @Failure 	 404 {object} core_http_response.ErrorResponse  "User not found"
// @Failure 	 500 {object} core_http_response.ErrorResponse  "Internal server error"
// @Router 	 	 /auth/refresh [post]
func (h *AuthHTTPHandler) RefreshToken(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request RefreshRequest
	if err := core_http_request.DecodeEndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	token, err := h.authService.RefreshAccessToken(ctx, request.RefreshToken)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to refresh access token",
		)

		return
	}

	response := RefreshResponse{Token: token}
	responseHandler.JSONResponse(response, http.StatusOK)
}
