package wazuh

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is the SAE adapter to communicate with the isolated Wazuh Manager API.
type Client struct {
	BaseURL    string
	Username   string
	Password   string
	HTTPClient *http.Client
	Token      string
}

// NewClient creates a new Wazuh API client.
func NewClient(baseURL, username, password string, insecureSkipVerify bool) *Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
	}
	return &Client{
		BaseURL:  baseURL,
		Username: username,
		Password: password,
		HTTPClient: &http.Client{
			Timeout:   10 * time.Second, // Timeout to prevent hangs
			Transport: transport,
		},
	}
}

// Authenticate fetches a JWT token from the Wazuh API.
func (c *Client) Authenticate(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/security/user/authenticate", c.BaseURL), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("wazuh authentication failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("wazuh auth error (status %d): %s", resp.StatusCode, string(body))
	}

	var authResp struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("failed to decode wazuh token: %w", err)
	}

	c.Token = authResp.Data.Token
	return nil
}

// HealthCheck verifies the Wazuh Manager is online and accessible.
func (c *Client) HealthCheck(ctx context.Context) error {
	if c.Token == "" {
		if err := c.Authenticate(ctx); err != nil {
			return err
		}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/manager/status", c.BaseURL), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("wazuh health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		// Token might have expired, clear it so next call re-authenticates
		c.Token = ""
		return fmt.Errorf("unauthorized (token expired)")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("wazuh manager not healthy, status: %d", resp.StatusCode)
	}

	return nil
}

// GetAgents returns the list of connected agents
func (c *Client) GetAgents(ctx context.Context) (*APIResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/agents", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch agents: %w", err)
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode agents response: %w", err)
	}

	return &apiResp, nil
}
