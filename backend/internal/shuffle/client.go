package shuffle

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	WebhookURL string
	HTTPClient *http.Client
}

func NewClient(webhookURL string) *Client {
	return &Client{
		WebhookURL: webhookURL,
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
	}
}

type ActionPayload struct {
	CorrelationID string `json:"correlation_id"`
	Action        string `json:"action"`
	Target        string `json:"target"`
}

// ExecuteWorkflow triggers a Shuffle workflow.
func (c *Client) ExecuteWorkflow(ctx context.Context, payload ActionPayload) error {
	if c.WebhookURL == "" {
		return fmt.Errorf("shuffle webhook URL not configured")
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", c.WebhookURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("shuffle webhook failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("shuffle webhook returned HTTP %d", resp.StatusCode)
	}

	return nil
}
