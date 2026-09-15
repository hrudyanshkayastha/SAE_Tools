package cortex

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

type JobPayload struct {
	Data    string `json:"data"`
	DataType string `json:"dataType"`
	Tlp     int    `json:"tlp"`
}

type JobResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func (c *Client) SubmitJob(ctx context.Context, analyzerName string, payload JobPayload) (string, error) {
	if c.URL == "" {
		return "", fmt.Errorf("Cortex URL not configured")
	}

	body, _ := json.Marshal(payload)
	endpoint := fmt.Sprintf("%s/api/analyzer/%s/run", c.URL, analyzerName)
	
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Cortex API failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("Cortex API returned HTTP %d", resp.StatusCode)
	}

	var res JobResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	return res.ID, nil
}
