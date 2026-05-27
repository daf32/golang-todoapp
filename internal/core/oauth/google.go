package core_oauth

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

const googleIssuer = "https://accounts.google.com"

type GoogleProvider struct {
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
}

func NewGoogleProvider(
	ctx context.Context,
	config GoogleConfig,
	redirectURL string,
) (*GoogleProvider, error) {
	provider, err := oidc.NewProvider(ctx, googleIssuer)
	if err != nil {
		return nil, fmt.Errorf("init google oidc provider: %w", err)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  redirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: config.ClientID})

	return &GoogleProvider{
		oauth2Config: oauth2Config,
		verifier:     verifier,
	}, nil
}

func (p *GoogleProvider) AuthCodeURL(state, codeVerifier string) string {
	return p.oauth2Config.AuthCodeURL(
		state,
		oauth2.AccessTypeOnline,
		oauth2.S256ChallengeOption(codeVerifier),
	)
}

func (p *GoogleProvider) Exchange(
	ctx context.Context,
	code, codeVerifier string,
) (UserInfo, error) {
	token, err := p.oauth2Config.Exchange(
		ctx,
		code,
		oauth2.VerifierOption(codeVerifier),
	)
	if err != nil {
		return UserInfo{}, fmt.Errorf("exchage code: %w", err)
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return UserInfo{}, fmt.Errorf("id_token missing from response")
	}

	idToken, err := p.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return UserInfo{}, fmt.Errorf("verify id_token: %w", err)
	}

	var claims struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
	}

	if err := idToken.Claims(&claims); err != nil {
		return UserInfo{}, fmt.Errorf("parse id_token claims: %w", err)
	}

	if claims.Sub == "" || claims.Email == "" {
		return UserInfo{}, fmt.Errorf("id_token missing sub or email")
	}

	return UserInfo{
		Sub:           claims.Sub,
		Email:         claims.Email,
		EmailVerified: claims.EmailVerified,
		Name:          claims.Name,
	}, nil
}

func (p *GoogleProvider) Name() string { return ProviderGoogle }
