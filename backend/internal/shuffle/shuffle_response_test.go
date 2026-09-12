package shuffle

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestShuffleClient_MissingConfig(t *testing.T) {
	os.Setenv("SAE_SHUFFLE_API_URL", "")
	os.Setenv("SAE_SHUFFLE_AUTH_TOKEN", "")
	client := NewClient("http://mock-webhook")

	res, err := client.PollExecutionStatus(context.Background(), "exec-123", 10*time.Millisecond)
	if err == nil || res.Status != "VERIFICATION_FAILED" {
		t.Fatalf("Expected missing config failure, got status %s", res.Status)
	}
}

func TestShuffleClient_WebhookResponses(t *testing.T) {
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

	// 2. Malformed Webhook (missing execution_id)
	tsMalformed := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer tsMalformed.Close()

	clientMalformed := NewClient(tsMalformed.URL)
	_, err = clientMalformed.ExecuteWorkflow(context.Background(), ActionPayload{CorrelationID: "2"})
	if err == nil {
		t.Fatal("Expected malformed webhook failure (no execution_id)")
	}

	// 3. Failed Execution (HTTP 500)
	tsFail := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer tsFail.Close()

	clientFail := NewClient(tsFail.URL)
	_, err = clientFail.ExecuteWorkflow(context.Background(), ActionPayload{CorrelationID: "2"})
	if err == nil {
		t.Fatal("Expected HTTP failure, got success")
	}
}

func TestShuffleClient_PollingStates(t *testing.T) {
	os.Setenv("SAE_SHUFFLE_API_URL", "mock")
	os.Setenv("SAE_SHUFFLE_AUTH_TOKEN", "mock")

	// 1. SUCCESS
	apiSuccess := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "SUCCESS", "result": "Done"}`))
	}))
	defer apiSuccess.Close()
	
	client := NewClient("mock")
	client.APIURL = apiSuccess.URL
	client.AuthToken = "mock"
	res, _ := client.PollExecutionStatus(context.Background(), "exec-123", 10*time.Millisecond)
	if res.Status != "SUCCEEDED" {
		t.Fatalf("Expected SUCCEEDED, got %s", res.Status)
	}

	// 2. FAILURE
	apiFailure := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "FAILURE", "result": "Error occurred"}`))
	}))
	defer apiFailure.Close()
	client.APIURL = apiFailure.URL
	res, _ = client.PollExecutionStatus(context.Background(), "exec-123", 10*time.Millisecond)
	if res.Status != "FAILED" {
		t.Fatalf("Expected FAILED, got %s", res.Status)
	}

	// 3. Unknown State
	apiUnknown := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "UNKNOWN_GARBAGE", "result": "???"}`))
	}))
	defer apiUnknown.Close()
	client.APIURL = apiUnknown.URL
	res, _ = client.PollExecutionStatus(context.Background(), "exec-123", 10*time.Millisecond)
	if res.Status != "VERIFICATION_FAILED" {
		t.Fatalf("Expected VERIFICATION_FAILED for unknown state, got %s", res.Status)
	}

	// 4. HTTP/API failure (500)
	apiError := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer apiError.Close()
	client.APIURL = apiError.URL
	res, err := client.PollExecutionStatus(context.Background(), "exec-123", 10*time.Millisecond)
	if res.Status != "VERIFICATION_FAILED" || err == nil {
		t.Fatalf("Expected VERIFICATION_FAILED with error")
	}

	// 5. EXECUTING -> SUCCESS (Dynamic)
	calls := 0
	apiDynamic := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusOK)
		if calls < 3 {
			w.Write([]byte(`{"status": "EXECUTING"}`))
		} else {
			w.Write([]byte(`{"status": "SUCCESS"}`))
		}
	}))
	defer apiDynamic.Close()
	client.APIURL = apiDynamic.URL
	res, _ = client.PollExecutionStatus(context.Background(), "exec-123", 10*time.Millisecond)
	if res.Status != "SUCCEEDED" || calls != 3 {
		t.Fatalf("Expected SUCCEEDED after 3 calls, got %s in %d calls", res.Status, calls)
	}

	// 6. Timeout
	apiTimeout := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "EXECUTING"}`))
	}))
	defer apiTimeout.Close()
	client.APIURL = apiTimeout.URL
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	res2, err2 := client.PollExecutionStatus(ctx, "exec-456", 15*time.Millisecond)
	if err2 == nil || res2.Status != "TIMEOUT" {
		t.Fatalf("Expected TIMEOUT, got %s with err %v", res2.Status, err2)
	}
}
