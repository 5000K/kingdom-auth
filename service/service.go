package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/5000K/kingdom-auth/config"
	"github.com/5000K/kingdom-auth/core"
	"github.com/5000K/kingdom-auth/db"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

var closeWindowPage = []byte("<html><script>window.close();</script></html>")

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

		key: []byte(config.Token.KeyPhrase),

		log: slog.With("source", "auth-service"),
	}, nil
}

func (s *Service) getRedirectUrl(providerName string) string {
	return fmt.Sprintf("%s/oauth/end/%s", s.config.MainService.PublicUrl, providerName)
}

func (s *Service) createRefreshTokenFor(user *db.User) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub":                        fmt.Sprintf("%d", user.ID),
		"iss":                        s.config.Token.Issuer,
		"exp":                        time.Now().Add(time.Second * time.Duration(s.config.Token.RefreshTokenTTL)).Unix(),
		"iat":                        time.Now().Unix(),
		core.KingdomAuthVersionClaim: core.KingdomAuthVersion,
	})

	return t.SignedString(s.key)
}

func (s *Service) readRefreshToken(token string) (jwt.MapClaims, error) {
	contents := jwt.MapClaims{}
	tkn, err := jwt.ParseWithClaims(token, &contents, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, core.ErrInvalidSignature
		}

		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, core.ErrTokenExpired
		}

		s.log.Debug("failed to parse token", "err", err)
		return nil, core.ErrFailedToParseToken
	}

	if !tkn.Valid {
		return nil, core.ErrTokenInvalid
	}

	return contents, nil
}

func (s *Service) createAuthTokenFor(user *db.User) (string, int64, error) {
	aud := s.config.Token.DefaultAudience

	pud, err := user.GetPrivateUserdata()

	if err == nil {
		potentialAud, ok := pud["aud"].(string)
		if ok {
			aud = potentialAud
		}
	}

	exp := time.Now().Add(time.Second * time.Duration(s.config.Token.AuthTokenTTL)).Unix()

	t := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub":                        fmt.Sprintf("%d", user.ID),
		"aud":                        aud,
		"iss":                        s.config.Token.Issuer,
		"exp":                        exp,
		"iat":                        time.Now().Unix(),
		"public-data":                user.PublicData,
		core.KingdomAuthVersionClaim: core.KingdomAuthVersion,
	})

	tk, err := t.SignedString(s.key)

	return tk, exp, err
}

func (s *Service) readAuthToken(token string) (jwt.MapClaims, error) {
	contents := jwt.MapClaims{}
	tkn, err := jwt.ParseWithClaims(token, &contents, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, core.ErrInvalidSignature
		}

		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, core.ErrTokenExpired
		}

		s.log.Debug("failed to parse token", "err", err)
		return nil, core.ErrFailedToParseToken
	}

	if !tkn.Valid {
		return nil, core.ErrTokenInvalid
	}

	return contents, nil
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

		c.JSON(http.StatusOK, gin.H{
			"providers": list,
		})
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

		c.JSON(http.StatusNotFound, gin.H{
			"error": "provider not found",
		})
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

				j, err := s.createRefreshTokenFor(user)
				if err != nil {
					c.Writer.WriteHeader(http.StatusInternalServerError)
					s.log.Info("create jwt error", "error", err)
					return
				}

				c.SetCookie(s.config.CookieName, j, 3600*24, "/", s.config.CookieDomain, true, true)

				// write close window page

				c.Writer.WriteHeader(http.StatusOK)
				c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
				c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				c.Writer.Header().Set("Pragma", "no-cache")
				c.Writer.Header().Set("Expires", "0")

				_, _ = c.Writer.Write(closeWindowPage)
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{
			"error": "provider not found",
		})
	})

	r.GET("/token", func(c *gin.Context) {
		cookieString, err := c.Cookie(s.config.CookieName)

		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
		}

		tk, err := s.readRefreshToken(cookieString)

		if err != nil {
			if errors.Is(err, core.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "token expired",
				})
				return
			} else if errors.Is(err, core.ErrTokenInvalid) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "token invalid",
				})
				return
			} else if errors.Is(err, core.ErrInvalidSignature) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "token signature invalid",
				})
				return
			}
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		version := tk[core.KingdomAuthVersionClaim]

		if version != core.KingdomAuthVersion {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":    "version mismatch: token is from another format (older or newer)",
				"expected": core.KingdomAuthVersion,
				"actual":   version,
			})
			return
		}

		iss, err := tk.GetIssuer()
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			s.log.Info("get issuer error", "error", err)
			return
		}

		if iss != s.config.Token.Issuer {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":    "issuer mismatch",
				"expected": s.config.Token.Issuer,
				"actual":   iss,
			})
			return
		}

		uidS, err := tk.GetSubject()

		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			s.log.Info("get subject error", "error", err)
			return
		}

		uid, err := strconv.ParseUint(uidS, 10, 32)

		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		// find user
		user, err := s.db.GetUser(uint32(uid))

		if err != nil {
			c.Writer.WriteHeader(http.StatusNotFound)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token valid, but user not found",
			})
			return
		}

		// check for expiry, send a new refresh token, if token is old enough
		expiry, err := tk.GetExpirationTime()

		if err == nil {
			issueDate, err := tk.GetIssuedAt()
			if err == nil {
				timeSinceIssue := uint(issueDate.Time.Sub(expiry.Time).Seconds())

				if timeSinceIssue > s.config.Token.MinAgeForRefresh {
					j, err := s.createRefreshTokenFor(user)
					if err != nil {
						c.Writer.WriteHeader(http.StatusInternalServerError)
						s.log.Info("create jwt error", "error", err)
						return
					}

					c.SetCookie(s.config.CookieName, j, 3600*24, "/", s.config.CookieDomain, true, true)
				}
			}
		}

		at, exp, err := s.createAuthTokenFor(user)

		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			s.log.Info("create jwt error", "error", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": at,
			"exp":   exp,
		})
	})

	r.Use(gin.Recovery())

	err := r.Run(fmt.Sprintf("0.0.0.0:%d", s.config.MainService.Port))

	if err != nil {
		s.log.Error("error running main service", "error", err)
	}
}
