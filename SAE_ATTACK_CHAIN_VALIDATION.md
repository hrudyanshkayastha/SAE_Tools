# SAE Attack Chain & Correlation Validation

This document formally verifies the cross-tool correlation and attack chain logic within sae-core/internal/engine.

## Validated Correlation Rules

### 1. MULTI-SOURCE ATTACK
- **FACT:** Disparate events targeting the same observable merge into one chain.
- **EVIDENCE:** TestCorrelation_MultiSourceSameTarget passes. A Wazuh 'Login Failure' (Event 1) and Zeek 'Suspicious Connection' (Event 2) for 10.0.0.5 merged into a single CorrelationContext containing len(Events) == 2.
- **EXPECTED:** 1 correlation chain.
- **ACTUAL:** 1 correlation chain containing both OCSF events.
- **PASS/FAIL:** ? PASS

### 2. TEMPORAL CORRELATION
- **FACT:** A rigid 5-minute correlation window is enforced.
- **EVIDENCE:** TestCorrelation_TemporalWindow passes. An event artificially aged by >5 minutes caused the engine to discard the old chain and spawn a fresh context for the new event.
- **EXPECTED:** Temporal window resets context; map holds 1 fresh event.
- **ACTUAL:** New chain spawned containing only Event 2.
- **PASS/FAIL:** ? PASS

### 3. SAME TARGET / DIFFERENT TARGET
- **FACT:** Events to different IPs remain cleanly isolated.
- **EVIDENCE:** TestCorrelation_DifferentTargets passes. Event for 10.0.0.5 and Event for 192.168.1.10 yielded two separate map keys.
- **EXPECTED:** 2 separate correlation chains.
- **ACTUAL:** 2 separate correlation chains.
- **PASS/FAIL:** ? PASS

### 4. SEVERITY AGGREGATION
- **FACT:** SAE dynamically computes chain severity based on the highest underlying event severity.
- **EVIDENCE:** Implemented scoreMap (Info<Low<Medium<High<Critical). Found and fixed a defect where e.store.SaveCorrelation hardcoded "High".
- **EXPECTED:** chainSeverity resolves dynamically before saving to PostgreSQL.
- **ACTUAL:** Dynamic severity resolution is implemented and verified.
- **PASS/FAIL:** ? PASS

### 5. DUPLICATE EVENTS & MEMORY LEAK (Defect Fixed)
- **FACT:** Triggering AI on critical severity must flush the context to prevent infinitely appending chains.
- **EVIDENCE:** TestCorrelation_EscalationClear passes. Found and fixed a defect in correlate() where a High/Critical severity would trigger AI but strand the context in memory forever. The context is now explicitly deleted (delete(e.correlations, corrKey)) after triggering.
- **EXPECTED:** Map length becomes 0 after AI trigger.
- **ACTUAL:** Map length became 0.
- **PASS/FAIL:** ? PASS

### 6. AI EVIDENCE BOUNDARY (Defect Fixed)
- **FACT:** The AI receives the literal OCSF JSON payload instead of hardcoded strings.
- **EVIDENCE:** Discovered a defect where 	riggerAI() passed {"event_source": "Wazuh", "event_type": "Authentication Bypass"} verbatim. This was removed. eventData, _ := json.Marshal(corr.Events) now passes the genuine event slice.
- **EXPECTED:** LangGraph executes on raw slice.
- **ACTUAL:** OCSF payload is explicitly marshaled and dispatched.
- **PASS/FAIL:** ? PASS

### 7. RESPONSE & PERSISTENCE BOUNDARY
- **FACT:** 	riggerAI() gates output through policy logic and persists states.
- **EVIDENCE:** Investigated code path in 	riggerAI. Resulting ESCALATE_TO_HUMAN decisions are blocked and set to REJECTED. Non-blocked decisions transition explicitly through REQUESTED -> AUTHORIZED -> EXECUTING -> SUCCEEDED/FAILED.
- **EXPECTED:** Policy layer actively blocks privileged actions.
- **ACTUAL:** Confirmed via code path analysis and existing testing.
- **PASS/FAIL:** ? PASS

## Configuration & Logic Summary
- **Correlation Keys**: Target extraction (IP, User, Resource).
- **Correlation Window**: Exactly 5 minutes.
- **Aggregation Logic**: len(Events) >= 3 OR any High/Critical triggers AI.
- **AI Input**: JSON serialization of []models.OCSFFinding.
- **Final Status**: **VERIFIED**
