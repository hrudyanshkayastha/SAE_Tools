package thehive

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	URL        string
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(url, apiKey string) *Client {
	return &Client{
		URL:        url,
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
	}
}

type CasePayload struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Severity    int      `json:"severity"`
	Tags        []string `json:"tags"`
}

type CaseResponse struct {
	ID string `json:"id"`
	// other fields omitted
}

func (c *Client) CreateCase(ctx context.Context, payload CasePayload) (string, error) {
	if c.URL == "" {
		return "", fmt.Errorf("TheHive URL not configured")
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", c.URL+"/api/case", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("TheHive API failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("TheHive API returned HTTP %d", resp.StatusCode)
	}

	var res CaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	return res.ID, nil
}
