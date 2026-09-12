# SAE System Architecture

## Core Pipeline Flow

1. **Ingestion Layer (Go Goroutines)**
   - 13 distinct watchers monitor log files from Wazuh, Zeek, Suricata, Falco, KubeArmor, Trivy, ScoutSuite, Shuffle, TheHive, Cortex, Ollama, LangGraph, Garak.
   - Files are parsed and normalized into the Open Cybersecurity Schema Framework (`models.OCSFFinding`).
   
2. **Event Bus (Redis Streams)**
   - OCSF Events are enriched with UUID `EventID` and published to the `sae_events` Redis stream using `XADD`.
   - Replaces the legacy stdout-only termination model.

3. **Storage Layer (PostgreSQL)**
   - Events are persisted to `ocsf_events` with `ON CONFLICT DO NOTHING` for idempotent duplicate handling.
   - ClickHouse was determined to be **BLOCKED** due to download/runtime size constraints on the native host.

4. **Correlation Engine (Go)**
   - A consumer group (`sae_group`) continuously reads events via `XREADGROUP`.
   - Deterministic deterministic mapping extracts target vectors (`IP`, `User`, `Resource`).
   - Grouping events under shared `CorrelationID`s allows multi-tool fusion (e.g., Wazuh + Zeek + Suricata).

5. **AI Reasoning (LangGraph + Ollama)**
   - Triggered asynchronously when correlation threshold exceeds 3, or a single Critical event is observed.
   - Evaluates risk score, validates attack context, and recommends actions.

6. **Policy Engine & Response**
   - Implements hard safety boundaries. Rejects `block_ip` or `isolate_host` if suggested by the AI, requiring manual auth.
   - Forwards approved actions to the Response orchestrator (Shuffle - currently BLOCKED by Docker runtime constraints).

## Failure Resiliency
- **Redis Disconnects**: `engine.go` retries every 1s.
- **Postgres Disconnects**: `/ready` API reflects `503 Service Unavailable`.
- **Malformed Events**: Dropped gracefully at JSON unmarshal.
- **AI Outages**: Fails safe; logs error without breaking the ingestion pipeline.
