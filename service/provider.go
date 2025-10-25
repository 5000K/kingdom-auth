package service

import (
	"github.com/5000K/kingdom-auth/config"
	"golang.org/x/oauth2"
)

type Provider struct {
	Name   string
	Config oauth2.Config
}

func NewProvider(config *config.OAuthConfig, redirectUrl string) *Provider {
	ep := oauth2.Endpoint{
		AuthURL:  config.AuthUrl,
		TokenURL: config.TokenUrl,
	}

	return &Provider{
		Name: config.Name,
		Config: oauth2.Config{
			ClientID:     config.ClientId,
			ClientSecret: config.ClientSecret,
			Endpoint:     ep,
			RedirectURL:  redirectUrl,
		},
	}
}
