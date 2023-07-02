package authenticator

import (
	"context"
	"errors"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	env "invisibleprogrammer.com/invisibleurl/environment"
)

const AUTH0_DOMAIN = "AUTH0_DOMAIN"
const AUTH0_CLIENT_ID = "AUTH0_CLIENT_ID"
const AUTH0_CLIENT_SECRET = "AUTH0_CLIENT_SECRET"
const AUTH0_CALLBACK_URL = "AUTH0_CALLBACK_URL"

type Authenticator struct {
	*oidc.Provider
	oauth2.Config
}

func New() (*Authenticator, error) {
	provider, err := oidc.NewProvider(
		context.Background(),
		"https://"+os.Getenv(env.AUTH0_DOMAIN)+"/",
	)
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     os.Getenv(env.AUTH0_CLIENT_ID),
		ClientSecret: os.Getenv(env.AUTH0_CLIENT_SECRET),
		RedirectURL:  os.Getenv(env.AUTH0_CALLBACK_URL),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
	}, nil
}

func (a *Authenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDTOken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDTOken)
}
