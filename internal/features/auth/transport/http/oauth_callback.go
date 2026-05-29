package auth_transport_http

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"net/url"

	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	core_logger "github.com/daf32/golang-todoapp/internal/core/logger"
	core_http_response "github.com/daf32/golang-todoapp/internal/core/transport/http/response"
)

const oauthSuccessRedirectPath = "/oauth-callback"

type OAuthCallbackResponse struct {
	AccessToken  string `json:"access_token"  example:"access_token"`
	RefreshToken string `json:"refresh_token" example:"refresh_token"`
}

// OAuthCallback godoc
// @Summary      OAuth callback
// @Description  Handles the provider's redirect after consent. Verifies state, exchanges the code via PKCE, finds or creates the user and issues tokens.
// @Tags         auth
// @Produce      json
// @Param        provider path  string true "OAuth provider name (e.g. google)"
// @Param        code  query string true "Authorization code returned by the provider"
// @Param        state query string true "State value returned by the provider, must match cookie"
// @Success      200 {object} OAuthCallbackResponse                      "Successful login via OAuth"
// @Failure      400 {object} core_http_response.ErrorResponse           "Missing/invalid state, verifier, code or unknown provider"
// @Failure      500 {object} core_http_response.ErrorResponse           "Internal server error"
// @Router       /auth/oauth/{provider}/callback [get]
func (h *AuthHTTPHandler) OAuthCallback(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	defer h.cookieManager.Clear(rw, oauthStateCookie)
	defer h.cookieManager.Clear(rw, oauthVerifierCookie)

	stateCookie, err := r.Cookie(oauthStateCookie)
	if err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("missing oauth state cookie: %w", core_errors.ErrInvalidArgument),
			"oauth callback",
		)
		return
	}

	verifierCookie, err := r.Cookie(oauthVerifierCookie)
	if err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("missing oauth verifier cookie: %w", core_errors.ErrInvalidArgument),
			"oauth callback",
		)
		return
	}

	queryState := r.URL.Query().Get("state")
	if queryState == "" || subtle.ConstantTimeCompare([]byte(queryState), []byte(stateCookie.Value)) != 1 {
		responseHandler.ErrorResponse(
			fmt.Errorf("state mismatch: %w", core_errors.ErrInvalidArgument),
			"oauth callback",
		)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		responseHandler.ErrorResponse(
			fmt.Errorf("missing code: %w", core_errors.ErrInvalidArgument),
			"oauth callback",
		)
		return
	}

	providerName := r.PathValue("provider")

	accessToken, refreshToken, err := h.authService.LoginUserWithOAuth(
		ctx,
		providerName,
		code,
		verifierCookie.Value,
	)
	if err != nil {
		responseHandler.ErrorResponse(err, "login with oauth")
		return
	}

	// Redirect to the frontend with tokens in the URL fragment.
	// Fragments aren't sent to the server on subsequent requests and don't
	// appear in server access logs, so they're a safer transport than query
	// strings for short-lived secrets handed off to the SPA.
	fragment := url.Values{}
	fragment.Set("access_token", accessToken)
	fragment.Set("refresh_token", refreshToken.Token)
	http.Redirect(rw, r, h.appBaseURL+oauthSuccessRedirectPath+"#"+fragment.Encode(), http.StatusFound)
}
