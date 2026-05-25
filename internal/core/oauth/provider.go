package core_oauth

import "context"

const (
	ProviderGoogle = "google"
	ProviderApple  = "apple"
)

type UserInfo struct {
	Sub           string
	Email         string
	EmailVerified bool
	Name          string
}

type Provider interface {
	Name() string
	AuthCodeURL(state, codeVerifier string) string
	Exchange(ctx context.Context, code, codeVerifier string) (UserInfo, error)
}

