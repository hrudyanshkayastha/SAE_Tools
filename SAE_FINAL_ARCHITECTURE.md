# SAE Final Architecture Report

## Architecture Design
The core backend uses a unified data pipeline mapping disparate telemetry files into an open schema (OCSF), transmitting them via Redis Streams, persisting them to PostgreSQL, correlating targets natively in Go memory, extracting insight via a local LLM, checking constraints in a Deterministic Policy Layer, and finally attempting response orchestrator delegation.

```text
Log File (JSON) 
  --> [Watcher Goroutines]
    --> OCSFFinding Mapper 
      --> [Redis XADD sae_events]
        --> [Consumer Group sae_group]
          --> PostgreSQL sae_telemetry
          --> Correlation Engine 
            --> LangGraph + Ollama
              --> Policy Constraint Engine
                --> PostgreSQL decisions
                  --> REST API (Client Consumers)
```

## Resilience & Bounds
*   **Idempotency**: All database layers use UUID `EventID` mappings and enforce `ON CONFLICT DO NOTHING`.
*   **AI Sandboxing**: The `Policy Engine` physically traps `block_ip` or `isolate_host` actions, emitting `DANGER: Requiring human authorization` rather than orchestrating blindly.
*   **Network Tolerance**: Re-polls Redis and PostgreSQL gracefully.

## Final Storage Design
*   **PostgreSQL**: Handles operational logic (Telemetry `sae_telemetry`, Correlated Incidents `correlations`, and Responses `decisions`).
*   **ClickHouse**: **RUNTIME BLOCKED**. It was functionally configured in code via `clickhouse-go`, but authentication and volume mapping within native Docker Compose proved too volatile without explicit provisioning scripts. PostgreSQL safely handles both data lakes.
