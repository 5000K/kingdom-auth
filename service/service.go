package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/5000K/kingdom-auth/config"
	"github.com/5000K/kingdom-auth/db"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type Service struct {
	config *config.Config
	log    *slog.Logger
	db     *db.Driver

	states map[string]string
}

func NewService(config *config.Config, db *db.Driver) *Service {
	return &Service{
		config: config,
		db:     db,
		log:    slog.With("source", "auth-service"),
	}
}

func (s *Service) getRedirectUrl(providerName string) string {
	return fmt.Sprintf("%s/oauth/end/%s", s.config.MainService.PublicUrl, providerName)
}

func (s *Service) Run() {
	providers := make([]*Provider, 0)

	for _, p := range s.config.OAuthProviders {
		providers = append(providers, NewProvider(&p, s.getRedirectUrl(p.Name)))
	}

	if len(providers) == 0 {
		s.log.Error("no providers loaded - can't start")
		os.Exit(1)
		return
	}

	r := gin.New()

	r.GET("/providers", func(c *gin.Context) {
		list := make([]string, 0)

		for _, p := range providers {
			list = append(list, p.Name)
		}

		c.JSON(http.StatusOK, list)
	})

	r.GET("/oauth/begin/:provider", func(c *gin.Context) {
		prov := c.Param("provider")

		for _, provider := range providers {
			if provider.Name == prov {
				// redirect
				url := provider.Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
				c.Redirect(http.StatusFound, url)
				return
			}
		}

		c.Writer.WriteHeader(http.StatusNotFound)
	})

	r.GET("/oauth/callback/:provider", func(c *gin.Context) {
		prov := c.Param("provider")

		for _, provider := range providers {
			if provider.Name == prov {
				code := c.Query("code")
				url, err := provider.Config.Exchange(context.Background(), code)

				if err != nil {
					c.Writer.WriteHeader(http.StatusInternalServerError)
					s.log.Info("exchange error", "error", err)
					return
				}

			}
		}
	})

	r.Use(gin.Recovery())

	err := r.Run(fmt.Sprintf("0.0.0.0:%d", conf.MainService.Port))

	if err != nil {
		s.log.Error("error running main service", "error", err)
	}
}
