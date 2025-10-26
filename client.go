package kingdomauth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/5000K/kingdom-auth/core"
	"github.com/golang-jwt/jwt/v5"
)

// Client is a client for Kingdom Auth service, intended to be used by other services.
type Client struct {
	baseURL string
	secret  string

	providers []string
	publicKey *rsa.PublicKey

	log *slog.Logger
}

func NewClient(baseURL string, secret string, publicKeyPath string) (*Client, error) {
	log := slog.With("source", "kingdomauth.Client")

	if !strings.HasPrefix(baseURL, "https://") {

		if !strings.HasPrefix(baseURL, "http://") {
			return nil, fmt.Errorf("kingdomauth baseURL must start with http:// or https://")
		}

		log.Warn("baseURL does not use https - this is supported but not recommended")
	}

	if strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL[:len(baseURL)-1]
	}

	// Load public key
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	publicKeyBlock, _ := pem.Decode(publicKeyData)
	if publicKeyBlock == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	client := &Client{
		baseURL:   baseURL,
		secret:    secret,
		providers: make([]string, 0),
		publicKey: publicKey,
		log:       log,
	}

	err = client.loadProviders()

	if err != nil {
		return nil, err
	}

	return client, nil
}

type providersAnswer struct {
	Providers []string `json:"providers"`
}

func (c *Client) loadProviders() error {
	url := c.baseURL + "/providers"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.secret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// read answer
	resBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	var answer providersAnswer
	err = json.Unmarshal(resBody, &answer)
	if err != nil {
		return err
	}

	c.providers = answer.Providers

	return nil
}

// ValidateToken validates a JWT token using the public key and returns the claims.
// It verifies the signature and checks if the token is expired.
// Returns jwt.MapClaims on success, or an error if validation fails.
func (c *Client) ValidateToken(token string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	tkn, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return c.publicKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, core.ErrInvalidSignature
		}

		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, core.ErrTokenExpired
		}

		c.log.Debug("failed to parse token", "err", err)
		return nil, core.ErrFailedToParseToken
	}

	if !tkn.Valid {
		return nil, core.ErrTokenInvalid
	}

	return claims, nil
}
