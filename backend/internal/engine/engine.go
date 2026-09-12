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

	"github.com/google/uuid"
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

	e.correlate(ctx, event)
}

func (e *Engine) correlate(ctx context.Context, event models.OCSFFinding) {
	e.mu.Lock()
	defer e.mu.Unlock()

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
	event.CorrelationID = ctxData.ID
	e.store.SaveTelemetry(event)
	e.store.SaveCorrelation(ctxData.ID, ctxData.Target, len(ctxData.Events), event.Severity, "ACTIVE")

	if event.Severity == "Critical" || event.Severity == "High" {
		e.triggerAI(ctxData)
	} else if len(ctxData.Events) >= 3 {
		e.triggerAI(ctxData)
		delete(e.correlations, corrKey)
	}
}

func (e *Engine) triggerAI(corr *CorrelationContext) {
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

	// Save initial Decision to Postgres as REQUESTED
	e.store.SaveDecision(corr.ID, result.RiskScore, result.Validation, result.Decision)
	e.store.SaveCorrelation(corr.ID, corr.Target, len(corr.Events), "High", "RESOLVED")
	e.store.SaveResponseState(corr.ID, "REQUESTED")

	// Policy Engine bounds checking
	if result.Decision == "block_ip" || result.Decision == "isolate_host" {
		log.Printf("[POLICY ENGINE] DANGER: LLM recommended highly privileged action: %s. Action blocked.", result.Decision)
		e.store.SaveResponseState(corr.ID, "REJECTED")
		return
	}
	
	e.store.SaveResponseState(corr.ID, "AUTHORIZED")
	
	if result.Decision == "monitor" || result.Decision == "none" {
		e.store.SaveResponseState(corr.ID, "COMPLETED_NO_ACTION")
		return
	}

	payload := shuffle.ActionPayload{
		CorrelationID: corr.ID,
		Action:        result.Decision,
		Target:        corr.Target,
	}

	execCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	e.store.SaveResponseState(corr.ID, "EXECUTING")
	execID, err := e.shuffleClient.ExecuteWorkflow(execCtx, payload)
	if err != nil {
		log.Printf("[RESPONSE] Shuffle execution failed: %v", err)
		e.store.SaveResponseState(corr.ID, "FAILED")
		
		failEvent := models.OCSFFinding{
			EventID:      uuid.New().String(),
			CorrelationID: corr.ID,
			ActivityName: "Response Execution Failed",
			Severity:     "High",
			Message:      fmt.Sprintf("Shuffle execution failed for action %s: %v", result.Decision, err),
		}
		e.store.PublishEvent(context.Background(), failEvent)
		return
	}

	// Now independently verify the execution (Closed Loop)
	log.Printf("[RESPONSE] Polling for action %s verification (ExecID: %s)...", result.Decision, execID)
	
	verifyCtx, vCancel := context.WithTimeout(context.Background(), 10*time.Second) // 10 second bounded polling window
	defer vCancel()
	
	verifyResult, vErr := e.shuffleClient.PollExecutionStatus(verifyCtx, execID, 2*time.Second) // Poll every 2 seconds
	var finalState string
	var verifMsg string
	
	if vErr != nil && verifyResult.Status == "TIMEOUT" {
		finalState = "TIMEOUT"
		verifMsg = "Execution timed out during verification"
	} else if vErr != nil {
		finalState = "VERIFICATION_FAILED"
		verifMsg = fmt.Sprintf("Failed to verify execution status: %v", vErr)
	} else if verifyResult.Status == "SUCCEEDED" {
		finalState = "SUCCEEDED"
		verifMsg = fmt.Sprintf("Action %s verified successful via Shuffle API", result.Decision)
	} else {
		finalState = "FAILED"
		verifMsg = fmt.Sprintf("Execution reported failure: %s", verifyResult.Message)
	}

	e.store.SaveResponseState(corr.ID, finalState)
	
	// Generate Verification OCSF Event
	verifEvent := models.OCSFFinding{
		EventID:      uuid.New().String(),
		CorrelationID: corr.ID,
		ActivityName: "Response Verification",
		Severity:     "Info",
		Status:       finalState,
		Message:      verifMsg,
		Time:         time.Now(),
	}
	verifEvent.Observables = append(verifEvent.Observables, models.Observable{Type: "ExecutionID", Value: execID})
	e.store.PublishEvent(context.Background(), verifEvent)
}
