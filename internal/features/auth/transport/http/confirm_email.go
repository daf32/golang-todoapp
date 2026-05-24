package auth_transport_http

import (
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	core_logger "github.com/daf32/golang-todoapp/internal/core/logger"
	core_http_response "github.com/daf32/golang-todoapp/internal/core/transport/http/response"
	"go.uber.org/zap"
)

const emailConfirmedRedirectPath = "/email-conifrmed"

// ConfirmEmail   godoc
// @Summary       Confirm email address
// @Description   Validates the one-time token sent by email and marks the user as verified
// @Tags          auth
// @Produce       json
// @Param         token query string true "Email confirmation token"
// @Success       200
// @Failure       400 {object} core_http_response.ErrorResponse "Token expired or invalid"
// @Failure       404 {object} core_http_response.ErrorResponse "Token not found"
// @Failure       500 {object} core_http_response.ErrorResponse "Internal server error"
func (h *AuthHTTPHandler) ConfirmEmail(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	redirectBase := h.appBaseURL + emailConfirmedRedirectPath

	token := r.URL.Query().Get("token")
	if token == "" {
		responseHandler.ErrorResponse(
			fmt.Errorf("token query param is required: %w", core_errors.ErrInvalidArgument),
			"failed to confirm email",
		)

		return
	}

	if err := h.authService.ConfirmEmail(ctx, token); err != nil {
		log.Warn("failed to confirm email", zap.Error(err))

		reason := "invalid"
		if errors.Is(err, core_errors.ErrInvalidArgument) {
			reason = "expired"
		}

		http.Redirect(rw, r, redirectBase+"?status=error&reason"+reason, http.StatusFound)
		return
	}

	http.Redirect(rw, r, redirectBase+"?status=success", http.StatusFound)
}
