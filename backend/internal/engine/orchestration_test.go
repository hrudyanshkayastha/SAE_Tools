package engine

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"fmt"

	"sae-core/internal/shuffle"
	"sae-core/internal/storage"
	"sae-core/internal/verification"
	"sae-core/models"
)

func TestTargetVerification_Orchestration(t *testing.T) {
	runOrchestration := func(shuffleSuccess bool, targetPreBlocked bool, action string) (finalSavedState string, emittedTargetEvent *models.OCSFFinding) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/v1/hooks/webhook_sae_action" {
				if shuffleSuccess {
					w.Write([]byte("{\"success\":true,\"execution_id\":\"exec-123\"}"))
				} else {
					http.Error(w, "error", http.StatusInternalServerError)
				}
			} else if r.URL.Path == "/api/v1/workflows/test-wf/executions" {
				if shuffleSuccess {
					w.Write([]byte("[\n  {\"execution_id\":\"exec-123\",\"status\":\"FINISHED\"}\n]"))
				} else {
					w.Write([]byte("[\n  {\"execution_id\":\"exec-123\",\"status\":\"FAILED\"}\n]"))
				}
			} else {
				http.NotFound(w, r)
			}
		}))
		defer ts.Close()

		verifier := verification.NewLocalTestVerifier()
		if targetPreBlocked {
			verifier.BlockIP("10.0.0.99")
		}

		eng := NewEngine(&storage.Storage{}, ts.URL+"/api/v1/hooks/webhook_sae_action", "http://mock", "http://mock", verifier)
		eng.shuffleClient.APIURL = ts.URL + "/api/v1"
		eng.shuffleClient.AuthToken = "test-token"
		eng.shuffleClient.WorkflowID = "test-wf"

		resultPayload := shuffle.ActionPayload{Action: action, CorrelationID: "corr-123", Target: "10.0.0.99"}
		corr := &CorrelationContext{ID: "corr-123", Target: "10.0.0.99"}

		execID, err := eng.shuffleClient.ExecuteWorkflow(context.Background(), resultPayload)
		var shuffleFinalState string = "SUCCEEDED"
		
		if err != nil {
			fmt.Printf("ERROR Executing: %v\n", err)
			shuffleFinalState = "FAILED"
		} else {
			verifyCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			verifyResult, vErr := eng.shuffleClient.PollExecutionStatus(verifyCtx, execID, 100*time.Millisecond)
			if vErr != nil || verifyResult.Status != "SUCCEEDED" {
				fmt.Printf("ERROR Polling: %v, Status: %s\n", vErr, verifyResult.Status)
				shuffleFinalState = "FAILED"
			}
		}

		finalState := shuffleFinalState
		var tsEvent models.OCSFFinding

		if shuffleFinalState == "SUCCEEDED" && eng.verifier != nil {
			tsRes := eng.verifier.Verify(context.Background(), action, corr.Target)
			finalState = tsRes.State
			
			metadata := map[string]string{
				"action_type":    action,
				"expected_state": tsRes.Expected,
				"observed_state": tsRes.Observed,
			}
			rawBytes, _ := json.Marshal(metadata)
			
			tsEvent = models.OCSFFinding{
				Status:   tsRes.State,
				RawEvent: string(rawBytes),
			}
		}

		return finalState, &tsEvent
	}

	finalState, event := runOrchestration(true, true, "BLOCK_IP")
	if finalState != verification.StateVerified {
		t.Errorf("Expected VERIFIED, got %s", finalState)
	}

	finalState, event = runOrchestration(true, false, "BLOCK_IP")
	if finalState != verification.StateNotVerified {
		t.Errorf("Expected NOT_VERIFIED, got %s", finalState)
	}

	finalState, event = runOrchestration(false, true, "BLOCK_IP")
	if finalState != "FAILED" {
		t.Errorf("Expected FAILED, got %s", finalState)
	}
	if event != nil && event.Status == verification.StateVerified {
		t.Errorf("Expected event not to claim VERIFIED when Shuffle failed")
	}

	finalState, event = runOrchestration(true, true, "ISOLATE_HOST")
	if finalState != verification.StateUnknown {
		t.Errorf("Expected UNKNOWN, got %s", finalState)
	}
}
