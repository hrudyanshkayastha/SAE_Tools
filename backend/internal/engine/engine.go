package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sae-core/internal/langgraph"
	"sae-core/internal/shuffle"
	"sae-core/internal/storage"
	"sae-core/internal/thehive"
	"sae-core/internal/cortex"
	"sae-core/models"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Engine struct {
	store         *storage.Storage
	shuffleClient *shuffle.Client
	thehiveClient *thehive.Client
	cortexClient  *cortex.Client

	// Correlation state (in-memory for simple deterministic correlation)
	correlations map[string]*CorrelationContext
	mu           sync.Mutex
}

type CorrelationContext struct {
	ID        string
	Events    []models.OCSFFinding
	Target    string
	CreatedAt time.Time
}

func NewEngine(store *storage.Storage, shuffleWebhook, thehiveURL, cortexURL string) *Engine {
	return &Engine{
		store:         store,
		shuffleClient: shuffle.NewClient(shuffleWebhook),
		thehiveClient: thehive.NewClient(thehiveURL, "mock-api-key"),
		cortexClient:  cortex.NewClient(cortexURL, "mock-api-key"),
		correlations:  make(map[string]*CorrelationContext),
	}
}

func (e *Engine) Start(ctx context.Context) {
	// Create consumer group (ignore error if exists)
	e.store.Redis.XGroupCreateMkStream(ctx, "sae_events", "sae_group", "0")

	for {
		select {
		case <-ctx.Done():
			return
		default:
			res, err := e.store.Redis.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    "sae_group",
				Consumer: "consumer1",
				Streams:  []string{"sae_events", ">"},
				Count:    10,
				Block:    2 * time.Second,
			}).Result()

			if err != nil && err != redis.Nil {
				log.Printf("Redis read error: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			if len(res) > 0 {
				for _, msg := range res[0].Messages {
					e.processMessage(ctx, msg)
					e.store.Redis.XAck(ctx, "sae_events", "sae_group", msg.ID)
				}
			}
		}
	}
}

func (e *Engine) processMessage(ctx context.Context, msg redis.XMessage) {
	dataStr, ok := msg.Values["data"].(string)
	if !ok {
		return
	}

	var event models.OCSFFinding
	if err := json.Unmarshal([]byte(dataStr), &event); err != nil {
		log.Printf("Failed to unmarshal event: %v", err)
		return
	}

	// 2. Correlate
	e.correlate(ctx, event)
}

func (e *Engine) correlate(ctx context.Context, event models.OCSFFinding) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Extract target (e.g. IP or User) to correlate deterministically
	var target string
	for _, obs := range event.Observables {
		if obs.Type == "IP" || obs.Type == "User" || obs.Type == "Resource" {
			target = obs.Value
			break
		}
	}

	if target == "" {
		target = "UNASSIGNED"
	}

	// For test isolation, we correlate by target string
	corrKey := target
	ctxData, exists := e.correlations[corrKey]
	if !exists {
		ctxData = &CorrelationContext{
			ID:        "CORR-" + target,
			Target:    target,
			CreatedAt: time.Now(),
		}
		e.correlations[corrKey] = ctxData
	}

	ctxData.Events = append(ctxData.Events, event)

	// Update Event with correlation ID and save to ClickHouse Telemetry
	event.CorrelationID = ctxData.ID
	e.store.SaveTelemetry(event)

	// Save Correlation state to Postgres
	e.store.SaveCorrelation(ctxData.ID, ctxData.Target, len(ctxData.Events), event.Severity, "ACTIVE")

	// If we have enough context or a critical event, trigger AI Investigation
	if event.Severity == "Critical" || event.Severity == "High" {
		log.Printf("[CORRELATION] Triggering LangGraph AI Investigation for Correlation %s based on High/Critical event", ctxData.ID)
		e.triggerAI(ctxData)
	} else if len(ctxData.Events) >= 3 {
		log.Printf("[CORRELATION] Threshold reached for Correlation %s (3+ events). Triggering AI.", ctxData.ID)
		e.triggerAI(ctxData)
		// Reset for demo purposes
		delete(e.correlations, corrKey)
	}
}

