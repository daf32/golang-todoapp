package auth_transport_http

import (
	"net/http"

	core_logger "github.com/daf32/golang-todoapp/internal/core/logger"
	core_http_request "github.com/daf32/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/daf32/golang-todoapp/internal/core/transport/http/response"
)

type ResendConfirmationRequest struct {
	Email string `json:"email" validate:"required,email" example:"user@example.com"`
}

// ResendConfirmationEmail  godoc
// @Summary                 Resend confirmation email
// @Description             Sends a new confirmation email if the address exists and is not yet verified
// @Tags                    auth
// @Accept                  json
// @Param                   request body ResendConfirmationEmailRequest true "Email address"
// @Success                 204
// @Failure                 400 {object} core_http_response.ErrorResponse "Bad request"
// @Router                  /auth/resend-confirmation [post]
func (h *AuthHTTPHandler) ResendConfirmationEmail(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request ResendConfirmationRequest
	if err := core_http_request.DecodeEndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)
	}

	confirmationURL := h.appBaseURL + h.apiVersion.Path(confirmEmailPath) + "?token="

	_ = h.authService.ResendConfirmationEmail(ctx, request.Email, confirmationURL)

	responseHandler.NoContentResponse()
}
