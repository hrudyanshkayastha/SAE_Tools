package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sae-core/internal/cortex"
	"sae-core/internal/langgraph"
	"sae-core/internal/shuffle"
	"sae-core/internal/storage"
	"sae-core/internal/thehive"
	"sae-core/internal/verification"
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
	verifier      verification.TargetVerifier

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

func NewEngine(store *storage.Storage, shuffleWebhook, thehiveURL, cortexURL string, verifier verification.TargetVerifier) *Engine {
	return &Engine{
		store:         store,
		shuffleClient: shuffle.NewClient(shuffleWebhook),
		thehiveClient: thehive.NewClient(thehiveURL, "mock-api-key"),
		cortexClient:  cortex.NewClient(cortexURL, "mock-api-key"),
		verifier:      verifier,
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

func (e *Engine) correlate(ctx context.Context, event models.OCSFFinding) *langgraph.GraphOutput {
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

	// Temporal Correlation: 5-minute window
	if exists && time.Since(ctxData.CreatedAt) > 5*time.Minute {
		// Window expired, force a new chain
		exists = false
	}

	if !exists {
		ctxData = &CorrelationContext{
			ID:        "CORR-" + target + "-" + uuid.New().String()[:8],
			Target:    target,
			CreatedAt: time.Now(),
		}
		e.correlations[corrKey] = ctxData
	}

	ctxData.Events = append(ctxData.Events, event)
	event.CorrelationID = ctxData.ID
	e.store.SaveTelemetry(event)

	// Severity aggregation: escalate if any event in chain is higher severity
	// Simple mapping: Info<Low<Medium<High<Critical
	chainSeverity := "Info"
	scoreMap := map[string]int{"Info": 0, "Low": 1, "Medium": 2, "High": 3, "Critical": 4}
	currentHighest := 0
	for _, ev := range ctxData.Events {
		if scoreMap[ev.Severity] > currentHighest {
			currentHighest = scoreMap[ev.Severity]
			chainSeverity = ev.Severity
		}
	}

	e.store.SaveCorrelation(ctxData.ID, ctxData.Target, len(ctxData.Events), chainSeverity, "ACTIVE")

	// Trigger criteria: High/Critical severity OR >= 3 correlated events
	if currentHighest >= 3 || len(ctxData.Events) >= 3 {
		res, _ := e.triggerAI(ctxData)
		delete(e.correlations, corrKey)
		return res
	}
	return nil
}

func (e *Engine) triggerAI(corr *CorrelationContext) (*langgraph.GraphOutput, error) {
	// AI Evidence Boundary: Send the REAL correlated events
	// Instead of hardcoded Wazuh strings, we pass the actual OCSF array
	eventData, _ := json.Marshal(corr.Events)

	result, err := langgraph.ExecuteReasoningGraph(eventData)
	if err != nil {
		log.Printf("[AI REASONING] LangGraph execution failed: %v", err)
		return nil, err
	}

	// Save initial Decision to Postgres as REQUESTED
	e.store.SaveDecision(corr.ID, result.RiskScore, result.Validation, result.Decision)
	// Calculate chain severity for resolution record
	chainSeverity := "Info"
	scoreMap := map[string]int{"Info": 0, "Low": 1, "Medium": 2, "High": 3, "Critical": 4}
	currentHighest := 0
	for _, ev := range corr.Events {
		if scoreMap[ev.Severity] > currentHighest {
			currentHighest = scoreMap[ev.Severity]
			chainSeverity = ev.Severity
		}
	}

	// Update correlation to RESOLVED
	e.store.SaveCorrelation(corr.ID, corr.Target, len(corr.Events), chainSeverity, "RESOLVED")
	e.store.SaveResponseState(corr.ID, "REQUESTED")

	// Policy Engine bounds checking
	validActions := map[string]bool{
		"LOG_AND_MONITOR":   true,
		"ESCALATE_TO_HUMAN": true,
		"BLOCK_IP":          true,
		"ISOLATE_HOST":      true,
		"KILL_PROCESS":      true,
		"DISABLE_ACCOUNT":   true,
	}

	if !validActions[result.Decision] {
		log.Printf("[POLICY ENGINE] DANGER: Unrecognized or malformed decision: %s. Failing closed.", result.Decision)
		e.store.SaveResponseState(corr.ID, "REJECTED")
		return result, nil
	}

	if result.Decision == "ESCALATE_TO_HUMAN" {
		log.Printf("[POLICY ENGINE] AI requested manual investigation. Halting automated response.")
		e.store.SaveResponseState(corr.ID, "REJECTED")
		return result, nil
	}

	if result.Decision == "LOG_AND_MONITOR" {
		e.store.SaveResponseState(corr.ID, "COMPLETED_NO_ACTION")
		return result, nil
	}

	// For actionable decisions, check strict policy constraints
	// Rule: Action is authorized ONLY if RiskScore >= 80 AND chainSeverity is High/Critical
	if result.RiskScore < 80 || (chainSeverity != "High" && chainSeverity != "Critical") {
		log.Printf("[POLICY ENGINE] DANGER: Unauthorized action %s recommended for low-confidence/low-severity evidence. Blocked.", result.Decision)
		e.store.SaveResponseState(corr.ID, "REJECTED")
		return result, nil
	}

	e.store.SaveResponseState(corr.ID, "AUTHORIZED")

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
			EventID:       uuid.New().String(),
			CorrelationID: corr.ID,
			ActivityName:  "Response Execution Failed",
			Severity:      "High",
			Message:       fmt.Sprintf("Shuffle execution failed for action %s: %v", result.Decision, err),
		}
		e.store.PublishEvent(context.Background(), failEvent)
		return result, err
	}

	// Now independently verify the execution (Closed Loop)
	log.Printf("[RESPONSE] Polling for action %s verification (ExecID: %s)...", result.Decision, execID)

	verifyCtx, vCancel := context.WithTimeout(context.Background(), 10*time.Second) // 10 second bounded polling window
	defer vCancel()

	verifyResult, vErr := e.shuffleClient.PollExecutionStatus(verifyCtx, execID, 2*time.Second) // Poll every 2 seconds
	var finalState string = "SUCCEEDED"
	var verifMsg string
	if vErr != nil {
		verifMsg = fmt.Sprintf("Shuffle execution failed for action %s: %s", result.Decision, vErr.Error())
		log.Printf("[RESPONSE] %s", verifMsg)
		finalState = "FAILED"
	} else if verifyResult.Status == "SUCCEEDED" {
		verifMsg = fmt.Sprintf("Shuffle successfully verified execution of %s (ExecID: %s)", result.Decision, verifyResult.ExecutionID)
		log.Printf("[RESPONSE] %s", verifMsg)
		finalState = "SUCCEEDED"
	} else {
		verifMsg = fmt.Sprintf("Shuffle verification failed/timed-out for %s (Status: %s, Message: %s)", result.Decision, verifyResult.Status, verifyResult.Message)
		log.Printf("[RESPONSE] %s", verifMsg)
		finalState = verifyResult.Status
	}

	e.store.SaveResponseState(corr.ID, finalState)

	// Generate Base Verification OCSF Event
	verifEvent := models.OCSFFinding{
		EventID:       uuid.New().String(),
		CorrelationID: corr.ID,
		ActivityName:  "Response Verification",
		Severity:      "Info",
		Status:        finalState,
		Message:       verifMsg,
		Time:          time.Now(),
	}
	verifEvent.Observables = append(verifEvent.Observables, models.Observable{Type: "ExecutionID", Value: execID})
	e.store.PublishEvent(context.Background(), verifEvent)

	// TARGET-STATE VERIFICATION
	if finalState == "SUCCEEDED" && e.verifier != nil {
		tsRes := e.verifier.Verify(context.Background(), result.Decision, corr.Target)

		metadata := map[string]string{
			"action_type":    result.Decision,
			"expected_state": tsRes.Expected,
			"observed_state": tsRes.Observed,
			"evidence":       tsRes.Evidence,
		}
		rawBytes, _ := json.Marshal(metadata)

		tsEvent := models.OCSFFinding{
			EventID:       uuid.New().String(),
			CorrelationID: corr.ID,
			ActivityName:  "Target-State Verification",
			Severity:      "Info",
			Status:        tsRes.State,
			Message:       fmt.Sprintf("Target-State Verification Result: %s. Expected: %s, Observed: %s.", tsRes.State, tsRes.Expected, tsRes.Observed),
			Time:          tsRes.Timestamp,
			RawEvent:      string(rawBytes),
		}
		tsEvent.Observables = append(tsEvent.Observables, models.Observable{Type: "Target", Value: tsRes.Target})
		tsEvent.Observables = append(tsEvent.Observables, models.Observable{Type: "ActionID", Value: execID})
		
		e.store.PublishEvent(context.Background(), tsEvent)

		// Set final correlation state based on the Target-State reality, not just Shuffle
		e.store.SaveResponseState(corr.ID, tsRes.State)
		if tsRes.State != verification.StateVerified {
			log.Printf("[VERIFY] Target-state verification failed for %s on %s: %s", result.Decision, corr.Target, tsRes.Evidence)
		} else {
			log.Printf("[VERIFY] Target-state VERIFIED for %s on %s", result.Decision, corr.Target)
		}
	}

	return result, nil
}
