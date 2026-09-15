package ueba

import (
	"testing"
	"time"
)

// BenchmarkEvaluateDeviation measures the mathematical latency of the UEBA deviation algorithm.
func BenchmarkEvaluateDeviation(b *testing.B) {
	currentMinute := time.Now().Truncate(time.Minute).Unix()

	// Build a mock entity with 60 minutes of history
	stat := &EntityStat{
		Type:    "User",
		Value:   "admin",
		Buckets: make(map[int64]int),
	}
	
	for i := 1; i <= 60; i++ {
		stat.Buckets[currentMinute - int64(i*60)] = 2
	}
	stat.Buckets[currentMinute] = 10 // Spike

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EvaluateDeviation(stat, currentMinute)
	}
}
