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
	os.Setenv("SAE_SHUFFLE_WORKFLOW_ID", "")
	client := NewClient("http://mock-webhook")

	res, err := client.PollExecutionStatus(context.Background(), "exec-123", 10*time.Millisecond)
	if err == nil || res.Status != "VERIFICATION_FAILED" {
		t.Fatalf("Expected missing config failure, got status %s", res.Status)
	}
}

func TestShuffleClient_WebhookResponses(t *testing.T) {
	tsSuccess := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "execution_id": "exec-123"}`))
	}))
	defer tsSuccess.Close()
	client := NewClient(tsSuccess.URL)
	execID, _ := client.ExecuteWorkflow(context.Background(), ActionPayload{})
	if execID != "exec-123" {
		t.Fatal("Expected exec-123")
	}

	tsFail := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer tsFail.Close()
	clientFail := NewClient(tsFail.URL)
	_, err := clientFail.ExecuteWorkflow(context.Background(), ActionPayload{})
	if err == nil {
		t.Fatal("Expected HTTP failure")
	}
}

func TestShuffleClient_PollingStates(t *testing.T) {
	os.Setenv("SAE_SHUFFLE_API_URL", "mock")
	os.Setenv("SAE_SHUFFLE_AUTH_TOKEN", "mock")
	os.Setenv("SAE_SHUFFLE_WORKFLOW_ID", "wf-123")

	// 1. SUCCESS
	apiSuccess := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"execution_id":"exec-123", "status": "SUCCESS", "result": "Done"}]`))
	}))
	defer apiSuccess.Close()
	
	client := NewClient("mock")
	client.APIURL = apiSuccess.URL
	client.AuthToken = "mock"
	client.WorkflowID = "wf-123"
	res, _ := client.PollExecutionStatus(context.Background(), "exec-123", 10*time.Millisecond)
	if res.Status != "SUCCEEDED" {
		t.Fatalf("Expected SUCCEEDED, got %s", res.Status)
	}

	// 2. EXECUTING -> SUCCESS
	calls := 0
	apiDynamic := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusOK)
		if calls < 3 {
			w.Write([]byte(`[{"execution_id":"exec-123", "status": "EXECUTING"}]`))
		} else {
			w.Write([]byte(`[{"execution_id":"exec-123", "status": "SUCCESS"}]`))
		}
	}))
	defer apiDynamic.Close()
	client.APIURL = apiDynamic.URL
	res, _ = client.PollExecutionStatus(context.Background(), "exec-123", 10*time.Millisecond)
	if res.Status != "SUCCEEDED" || calls != 3 {
		t.Fatalf("Expected SUCCEEDED after 3 calls, got %s", res.Status)
	}

	// 3. Timeout
	apiTimeout := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"execution_id":"exec-456", "status": "EXECUTING"}]`))
	}))
	defer apiTimeout.Close()
	client.APIURL = apiTimeout.URL
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	res2, err2 := client.PollExecutionStatus(ctx, "exec-456", 15*time.Millisecond)
	if err2 == nil || res2.Status != "TIMEOUT" {
		t.Fatalf("Expected TIMEOUT, got %s", res2.Status)
	}
}
