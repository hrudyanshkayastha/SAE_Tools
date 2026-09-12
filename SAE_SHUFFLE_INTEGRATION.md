# SAE Shuffle Integration

## 1. Unified Event Model
Shuffle execution telemetry is mapped into the SAE Open Cybersecurity Schema Framework (OCSF) standard.
- **ActivityName**: Mapped as `Response Execution`.
- **Severity Mapping**: Extracted from the SOAR execution state (`SUCCESS` -> `Info`, `FAILED`/`ERROR` -> `High`, `ABORTED` -> `Medium`).
- **Product Metadata**: Source mapped as `SAE Tool (Shuffle Engine)`.

## 2. Integration Boundary
The integration ensures decoupled operation:
- **Execution Model**: Shuffle operates as an independent Dockerized orchestration SOAR stack. SAE does NOT natively link Shuffle (MIT) code, retaining strong architecture modularity.
- **Data Channel**: Shuffle workflow outputs are consumed from `shuffle/logs/shuffle_results.json`.
- **Consumer**: An isolated watcher polling the JSON file size and modified-time.

## 3. Failure Isolation
- The `consumer.go` safely traps unmarshal errors and logs them, ensuring that if a SOAR response action yields a malformed result payload, the rest of the ingestion daemon remains highly available.
