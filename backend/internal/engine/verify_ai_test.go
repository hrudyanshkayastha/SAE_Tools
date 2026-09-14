package engine

import (
	"context"
	"testing"
	"time"

	"sae-core/internal/storage"
	"sae-core/models"
)

func TestVerifyAI_ProductionPath_Attack(t *testing.T) {
	eng := NewEngine(&storage.Storage{}, "http://mock", "http://mock", "http://mock")
	ctx := context.Background()

	event := models.OCSFFinding{
		EventID:      "test-1234-uuid",
		ActivityName: "Brute Force Authentication Bypass",
		Severity:     "High",
		Message:      "150 failed SSH logins from external IP 203.0.113.45",
		Observables: []models.Observable{{Type: "IP", Value: "203.0.113.45"}},
		Time: time.Now(),
	}

	// This triggers immediately because Severity is High
	result := eng.correlate(ctx, event)
	
	if result == nil {
		t.Fatalf("LangGraph execution failed to return a result!")
	}
	if result.Decision == "" {
		t.Errorf("Policy decision is missing")
	}
	if result.RiskScore < 80 {
		t.Errorf("Expected RiskScore >= 80, got: %d", result.RiskScore)
	}
	if result.Decision != "ESCALATE_TO_HUMAN" {
		t.Errorf("Expected ESCALATE_TO_HUMAN, got: %s", result.Decision)
	}
}

func TestVerifyAI_ProductionPath_Benign(t *testing.T) {
	eng := NewEngine(&storage.Storage{}, "http://mock", "http://mock", "http://mock")
	ctx := context.Background()

	event := models.OCSFFinding{
		EventID:      "test-benign",
		ActivityName: "VPN Login Success",
		Severity:     "Info", // Info won't trigger immediately unless we reach 3 events, wait! 
		// Actually, let's inject 3 benign events to force correlation trigger!
		Message:      "Successful VPN login",
		Observables: []models.Observable{{Type: "User", Value: "admin"}},
		Time: time.Now(),
	}

	eng.correlate(ctx, event)
	eng.correlate(ctx, event)
	result := eng.correlate(ctx, event) // 3rd event triggers AI
	
	if result == nil { t.Fatalf("Expected result on 3rd benign event") }
	if result.Decision != "LOG_AND_MONITOR" && result.Decision != "none" {
		t.Errorf("Expected LOG_AND_MONITOR, got: %s", result.Decision)
	}
	if result.RiskScore >= 80 {
		t.Errorf("Expected low risk score, got: %d", result.RiskScore)
	}
}

func TestVerifyAI_ProductionPath_CrossTool(t *testing.T) {
	eng := NewEngine(&storage.Storage{}, "http://mock", "http://mock", "http://mock")
	ctx := context.Background()

	event1 := models.OCSFFinding{
		EventID:      "test-multi-1",
		ActivityName: "Zeek Connection established",
		Severity:     "Low",
		Message:      "Inbound connection from 198.51.100.22",
		Observables: []models.Observable{{Type: "IP", Value: "198.51.100.22"}},
		Time: time.Now(),
	}
	event2 := models.OCSFFinding{
		EventID:      "test-multi-2",
		ActivityName: "Suricata ET DROP Dshield Block Listed Source",
		Severity:     "High", // High triggers it immediately!
		Message:      "Known malicious IP 198.51.100.22 attempted contact",
		Observables: []models.Observable{{Type: "IP", Value: "198.51.100.22"}},
		Time: time.Now(),
	}

	res1 := eng.correlate(ctx, event1)
	if res1 != nil {
		t.Fatalf("Expected no trigger on first Low event")
	}

	result := eng.correlate(ctx, event2) // Triggers and combines both!
	
	if result == nil { t.Fatalf("Expected AI trigger on High event") }
	if result.RiskScore < 50 {
		t.Errorf("Expected elevated risk score for known malicious IP, got: %d", result.RiskScore)
	}
}
