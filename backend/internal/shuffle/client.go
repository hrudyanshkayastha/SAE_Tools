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
	APIURL     string
	AuthToken  string
	HTTPClient *http.Client
}

func NewClient(webhookURL string) *Client {
	return &Client{
		WebhookURL: webhookURL,
		APIURL:     "http://localhost:5001/api/v1", // Default local Shuffle API
		AuthToken:  "mock-auth-token",
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
	}
}

type ActionPayload struct {
	CorrelationID string `json:"correlation_id"`
	Action        string `json:"action"`
	Target        string `json:"target"`
}

type ExecutionResult struct {
	ExecutionID string
	Status      string
	Message     string
}

type WebhookResponse struct {
	ExecutionID string `json:"execution_id"`
	Success     bool   `json:"success"`
}

type APIStatusResponse struct {
	Status string `json:"status"`
	Result string `json:"result"`
}

func (c *Client) ExecuteWorkflow(ctx context.Context, payload ActionPayload) (string, error) {
	if c.WebhookURL == "" {
		return "", fmt.Errorf("shuffle webhook URL not configured")
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", c.WebhookURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("shuffle webhook failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("shuffle webhook returned HTTP %d", resp.StatusCode)
	}

	var whResp WebhookResponse
	if err := json.NewDecoder(resp.Body).Decode(&whResp); err != nil {
		return "", fmt.Errorf("failed to parse shuffle webhook response: %w", err)
	}

	if whResp.ExecutionID == "" {
		return "", fmt.Errorf("shuffle webhook did not return an execution_id")
	}

	return whResp.ExecutionID, nil
}

func (c *Client) CheckStatus(ctx context.Context, executionID string) (ExecutionResult, error) {
	url := fmt.Sprintf("%s/executions/%s", c.APIURL, executionID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return ExecutionResult{Status: "VERIFICATION_FAILED"}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.AuthToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return ExecutionResult{Status: "VERIFICATION_FAILED"}, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return ExecutionResult{Status: "VERIFICATION_FAILED"}, fmt.Errorf("API returned HTTP %d", resp.StatusCode)
	}

	var apiResp APIStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return ExecutionResult{Status: "VERIFICATION_FAILED"}, fmt.Errorf("failed to parse API response: %w", err)
	}

	var mappedStatus string
	switch apiResp.Status {
	case "SUCCESS":
		mappedStatus = "SUCCEEDED"
	case "FAILURE":
		mappedStatus = "FAILED"
	case "EXECUTING":
		mappedStatus = "EXECUTING"
	default:
		mappedStatus = "VERIFICATION_FAILED"
	}

	return ExecutionResult{
		ExecutionID: executionID,
		Status:      mappedStatus,
		Message:     apiResp.Result,
	}, nil
}

func (c *Client) PollExecutionStatus(ctx context.Context, executionID string, interval time.Duration) (ExecutionResult, error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ExecutionResult{ExecutionID: executionID, Status: "TIMEOUT", Message: "Polling context timed out"}, ctx.Err()
		case <-ticker.C:
			res, err := c.CheckStatus(ctx, executionID)
			if err != nil {
				if ctx.Err() != nil {
					return ExecutionResult{ExecutionID: executionID, Status: "TIMEOUT", Message: "Polling context timed out"}, ctx.Err()
				}
				return res, err
			}
			if res.Status == "SUCCEEDED" || res.Status == "FAILED" || res.Status == "VERIFICATION_FAILED" {
				return res, nil
			}
			// If EXECUTING, loop continues
		}
	}
}
