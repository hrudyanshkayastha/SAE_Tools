
package langgraph

import (
	"encoding/json"
	"testing"
)

type SemanticTestCase struct {
	Name             string
	InputEvent       string
	ExpectedDecision string
	MinRiskScore     int
}

func TestSemanticEvaluation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping semantic evaluation in short mode as it requires a local LLM.")
	}

	testCases := []SemanticTestCase{
		{
			Name: "True Positive Brute Force",
			InputEvent: `{"activity_name": "Logon", "app_name": "SSH", "status": "Failure", "actor": {"user": {"name": "root"}}, "src_endpoint": {"ip": "185.15.20.100"}, "observables": [{"name": "failed_logins_1min", "value": "150"}], "message": "Multiple failed SSH login attempts for root from external IP"}`,
			ExpectedDecision: "ESCALATE_TO_HUMAN",
			MinRiskScore:     80, // High or Critical severity
		},
		{
			Name: "Benign Administrative Activity",
			InputEvent: `{"activity_name": "Logon", "app_name": "VPN", "status": "Success", "actor": {"user": {"name": "admin_bob"}}, "src_endpoint": {"ip": "10.0.50.25"}, "observables": [{"name": "mfa_used", "value": "true"}], "message": "Successful VPN login with MFA from internal IT subnet"}`,
			ExpectedDecision: "LOG_AND_MONITOR",
			MinRiskScore:     0, // Must be Info or Low
		},
		{
			Name: "Hallucination Trap",
			InputEvent: `{"activity_name": "Network Connection", "status": "Success", "message": "A random log about an apple pie baking contest in the cafeteria.", "observables": [{"name": "pie_flavor", "value": "apple"}]}`,
			ExpectedDecision: "LOG_AND_MONITOR",
			MinRiskScore:     0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := ExecuteReasoningGraph([]byte(tc.InputEvent))
			if err != nil {
				t.Fatalf("Failed to execute reasoning graph: %v", err)
			}

			// Validate LLM output format
			if result.Decision == "" {
				t.Fatalf("Empty decision returned by LLM pipeline")
			}

			if result.RiskScore < tc.MinRiskScore {
				t.Errorf("Expected minimum risk score %d, got %d", tc.MinRiskScore, result.RiskScore)
			}

			if tc.ExpectedDecision == "ESCALATE_TO_HUMAN" && result.RiskScore < 80 {
				t.Errorf("Expected high risk for %s, got score %d", tc.ExpectedDecision, result.RiskScore)
			}
			
			if tc.ExpectedDecision == "LOG_AND_MONITOR" && result.RiskScore >= 80 {
				t.Errorf("Expected low risk for %s, got score %d (High/Critical)", tc.ExpectedDecision, result.RiskScore)
			}

			if result.Decision != tc.ExpectedDecision {
				t.Errorf("Expected decision %q, got %q", tc.ExpectedDecision, result.Decision)
			}

			// Output the raw LLM recommendation for manual review in test logs
			llmRec, _ := json.MarshalIndent(result.LLMRecommendation, "", "  ")
			t.Logf("LLM Recommendation:\n%s", string(llmRec))
			t.Logf("Validation: %s", result.Validation)
		})
	}
}

