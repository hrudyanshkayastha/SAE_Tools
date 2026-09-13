package engine

import (
	"context"
	"fmt"
	"testing"

	"sae-core/internal/storage"
	"sae-core/models"
)

// TEST 1 - OCSF / CORRELATION MICROBENCHMARK
// Measures deterministic correlation throughput (events/sec) before the AI boundary.
func BenchmarkEngine_Correlate(b *testing.B) {
	eng := NewEngine(&storage.Storage{}, "http://mock", "http://mock", "http://mock")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Use a large target space to avoid hitting the 3-event AI boundary
		target := fmt.Sprintf("10.0.0.%d", i%(b.N/2 + 1))
		eng.correlate(ctx, models.OCSFFinding{
			EventID:      fmt.Sprintf("%d", i),
			ActivityName: "Login Failure",
			Severity:     "Low",
			Observables:  []models.Observable{{Type: "IP", Value: target}},
		})
	}
}
