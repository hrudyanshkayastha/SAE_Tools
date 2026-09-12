package ueba

import (
	"context"
	"fmt"
	"log"
	"math"
	"sae-core/internal/storage"
	"sae-core/models"
	"time"
	"encoding/json"
	"github.com/google/uuid"
)

type UEBAEngine struct {
	store *storage.Storage
}

func NewUEBAEngine(store *storage.Storage) *UEBAEngine {
	return &UEBAEngine{store: store}
}

// Start runs the UEBA anomaly detection loop periodically
func (u *UEBAEngine) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			u.analyzeBehavior(ctx)
		}
	}
}

type EntityStat struct {
	Type      string
	Value     string
	Buckets   map[int64]int // Unix minute -> count
	RiskScore int
}

func (u *UEBAEngine) analyzeBehavior(ctx context.Context) {
	// Query telemetry from the last 60 minutes
	rows, err := u.store.PG.QueryContext(ctx, "SELECT raw FROM sae_telemetry WHERE timestamp >= NOW() - INTERVAL '60 minutes'")
	if err != nil {
		log.Printf("[UEBA] Failed to query telemetry: %v", err)
		return
	}
	defer rows.Close()

	entities := make(map[string]*EntityStat)

	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err != nil {
			continue
		}
		var finding models.OCSFFinding
		if err := json.Unmarshal([]byte(raw), &finding); err != nil {
			continue
		}

		minuteBucket := finding.Time.Truncate(time.Minute).Unix()

		for _, obs := range finding.Observables {
			if obs.Type != "User" && obs.Type != "IP" {
				continue
			}
			key := obs.Type + ":" + obs.Value
			if _, exists := entities[key]; !exists {
				entities[key] = &EntityStat{
					Type:    obs.Type,
					Value:   obs.Value,
					Buckets: make(map[int64]int),
				}
			}
			entities[key].Buckets[minuteBucket] += 1
		}
	}

	currentMinute := time.Now().Truncate(time.Minute).Unix()

	for _, stat := range entities {
		isAnomaly, mean, stddev := EvaluateDeviation(stat, currentMinute)
		if isAnomaly {
			currentCount := stat.Buckets[currentMinute]
			log.Printf("[UEBA v2] ANOMALY DETECTED for %s %s: %d events this minute. Risk Score: %d (Baseline Mean: %.2f, StdDev: %.2f)", stat.Type, stat.Value, currentCount, stat.RiskScore, mean, stddev)
			
			// Publish anomaly to the Event Fabric
			anomaly := models.OCSFFinding{
				EventID:       uuid.New().String(),
				ActivityName:  "UEBA Behavioral Deviation",
				Time:          time.Now(),
				Severity:      "High",
				Status:        "New",
				Message:       fmt.Sprintf("Behavioral deviation detected for %s %s. Event rate spike: %d/min. Risk Score: %d", stat.Type, stat.Value, currentCount, stat.RiskScore),
			}
			anomaly.Observables = append(anomaly.Observables, models.Observable{
				Type:  stat.Type,
				Value: stat.Value,
			})
			anomaly.Metadata.Product = "SAE UEBA Engine v2"
			
			u.store.PublishEvent(ctx, anomaly)
		}
	}
}

// EvaluateDeviation calculates baseline and returns true if the current minute is an anomaly.
func EvaluateDeviation(stat *EntityStat, currentMinute int64) (bool, float64, float64) {
	var sum int
	var count int
	for t, c := range stat.Buckets {
		if t != currentMinute {
			sum += c
			count++
		}
	}

	if count < 3 {
		return false, 0, 0
	}

	mean := float64(sum) / float64(count)
	
	var variance float64
	for t, c := range stat.Buckets {
		if t != currentMinute {
			diff := float64(c) - mean
			variance += diff * diff
		}
	}
	variance = variance / float64(count)
	stddev := math.Sqrt(variance)

	currentCount := stat.Buckets[currentMinute]

	// Threshold: Mean + 3 standard deviations (or at least 5 absolute events to avoid low-volume noise)
	threshold := mean + (3 * stddev)
	if threshold < 5 {
		threshold = 5
	}

	isDeviation := float64(currentCount) > threshold

	// UEBA v2: Calculate Risk Magnitude (1-100)
	if isDeviation {
		magnitude := (float64(currentCount) - mean) / stddev
		stat.RiskScore = int(math.Min(100, magnitude * 10))
	} else {
		stat.RiskScore = 0
	}

	return isDeviation, mean, stddev
}
