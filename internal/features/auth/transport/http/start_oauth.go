package auth_transport_http

import (
	"net/http"
	"time"

	core_auth "github.com/daf32/golang-todoapp/internal/core/auth"
	core_logger "github.com/daf32/golang-todoapp/internal/core/logger"
	core_http_response "github.com/daf32/golang-todoapp/internal/core/transport/http/response"
)

const (
	oauthStateCookie    = "oauth_state"
	oauthVerifierCookie = "oauth_verifier"
	oauthCookieMaxAge   = 10 * time.Minute
)

// StartOAuth godoc
// @Summary      Start OAuth flow
// @Description  Generates state and PKCE verifier, stores them in HttpOnly cookies and redirects the user to the provider's consent screen
// @Tags         auth
// @Param        provider path  string true "OAuth provider name (e.g. google)"
// @Success      302 "Redirect to the provider's OAuth consent screen"
// @Failure      400 {object} core_http_response.ErrorResponse "Unknown OAuth provider"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /auth/oauth/{provider} [get]
func (h *AuthHTTPHandler) StartOAuth(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	state, err := core_auth.GenerateSecureToken(32)
	if err != nil {
		responseHandler.ErrorResponse(err, "generate oauth state")
		return
	}

	verifier, err := core_auth.GenerateSecureToken(32)
	if err != nil {
		responseHandler.ErrorResponse(err, "generate oauth verifier")
		return
	}

	providerName := r.PathValue("provider")

	url, err := h.authService.BuildOAuthURL(
		providerName,
		state,
		verifier,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"build OAuth URL",
		)

		return
	}

	h.cookieManager.Set(rw, oauthStateCookie, state, oauthCookieMaxAge)
	h.cookieManager.Set(rw, oauthVerifierCookie, verifier, oauthCookieMaxAge)

	http.Redirect(rw, r, url, http.StatusFound)
}
