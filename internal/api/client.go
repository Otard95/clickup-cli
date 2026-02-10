package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/otard95/clickup-cli/internal/config"
)

const baseURL = "https://api.clickup.com/api/v2"

type Client struct {
	cfg  *config.Config
	http *http.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg: cfg,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) TeamID() string {
	return c.cfg.TeamID
}

// request performs an HTTP request and decodes the JSON response into dest.
// params is a map of query parameters; values that are slices will be expanded
// into repeated keys (e.g. "assignees[]" => ["1","2"]).
func (c *Client) request(method, endpoint string, body io.Reader, params map[string]string, dest interface{}) error {
	u, err := url.Parse(baseURL + endpoint)
	if err != nil {
		return fmt.Errorf("invalid endpoint: %w", err)
	}

	if params != nil {
		q := u.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", c.cfg.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	if dest != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, dest); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}

	return nil
}

func (c *Client) Get(endpoint string, params map[string]string, dest interface{}) error {
	return c.request(http.MethodGet, endpoint, nil, params, dest)
}

func (c *Client) Put(endpoint string, body io.Reader, params map[string]string, dest interface{}) error {
	return c.request(http.MethodPut, endpoint, body, params, dest)
}

func (c *Client) Post(endpoint string, body io.Reader, params map[string]string, dest interface{}) error {
	return c.request(http.MethodPost, endpoint, body, params, dest)
}

// SetQueryArray adds repeated query params to a URL (e.g. assignees[]=1&assignees[]=2).
// This is a helper for building params maps for endpoints that need array params.
func SetQueryArray(endpoint string, key string, values []string) string {
	if len(values) == 0 {
		return endpoint
	}
	sep := "?"
	if strings.Contains(endpoint, "?") {
		sep = "&"
	}
	var parts []string
	for _, v := range values {
		parts = append(parts, url.QueryEscape(key)+"="+url.QueryEscape(v))
	}
	return endpoint + sep + strings.Join(parts, "&")
}
