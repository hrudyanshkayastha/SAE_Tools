package ollama

import (
	"testing"
)

func TestRealOllamaExecution(t *testing.T) {
	// The user requested explicit proof of ACTUAL Ollama execution
	client := NewClient("http://127.0.0.1:11434", "llama3.2:3b")

	resp, err := client.GenerateThreatAnalysis("Test event: unauthorized login attempt.")
	if err != nil {
		t.Fatalf("Real Ollama execution failed: %v", err)
	}

	if resp.Model != "llama3.2:3b" {
		t.Errorf("Expected model llama3.2:3b, got %s", resp.Model)
	}

	if len(resp.Response) == 0 {
		t.Errorf("Received empty response from real Ollama API")
	}

	// Just checking if it returned something parsable or raw
	t.Logf("Raw Ollama Response: %s", resp.Response)
	
	// Map to OCSF to prove end-to-end processing
	f, err := ParseAIResponseToOCSF(resp)
	if err != nil {
		t.Fatalf("Failed to parse into OCSF: %v", err)
	}

	if len(f) != 1 {
		t.Fatalf("Expected 1 finding")
	}
	
	if f[0].ActivityName != "AI Threat Investigation" {
		t.Errorf("OCSF Activity Name mismatch")
	}
	
	t.Logf("Final OCSF Finding: %+v", f[0])
}
