# SAE TheHive Integration

## 1. Unified Event Model
TheHive case/webhook telemetry is mapped into the SAE Open Cybersecurity Schema Framework (OCSF) standard.
- **ActivityName**: Mapped as `Case Management`.
- **Severity Mapping**: Extracted from TheHive's `severity` integer (`1` -> `Low`, `2` -> `Medium`, `3` -> `High`, `4` -> `Critical`).
- **Product Metadata**: Source mapped as `SAE Tool (TheHive Engine)`.

## 2. Integration Boundary
The integration ensures decoupled operation and strict AGPL-3.0 compliance:
- **Execution Model**: TheHive operates as an independent Dockerized application backend. SAE does NOT natively link TheHive (AGPL-3.0) code, retaining strong architecture modularity and license isolation.
- **Data Channel**: TheHive case payloads (via webhook JSON) are consumed from `thehive/logs/thehive_results.json`.
- **Consumer**: An isolated watcher polling the JSON file size and modified-time.

## 3. Failure Isolation
- The `consumer.go` safely traps unmarshal errors and logs them, ensuring that if TheHive emits a malformed result payload, the rest of the ingestion daemon remains highly available.