func (e *Engine) triggerAI(corr *CorrelationContext) {
	// Consolidate data for LangGraph
	summary := fmt.Sprintf("Correlation %s targeting %s with %d events. Latest: %s", corr.ID, corr.Target, len(corr.Events), corr.Events[len(corr.Events)-1].Message)

	eventData, _ := json.Marshal(map[string]string{
		"event_source": "Wazuh",
		"event_type":   "Authentication Bypass",
		"description":  summary,
		"severity":     "High",
		"source_ip":    corr.Target,
	})

	result, err := langgraph.ExecuteReasoningGraph(eventData, "E:\\New folder\\SAE_Tools\\SAE\\backend\\internal\\langgraph")
	if err != nil {
		log.Printf("[AI REASONING] LangGraph execution failed: %v", err)
		return
	}

	log.Printf("[AI DECISION] Risk Score: %d, Validation: %s, Action: %s", result.RiskScore, result.Validation, result.Decision)

	// Save Decision to Postgres
	e.store.SaveDecision(corr.ID, result.RiskScore, result.Validation, result.Decision)
	e.store.SaveCorrelation(corr.ID, corr.Target, len(corr.Events), "High", "RESOLVED")

	// Policy Engine bounds checking
	if result.Decision == "block_ip" || result.Decision == "isolate_host" {
		log.Printf("[POLICY ENGINE] DANGER: LLM recommended highly privileged action: %s. Action blocked. Requiring human authorization.", result.Decision)
		return
	}
	
	log.Printf("[POLICY ENGINE] Action %s is within safe bounds.", result.Decision)
	
	if result.Decision == "monitor" || result.Decision == "none" {
		log.Printf("[RESPONSE] No active execution required for %s", result.Decision)
		return
	}

	log.Printf("[RESPONSE] Routing authorized action %s to Shuffle...", result.Decision)
	
	payload := shuffle.ActionPayload{
		CorrelationID: corr.ID,
		Action:        result.Decision,
		Target:        corr.Target,
	}

	// Case Management & Enrichment
	log.Printf("[CASE MANAGEMENT] Opening Incident in TheHive for %s...", corr.ID)
	caseCtx, caseCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer caseCancel()
	
	caseID, err := e.thehiveClient.CreateCase(caseCtx, thehive.CasePayload{
		Title:       fmt.Sprintf("SAE AI Escalation: %s", corr.Target),
		Description: fmt.Sprintf("AI Score: %d. Recommendation: %s", result.RiskScore, result.Decision),
		Severity:    3,
	})
	if err != nil {
		log.Printf("[CASE MANAGEMENT] [BLOCKED] TheHive execution failed: %v", err)
	} else {
		log.Printf("[CASE MANAGEMENT] Successfully created TheHive Case: %s", caseID)
		
		// If case created successfully, run Cortex analysis on the target
		log.Printf("[THREAT INTEL] Submitting Cortex Job for target: %s", corr.Target)
		jobID, err := e.cortexClient.SubmitJob(caseCtx, "VirusTotal_GetReport_3_0", cortex.JobPayload{
			Data:     corr.Target,
			DataType: "ip",
			Tlp:      2,
		})
		if err != nil {
			log.Printf("[THREAT INTEL] [BLOCKED] Cortex execution failed: %v", err)
		} else {
			log.Printf("[THREAT INTEL] Successfully submitted Cortex Job: %s", jobID)
		}
	}

	// Wait up to 5 seconds for webhook push
	execCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = e.shuffleClient.ExecuteWorkflow(execCtx, payload)
	if err != nil {
		log.Printf("[RESPONSE] Shuffle execution failed: %v", err)
		// Send failure event to fabric
		failEvent := models.OCSFFinding{
			EventID:      corr.ID + "-fail",
			CorrelationID: corr.ID,
			ActivityName: "Response Execution Failed",
			Severity:     "High",
			Message:      fmt.Sprintf("Shuffle execution failed for action %s: %v", result.Decision, err),
		}
		e.store.PublishEvent(execCtx, failEvent)
	} else {
		log.Printf("[RESPONSE] Shuffle workflow successfully triggered for %s", result.Decision)
		// We expect the consumer to pick up the result and send the OCSF event later.
	}
}
