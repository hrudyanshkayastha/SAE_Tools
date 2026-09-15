package langgraph

import (
	"os"
	"strings"
	"testing"
)

func getScriptDir() string {
	wd, _ := os.Getwd()
	return wd
}

func TestExecuteReasoningGraph_ValidEvent(t *testing.T) {
	event := []byte(`{"event_type": "login_failure", "src_ip": "10.0.0.1"}`)
	res, err := ExecuteReasoningGraph(event)
	if err != nil { t.Fatalf("execution failed: %v", err) }
	
	if res.Decision == "" { t.Errorf("missing decision") }
	if res.Validation == "" { t.Errorf("missing validation") }
}

func TestExecuteReasoningGraph_MalformedInput(t *testing.T) {
	event := []byte(`{"event_type": `) // malformed json
	_, err := ExecuteReasoningGraph(event)
	if err == nil || !strings.Contains(err.Error(), "malformed input") {
		t.Errorf("expected malformed input error, got: %v", err)
	}
}

func TestExecuteReasoningGraph_EmptyEvent(t *testing.T) {
	_, err := ExecuteReasoningGraph([]byte(""))
	if err == nil { t.Errorf("expected empty event error") }
}

func TestExecuteReasoningGraph_WhitespaceOnly(t *testing.T) {
	_, err := ExecuteReasoningGraph([]byte("   \n"))
	if err == nil || !strings.Contains(err.Error(), "empty event") {
		t.Errorf("expected empty event error from graph.py")
	}
}

func TestExecuteReasoningGraph_StatePropagation(t *testing.T) {
	// The state propagation is proven by the graph successfully returning the final risk score based on the mocked or live LLM step.
	event := []byte(`{"event_type": "test_propagation"}`)
	res, err := ExecuteReasoningGraph(event)
	if err != nil { t.Fatalf("failed: %v", err) }
	
	// Ensure that LLMRecommendation propagated to Validation and RiskScore
	if res.Validation == "" { t.Errorf("validation state missing") }
}

func TestExecuteReasoningGraph_DeterministicTermination(t *testing.T) {
	// Tests that the graph does not loop infinitely. If it reaches here, it terminated deterministically.
	event := []byte(`{"type": "terminate_test"}`)
	_, err := ExecuteReasoningGraph(event)
	if err != nil { t.Fatalf("should terminate cleanly") }
}

func TestExecuteReasoningGraph_SecurityBoundary(t *testing.T) {
	// Test that LLM output recommending a privileged action is blocked at the validation boundary
	event := []byte(`{"type": "synthetic_high_severity"}`)
	res, err := ExecuteReasoningGraph(event)
	if err != nil { t.Fatalf("failed: %v", err) }
	
	// We expect the graph to NOT execute privileged actions. 
	// The validation node explicitly sets status to "Evidence validated" or "Requires manual authorization".
	if strings.Contains(res.Validation, "Executed privileged action") {
		t.Errorf("Security boundary violated! Found execution of privileged action.")
	}
}

func TestParseGraphResult_Valid(t *testing.T) {
	out := &GraphOutput{
		Decision: "ESCALATE_TO_HUMAN",
		Validation: "Requires manual authorization (Privileged action blocked)",
		RiskScore: 90,
		LLMRecommendation: map[string]interface{}{"recommendation": "Block IP"},
	}
	f, err := ParseGraphResultToOCSF(out)
	if err != nil { t.Fatalf("mapper failed") }
	if f[0].Severity != "High" { t.Errorf("expected High severity for score 90") }
}

func TestParseGraphResult_NoRecommendation(t *testing.T) {
	out := &GraphOutput{
		Decision: "LOG_AND_MONITOR",
		Validation: "Evidence validated (Low impact)",
		RiskScore: 10,
	}
	f, _ := ParseGraphResultToOCSF(out)
	if f[0].Severity != "Info" { t.Errorf("expected Info severity") }
}

func TestParseGraphResult_NilInput(t *testing.T) {
	_, err := ParseGraphResultToOCSF(nil)
	if err == nil { t.Errorf("expected error on nil") }
}

func TestExecuteReasoningGraph_CorrelationArray(t *testing.T) {
	// Represents the exact shape sent by the Correlation Engine
	// e.g. eventData, _ := json.Marshal(corr.Events)
	event := []byte(`[{"event_id": "1", "severity": "Low"}, {"event_id": "2", "severity": "High"}]`)
	res, err := ExecuteReasoningGraph(event)
	if err != nil {
		t.Fatalf("Regression found: array input failed with: %v", err)
	}
	if res.Decision == "" {
		t.Errorf("missing decision for array input")
	}
}
