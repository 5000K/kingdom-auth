package kingdomauth

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

// Client is a client for Kingdom Auth service, intended to be used by other services.
type Client struct {
	baseURL string
	secret  string

	providers []string

	log *slog.Logger
}

func NewClient(baseURL string, secret string) (*Client, error) {
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

	client := &Client{
		baseURL:   baseURL,
		secret:    secret,
		providers: make([]string, 0),
		log:       log,
	}

	err := client.loadProviders()

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
