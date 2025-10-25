package service

import (
	"context"

	"github.com/5000K/kingdom-auth/config"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Provider struct {
	Name         string
	config       oauth2.Config
	OICDProvider *oidc.Provider
}

func NewProvider(config *config.OAuthConfig, redirectUrl string) (*Provider, error) {
	oProv, err := oidc.NewProvider(context.Background(), config.Url)

	if err != nil {
		return nil, err
	}

	scopes := config.Scopes

	hasOICDScope := false

	for _, scope := range config.Scopes {
		if scope == oidc.ScopeOpenID {
			hasOICDScope = true
		}
	}

	if !hasOICDScope {
		scopes = append(scopes, oidc.ScopeOpenID)
	}

	return &Provider{
		Name: config.Name,
		config: oauth2.Config{
			Scopes:       config.Scopes,
			ClientID:     config.ClientId,
			ClientSecret: config.ClientSecret,
			Endpoint:     oProv.Endpoint(),
			RedirectURL:  redirectUrl,
		},
		OICDProvider: oProv,
	}, nil
}

func (p *Provider) getVerifier() *oidc.IDTokenVerifier {
	return p.OICDProvider.Verifier(&oidc.Config{ClientID: p.config.ClientID})
}
