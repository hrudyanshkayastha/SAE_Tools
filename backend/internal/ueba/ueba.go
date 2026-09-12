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
	Type  string
	Value string
	Buckets map[int64]int // Unix minute -> count
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
			log.Printf("[UEBA] ANOMALY DETECTED for %s %s: %d events this minute (Baseline Mean: %.2f, StdDev: %.2f)", stat.Type, stat.Value, currentCount, mean, stddev)
			
			// Publish anomaly to the Event Fabric
			anomaly := models.OCSFFinding{
				EventID:       uuid.New().String(),
				ActivityName:  "UEBA Behavioral Deviation",
				Time:          time.Now(),
				Severity:      "High",
				Status:        "New",
				Message:       fmt.Sprintf("Behavioral deviation detected for %s %s. Event rate spike: %d/min (Historical mean: %.2f/min, Stddev: %.2f)", stat.Type, stat.Value, currentCount, mean, stddev),
			}
			anomaly.Observables = append(anomaly.Observables, models.Observable{
				Type:  stat.Type,
				Value: stat.Value,
			})
			anomaly.Metadata.Product = "SAE UEBA Engine"
			
			u.store.PublishEvent(ctx, anomaly)
		}
	}
}

// EvaluateDeviation calculates baseline and returns true if the current minute is an anomaly.
func EvaluateDeviation(stat *EntityStat, currentMinute int64) (bool, float64, float64) {
	var total int
	var historicalBuckets int
	
	for ts, count := range stat.Buckets {
		if ts < currentMinute {
			total += count
			historicalBuckets++
		}
	}

	// Require at least 3 minutes of history to establish a baseline
	if historicalBuckets < 3 {
		return false, 0, 0
	}

	mean := float64(total) / float64(historicalBuckets)

	var varianceSum float64
	for ts, count := range stat.Buckets {
		if ts < currentMinute {
			varianceSum += math.Pow(float64(count)-mean, 2)
		}
	}
	variance := varianceSum / float64(historicalBuckets)
	stddev := math.Sqrt(variance)

	currentCount := stat.Buckets[currentMinute]

	// Threshold: count must be > mean + 2*stddev, and absolute count >= 5 to avoid low-volume noise
	threshold := mean + (2 * stddev)
	
	if float64(currentCount) > threshold && currentCount >= 5 {
		return true, mean, stddev
	}
	
	return false, mean, stddev
}
