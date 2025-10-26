package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/5000K/kingdom-auth/config"
	"github.com/5000K/kingdom-auth/db"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

type Service struct {
	config *config.Config
	log    *slog.Logger
	db     *db.Driver

	key []byte
}

func NewService(config *config.Config, db *db.Driver) (*Service, error) {

	return &Service{
		config: config,
		db:     db,

		key: []byte(config.KeyPhrase),

		log: slog.With("source", "auth-service"),
	}, nil
}

func (s *Service) getRedirectUrl(providerName string) string {
	return fmt.Sprintf("%s/oauth/end/%s", s.config.MainService.PublicUrl, providerName)
}

func (s *Service) createJwtFor(user *db.User) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"uid": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // todo: configurable
	})

	return t.SignedString(s.key)
}

func (s *Service) Run() {
	providers := make([]*Provider, 0)

	for _, p := range s.config.OAuthProviders {
		provider, err := NewProvider(&p, s.getRedirectUrl(p.Name))
		if err != nil {
			s.log.Error("Failed to create provider", "provider", p.Name, "error", err)
			return
		}

		providers = append(providers, provider)
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
				url := provider.config.AuthCodeURL("state", oauth2.AccessTypeOffline)
				c.Redirect(http.StatusFound, url)
				return
			}
		}

		c.Writer.WriteHeader(http.StatusNotFound)
	})

	r.GET("/oauth/end/:provider", func(c *gin.Context) {
		prov := c.Param("provider")

		for _, provider := range providers {
			if provider.Name == prov {
				code := c.Query("code")
				token, err := provider.config.Exchange(context.Background(), code)

				if err != nil {
					c.Writer.WriteHeader(http.StatusInternalServerError)
					s.log.Info("exchange error", "error", err)
					return
				}

				userInfo, err := provider.OICDProvider.UserInfo(context.Background(), provider.config.TokenSource(context.Background(), token))
				if err != nil {
					return
				}

				if err != nil {
					c.Writer.WriteHeader(http.StatusInternalServerError)
					s.log.Info("verifier error", "error", err)
					return
				}

				// try get auth
				auth, err := s.db.TryGetAuthentication(provider.Name, userInfo.Subject)

				if err != nil {
					// new user!
					usr, err := s.db.CreateUser()
					if err != nil {
						c.Writer.WriteHeader(http.StatusInternalServerError)
						s.log.Info("create user error", "error", err)
						return
					}

					auth, err = s.db.CreateAuthenticationFor(usr)
					if err != nil {
						c.Writer.WriteHeader(http.StatusInternalServerError)
						s.log.Info("create auth error", "error", err)
						return
					}

					auth.Provider = provider.Name
					auth.Subject = userInfo.Subject
					auth.Email = userInfo.Email
					err = s.db.UpdateAuthentication(auth)
					if err != nil {
						c.Writer.WriteHeader(http.StatusInternalServerError)
						s.log.Info("update auth error", "error", err)
						return
					}
				}

				user, err := s.db.GetUserFor(auth)

				if err != nil {
					c.Writer.WriteHeader(http.StatusInternalServerError)
					s.log.Info("get user error", "error", err)
					return
				}

				user.LastLogin = time.Now()
				_ = s.db.UpdateUser(user)

				j, err := s.createJwtFor(user)
				if err != nil {
					c.Writer.WriteHeader(http.StatusInternalServerError)
					s.log.Info("create jwt error", "error", err)
					return
				}

				c.SetCookie(s.config.CookieName, j, 3600*24, "/", s.config.CookieDomain, true, true)
			}
		}
	})

	r.Use(gin.Recovery())

	err := r.Run(fmt.Sprintf("0.0.0.0:%d", s.config.MainService.Port))

	if err != nil {
		s.log.Error("error running main service", "error", err)
	}
}
