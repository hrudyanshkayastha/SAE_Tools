# SAE ScoutSuite Integration

## 1. Unified Event Model
ScoutSuite JSON output is mapped directly into the SAE Open Cybersecurity Schema Framework (OCSF) standard.
- **ActivityName**: Mapped as `Cloud Assessment`.
- **Severity Mapping**: Mapped explicitly to standard SAE severities (`danger` -> `Critical`, `warning` -> `High`, `info` -> `Medium`).
- **Product Metadata**: Source explicitly mapped as `SAE Tool (ScoutSuite Engine)`.

## 2. Integration Boundary
The integration strictly enforces standalone tool execution:
- **Execution Model**: ScoutSuite runs as an independent executable via Python. SAE does NOT natively link ScoutSuite (GPLv2) code, ensuring zero complex licensing crossover.
- **Data Channel**: Output is logged to `scoutsuite/logs/scoutsuite_results.json` using the standard JSON schema.
- **Consumer**: An isolated watcher polling the JSON file size and modified-time.

## 3. Failure Isolation
- The `consumer.go` safely catches unmarshal errors and logs them locally without interrupting the primary ingestion pipeline or impacting other running engines.

## 4. Unsupported Scenarios Handled
- Malformed inputs, missing findings, and unknown severities are all bypassed deterministically without triggering panics.
