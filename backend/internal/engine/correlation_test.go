package engine

import (
	"context"
	"testing"
	"time"

	"sae-core/internal/storage"
	"sae-core/models"
)

func TestCorrelation_MultiSourceSameTarget(t *testing.T) {
	eng := NewEngine(&storage.Storage{}, "http://mock", "http://mock", "http://mock", nil)
	ctx := context.Background()

	// 1. Wazuh event
	eng.correlate(ctx, models.OCSFFinding{
		EventID:      "1",
		ActivityName: "Login Failure",
		Severity:     "Low",
		Observables:  []models.Observable{{Type: "IP", Value: "10.0.0.5"}},
	})

	// 2. Zeek event
	eng.correlate(ctx, models.OCSFFinding{
		EventID:      "2",
		ActivityName: "Suspicious Connection",
		Severity:     "Medium",
		Observables:  []models.Observable{{Type: "IP", Value: "10.0.0.5"}},
	})

	// They should merge into the same correlation context
	if len(eng.correlations) != 1 {
		t.Fatalf("Expected 1 correlation chain, got %d", len(eng.correlations))
	}
	
	ctxData := eng.correlations["10.0.0.5"]
	if len(ctxData.Events) != 2 {
		t.Fatalf("Expected 2 events in chain, got %d", len(ctxData.Events))
	}
}

func TestCorrelation_DifferentTargets(t *testing.T) {
	eng := NewEngine(&storage.Storage{}, "http://mock", "http://mock", "http://mock", nil)
	ctx := context.Background()

	eng.correlate(ctx, models.OCSFFinding{
		EventID:      "1",
		Severity:     "Low",
		Observables:  []models.Observable{{Type: "IP", Value: "10.0.0.5"}},
	})

	eng.correlate(ctx, models.OCSFFinding{
		EventID:      "2",
		Severity:     "Low",
		Observables:  []models.Observable{{Type: "IP", Value: "192.168.1.10"}},
	})

	if len(eng.correlations) != 2 {
		t.Fatalf("Expected 2 separate correlation chains, got %d", len(eng.correlations))
	}
}

func TestCorrelation_TemporalWindow(t *testing.T) {
	eng := NewEngine(&storage.Storage{}, "http://mock", "http://mock", "http://mock", nil)
	ctx := context.Background()

	// First event
	eng.correlate(ctx, models.OCSFFinding{
		EventID:      "1",
		Severity:     "Low",
		Observables:  []models.Observable{{Type: "IP", Value: "10.0.0.5"}},
	})

	// Manually age the context to simulate > 5 minutes
	eng.correlations["10.0.0.5"].CreatedAt = time.Now().Add(-6 * time.Minute)

	// Second event should force a new correlation chain due to temporal expiration
	eng.correlate(ctx, models.OCSFFinding{
		EventID:      "2",
		Severity:     "Low",
		Observables:  []models.Observable{{Type: "IP", Value: "10.0.0.5"}},
	})

	// Wait, the map key is overwritten with the new one.
	// So we still have 1 map key, but its Events slice should only have 1 event!
	if len(eng.correlations["10.0.0.5"].Events) != 1 {
		t.Fatalf("Expected temporal window to reset chain and have 1 event, got %d", len(eng.correlations["10.0.0.5"].Events))
	}
	if eng.correlations["10.0.0.5"].Events[0].EventID != "2" {
		t.Fatalf("Expected new chain to contain EventID 2")
	}
}

func TestCorrelation_EscalationClear(t *testing.T) {
	// We want to test that triggering AI (on 3rd event or High severity) clears the map
	// However, triggerAI calls langgraph/Ollama. This would block or fail if Ollama isn't running.
	// But triggerAI will fail quickly because the URL/path is blocked or we just mock triggerAI.
	// Actually, if it calls langgraph, we will see it fail gracefully and then return.
	// Then it deletes the key.
	eng := NewEngine(&storage.Storage{}, "http://mock", "http://mock", "http://mock", nil)
	ctx := context.Background()

	eng.correlate(ctx, models.OCSFFinding{
		EventID:      "1",
		Severity:     "Low",
		Observables:  []models.Observable{{Type: "IP", Value: "10.0.0.5"}},
	})

	eng.correlate(ctx, models.OCSFFinding{
		EventID:      "2",
		Severity:     "High", // Trigger AI immediately
		Observables:  []models.Observable{{Type: "IP", Value: "10.0.0.5"}},
	})

	if len(eng.correlations) != 0 {
		t.Fatalf("Expected correlation map to be cleared after High severity trigger, got %d", len(eng.correlations))
	}
}
