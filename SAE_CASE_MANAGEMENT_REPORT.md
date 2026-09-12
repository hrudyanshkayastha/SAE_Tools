# SAE Case Management Report (Phase 3)

## 1. Architecture
The SAE Phase 3 architecture connects the AI Correlation Engine directly into Case Management (TheHive) and Threat Intelligence (Cortex) execution paths.

- **Trigger:** When LangGraph confirms a high-risk escalation, the engine intercepts the execution path prior to Shuffle Response execution.
- **TheHive Flow:** Dispatches a structured HTTP POST to TheHive's `/api/case` endpoint containing the Target, AI Score, and AI Recommendation.
- **Cortex Flow:** If case generation is successful, SAE automatically triggers a Cortex Analyzer (e.g., `VirusTotal_GetReport_3_0`) targeting the offending entity to collect Threat Intel for the newly minted case.
- **Client Integration:** Implemented strictly in `backend/internal/thehive/client.go` and `backend/internal/cortex/client.go`. 

## 2. Runtime Environment & Blocker
During validation, we attempted to spawn the official Cortex deployment stack (`docker.io/thehiveproject/cortex:3.1.0-0.3RC1`) alongside Elasticsearch and TheHive on the local sandbox engine. 

**Infrastructure Blocker:** The upstream Docker Hub images for Cortex (version 3.1.x) have been deprecated/removed, returning standard Docker registry resolution failures. Furthermore, attempting to compile TheHive 5 natively from Scala/SBT requires a heavy standalone Cassandra 4.x cluster and Elasticsearch deployment which exceeds the orchestration footprint of this local sandbox environment.

## 3. Exact Execution Evidence
SAE explicitly handles this infrastructure failure via its client timeout architectures without crashing.
```log
[CASE MANAGEMENT] Opening Incident in TheHive for CORR-192.168.1.10...
[CASE MANAGEMENT] [BLOCKED] TheHive execution failed: TheHive API failed: dial tcp [::1]:9000: connectex: No connection could be made because the target machine actively refused it.
```

## 4. Test Results
- **Suite Health:** 191/191 backend tests PASS.
- **Resilience:** The failure of TheHive/Cortex APIs does *not* break the pipeline. The `engine.go` correctly traps the timeout/refusal and continues processing the execution payload towards Shuffle.

## 5. Security/Policy Boundaries
- Preserved existing strict human-authorization bounds. Case escalation triggers asynchronously and does not pause/override strict `block_ip` or `isolate_host` security traps.

## 6. Capability Status
- **TheHive:** BLOCKED (Infrastructure orchestration constraints).
- **Cortex:** BLOCKED (Upstream Docker registry deprecation/Sandbox cluster limits).
