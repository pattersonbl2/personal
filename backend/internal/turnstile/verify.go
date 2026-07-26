package turnstile

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultVerifyURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"

// Client verifies Cloudflare Turnstile tokens via the siteverify API.
type Client struct {
	Secret     string
	HTTPClient *http.Client
	VerifyURL  string
}

// Result is the subset of the siteverify response we care about.
type Result struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
	Hostname   string   `json:"hostname"`
}

// Verify checks a Turnstile response token. remoteIP is optional.
func (c *Client) Verify(ctx context.Context, token, remoteIP string) (*Result, error) {
	if c == nil || c.Secret == "" {
		return nil, fmt.Errorf("turnstile: secret not configured")
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return &Result{Success: false, ErrorCodes: []string{"missing-input-response"}}, nil
	}

	form := url.Values{}
	form.Set("secret", c.Secret)
	form.Set("response", token)
	if remoteIP != "" {
		form.Set("remoteip", remoteIP)
	}

	endpoint := c.VerifyURL
	if endpoint == "" {
		endpoint = defaultVerifyURL
	}
	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("turnstile: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("turnstile: siteverify request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	if err != nil {
		return nil, fmt.Errorf("turnstile: read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("turnstile: siteverify status %d", resp.StatusCode)
	}

	var result Result
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("turnstile: decode response: %w", err)
	}
	return &result, nil
}
