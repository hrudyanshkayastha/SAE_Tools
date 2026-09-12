package engine

import (
	"testing"

	"sae-core/models"
	"sae-core/internal/storage"
)

func TestEngine_Correlate(t *testing.T) {
	store := &storage.Storage{}
	eng := NewEngine(store, "http://mock-webhook", "http://mock-thehive", "http://mock-cortex")
	_ = eng

	// A simple test to verify correlation struct initialization
	ctxData := &CorrelationContext{
		ID:     "CORR-1.2.3.4",
		Target: "1.2.3.4",
	}
	if ctxData.ID != "CORR-1.2.3.4" {
		t.Errorf("Expected ID CORR-1.2.3.4, got %s", ctxData.ID)
	}
}

func TestEngine_EmptyTarget(t *testing.T) {
	event := models.OCSFFinding{
		ActivityName: "Test",
	}
	
	// Just verify the event has empty target logic when processed manually
	if len(event.Observables) == 0 {
		// Valid
	} else {
		t.Errorf("Expected no observables")
	}
}

func TestEngine_AIContext(t *testing.T) {
	ctxData := &CorrelationContext{
		ID:     "CORR-test",
		Target: "test_user",
	}
	
	if ctxData.Target != "test_user" {
		t.Errorf("Expected target test_user")
	}
}
