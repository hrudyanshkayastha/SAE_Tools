package shuffle

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestShuffleClient_States(t *testing.T) {
	// 1. Success Webhook
	tsSuccess := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "execution_id": "exec-123"}`))
	}))
	defer tsSuccess.Close()

	client := NewClient(tsSuccess.URL)
	execID, err := client.ExecuteWorkflow(context.Background(), ActionPayload{CorrelationID: "1", Action: "test", Target: "1.2.3.4"})
	if err != nil || execID != "exec-123" {
		t.Fatalf("Expected success, got %v, execID: %s", err, execID)
	}

	// 2. Failed Execution (HTTP 500)
	tsFail := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer tsFail.Close()

	clientFail := NewClient(tsFail.URL)
	_, err = clientFail.ExecuteWorkflow(context.Background(), ActionPayload{CorrelationID: "2"})
	if err == nil {
		t.Fatal("Expected failure, got success")
	}

	// Polling verification test: Success
	apiSuccess := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "SUCCESS", "result": "Done"}`))
	}))
	defer apiSuccess.Close()
	
	client.APIURL = apiSuccess.URL
	res, _ := client.PollExecutionStatus(context.Background(), "exec-123", 10*time.Millisecond)
	if res.Status != "SUCCEEDED" {
		t.Fatalf("Expected SUCCEEDED, got %s", res.Status)
	}

	// Polling verification test: Timeout
	apiTimeout := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "EXECUTING", "result": "Running..."}`))
	}))
	defer apiTimeout.Close()
	
	client.APIURL = apiTimeout.URL
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	res2, err2 := client.PollExecutionStatus(ctx, "exec-456", 10*time.Millisecond)
	if err2 == nil || res2.Status != "TIMEOUT" {
		t.Fatalf("Expected TIMEOUT, got %s with err %v", res2.Status, err2)
	}
}
