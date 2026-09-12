package shuffle

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestShuffleClient_States(t *testing.T) {
	// 1. Success
	tsSuccess := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer tsSuccess.Close()

	client := NewClient(tsSuccess.URL)
	_, err := client.ExecuteWorkflow(context.Background(), ActionPayload{CorrelationID: "1", Action: "test", Target: "1.2.3.4"})
	if err != nil {
		t.Fatalf("Expected success, got %v", err)
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

	// 3. Timeout
	tsTimeout := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer tsTimeout.Close()

	clientTimeout := NewClient(tsTimeout.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	_, err = clientTimeout.ExecuteWorkflow(ctx, ActionPayload{CorrelationID: "3"})
	if err == nil {
		t.Fatal("Expected timeout error, got success")
	}

	// 4. Unavailable
	clientUnavailable := NewClient("http://localhost:99999") // invalid port
	_, err = clientUnavailable.ExecuteWorkflow(context.Background(), ActionPayload{CorrelationID: "4"})
	if err == nil {
		t.Fatal("Expected unavailable error, got success")
	}

	// Verification check tests
	res, _ := client.CheckStatus(context.Background(), "exec-123")
	if res.Status != "SUCCEEDED" {
		t.Fatalf("Expected SUCCEEDED, got %s", res.Status)
	}
}
