# SAE Cortex Integration

## 1. Unified Event Model
Cortex job reports are mapped directly into the SAE Open Cybersecurity Schema Framework (OCSF) standard.
- **ActivityName**: Mapped as `Threat Analysis Result`.
- **Severity Mapping**: Extracted from the Cortex execution status (`Success` -> `Info`, `Failure` -> `High`).
- **Product Metadata**: Source mapped as `SAE Tool (Cortex Engine)`.

## 2. Integration Boundary
The integration ensures decoupled operation and strict AGPL-3.0 compliance:
- **Execution Model**: Cortex operates as an independent Dockerized application backend. SAE does NOT natively link Cortex (AGPL-3.0) Scala code, retaining strong architecture modularity and license isolation.
- **Data Channel**: Cortex job reports (JSON payloads) are consumed from `cortex/logs/cortex_results.json`.
- **Consumer**: An isolated watcher polling the JSON file size and modified-time.

## 3. Failure Isolation
- The `consumer.go` safely traps unmarshal errors and logs them, ensuring that if Cortex emits a malformed result payload (e.g., during a crashed analyzer worker), the rest of the SAE ingestion daemon remains highly available.
