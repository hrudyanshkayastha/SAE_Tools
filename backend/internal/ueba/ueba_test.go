package ueba

import (
	"testing"
	"time"
)

func TestEvaluateDeviation_NoAnomaly(t *testing.T) {
	currentMinute := time.Now().Truncate(time.Minute).Unix()

	stat := &EntityStat{
		Type:  "User",
		Value: "admin",
		Buckets: map[int64]int{
			currentMinute - 180: 2,
			currentMinute - 120: 3,
			currentMinute - 60:  2,
			currentMinute:       3, // normal
		},
	}

	isAnomaly, mean, _ := EvaluateDeviation(stat, currentMinute)
	if isAnomaly {
		t.Errorf("Expected false, got true")
	}
	
	expectedMean := float64(2+3+2) / 3.0
	if mean != expectedMean {
		t.Errorf("Expected mean %f, got %f", expectedMean, mean)
	}
}

func TestEvaluateDeviation_AnomalySpike(t *testing.T) {
	currentMinute := time.Now().Truncate(time.Minute).Unix()

	stat := &EntityStat{
		Type:  "User",
		Value: "admin",
		Buckets: map[int64]int{
			currentMinute - 180: 2,
			currentMinute - 120: 3,
			currentMinute - 60:  2,
			currentMinute:       20, // huge spike
		},
	}

	isAnomaly, _, _ := EvaluateDeviation(stat, currentMinute)
	if !isAnomaly {
		t.Errorf("Expected true for anomaly, got false")
	}
}

func TestEvaluateDeviation_InsufficientHistory(t *testing.T) {
	currentMinute := time.Now().Truncate(time.Minute).Unix()

	// Only 2 historical buckets, needs at least 3
	stat := &EntityStat{
		Type:  "User",
		Value: "admin",
		Buckets: map[int64]int{
			currentMinute - 120: 3,
			currentMinute - 60:  2,
			currentMinute:       20, 
		},
	}

	isAnomaly, _, _ := EvaluateDeviation(stat, currentMinute)
	if isAnomaly {
		t.Errorf("Expected false due to insufficient history, got true")
	}
}

func TestEvaluateDeviation_BelowAbsoluteThreshold(t *testing.T) {
	currentMinute := time.Now().Truncate(time.Minute).Unix()

	stat := &EntityStat{
		Type:  "User",
		Value: "admin",
		Buckets: map[int64]int{
			currentMinute - 180: 0,
			currentMinute - 120: 0,
			currentMinute - 60:  0,
			currentMinute:       2, // > mean + 2*stddev (which is 0), but absolute count is < 5
		},
	}

	isAnomaly, _, _ := EvaluateDeviation(stat, currentMinute)
	if isAnomaly {
		t.Errorf("Expected false because count is < 5 (absolute noise threshold), got true")
	}
}
