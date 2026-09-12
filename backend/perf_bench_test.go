package main

import (
	"context"
	"fmt"
	"sae-core/internal/storage"
	"sae-core/models"
	"testing"
	"time"

	"github.com/google/uuid"
)

// BenchmarkRedisPublish measures raw event ingestion throughput into the Redis stream.
func BenchmarkRedisPublish(b *testing.B) {
	store, err := storage.InitStorage()
	if err != nil {
		b.Fatalf("Failed to init storage: %v", err)
	}
	defer store.PG.Close()
	defer store.Redis.Close()

	ctx := context.Background()
	event := models.OCSFFinding{
		EventID:      uuid.New().String(),
		ActivityName: "Benchmark Test",
		Severity:     "Low",
		Message:      "Synthetic load testing",
		Time:         time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		event.EventID = uuid.New().String()
		err := store.PublishEvent(ctx, event)
		if err != nil {
			b.Fatalf("Redis publish failed: %v", err)
		}
	}
}

// BenchmarkPostgresWrite measures PostgreSQL insertion latency and throughput.
func BenchmarkPostgresWrite(b *testing.B) {
	store, err := storage.InitStorage()
	if err != nil {
		b.Fatalf("Failed to init storage: %v", err)
	}
	defer store.PG.Close()
	defer store.Redis.Close()

	event := models.OCSFFinding{
		EventID:       uuid.New().String(),
		CorrelationID: "CORR-BENCH",
		ActivityName:  "Benchmark DB Test",
		Severity:      "Low",
		Message:       "Synthetic load testing",
		Time:          time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		event.EventID = fmt.Sprintf("bench-event-%d", i)
		err := store.SaveTelemetry(event)
		if err != nil {
			b.Fatalf("Postgres write failed: %v", err)
		}
	}
}
