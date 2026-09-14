package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sae-core/models"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	PG    *sql.DB
	Redis *redis.Client
}

func InitStorage() (*Storage, error) {
	pgDSN := os.Getenv("SAE_PG_DSN")
	if pgDSN == "" {
		pgDSN = "postgres://sae:changeme_dev_only@localhost:5432/sae_db?sslmode=disable"
	}
	pgDB, err := sql.Open("postgres", pgDSN)
	if err != nil {
		return nil, fmt.Errorf("postgres open: %w", err)
	}

	pgDB.SetMaxOpenConns(25)
	pgDB.SetMaxIdleConns(5)
	pgDB.SetConnMaxLifetime(5 * time.Minute)

	_, err = pgDB.Exec(`CREATE TABLE IF NOT EXISTS correlations (
		correlation_id VARCHAR PRIMARY KEY,
		target VARCHAR,
		events_count INT,
		severity VARCHAR,
		status VARCHAR,
		created_at TIMESTAMP
	)`)
	if err != nil {
		return nil, fmt.Errorf("postgres correlations table init: %w", err)
	}

	_, err = pgDB.Exec(`CREATE TABLE IF NOT EXISTS decisions (
		correlation_id VARCHAR PRIMARY KEY,
		risk_score INT,
		validation TEXT,
		action VARCHAR,
		state VARCHAR DEFAULT 'REQUESTED',
		created_at TIMESTAMP
	)`)
	if err != nil {
		return nil, fmt.Errorf("postgres decisions table init: %w", err)
	}

	pgDB.Exec("ALTER TABLE decisions ADD COLUMN state VARCHAR DEFAULT 'REQUESTED'")

	_, err = pgDB.Exec(`CREATE TABLE IF NOT EXISTS sae_telemetry (
		event_id VARCHAR PRIMARY KEY,
		correlation_id VARCHAR,
		activity_name VARCHAR,
		severity VARCHAR,
		message TEXT,
		raw JSONB,
		timestamp TIMESTAMP
	)`)
	if err != nil {
		return nil, fmt.Errorf("postgres init table: %w", err)
	}

	redisAddr := os.Getenv("SAE_REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DialTimeout: 5 * time.Second,
		ReadTimeout: 3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})
	
	err = rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &Storage{PG: pgDB, Redis: rdb}, nil
}

func (s *Storage) PublishEvent(ctx context.Context, event models.OCSFFinding) error {
	if s == nil || s.Redis == nil {
		return nil
	}
	data, _ := json.Marshal(event)
	return s.Redis.XAdd(ctx, &redis.XAddArgs{
		Stream: "sae_events",
		Values: map[string]interface{}{"data": data},
	}).Err()
}

func (s *Storage) SaveTelemetry(event models.OCSFFinding) error {
	if s == nil || s.PG == nil {
		return nil
	}
	data, _ := json.Marshal(event)
	_, err := s.PG.Exec("INSERT INTO sae_telemetry (event_id, correlation_id, activity_name, severity, message, raw, timestamp) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (event_id) DO NOTHING",
		event.EventID, event.CorrelationID, event.ActivityName, event.Severity, event.Message, string(data), time.Now())
	return err
}

func (s *Storage) SaveCorrelation(corrID string, target string, count int, severity string, status string) error {
	if s == nil || s.PG == nil {
		return nil
	}
	_, err := s.PG.Exec("INSERT INTO correlations (correlation_id, target, events_count, severity, status, created_at) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (correlation_id) DO UPDATE SET events_count = EXCLUDED.events_count, status = EXCLUDED.status",
		corrID, target, count, severity, status, time.Now())
	return err
}

func (s *Storage) SaveDecision(corrID string, riskScore int, validation string, action string) error {
	if s == nil || s.PG == nil {
		return nil
	}
	_, err := s.PG.Exec("INSERT INTO decisions (correlation_id, risk_score, validation, action, created_at) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (correlation_id) DO UPDATE SET risk_score = EXCLUDED.risk_score, validation = EXCLUDED.validation, action = EXCLUDED.action",
		corrID, riskScore, validation, action, time.Now())
	return err
}

func (s *Storage) SaveResponseState(corrID string, state string) error {
	if s == nil || s.PG == nil {
		return nil
	}
	_, err := s.PG.Exec("UPDATE decisions SET state = $1 WHERE correlation_id = $2", state, corrID)
	return err
}
