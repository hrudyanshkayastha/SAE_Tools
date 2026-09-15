# SAE Trivy Integration

## 1. Unified Event Model
Trivy JSON output is mapped directly to the SAE Open Cybersecurity Schema Framework (OCSF) standard.
- **ActivityName**: Mapped as `Vulnerability Discovery`.
- **Severity Mapping**: Mapped explicitly to standard SAE severities (`CRITICAL` -> `Fatal`, `HIGH` -> `Critical`, `MEDIUM` -> `High`, `LOW` -> `Medium`, `UNKNOWN` -> `Info`).
- **Product Metadata**: Source mapped as `SAE Tool (Trivy Engine)`.

## 2. Integration Boundary
The integration strictly enforces standalone tool execution:
- **Execution Model**: Trivy runs as an independent executable. SAE does NOT natively link Trivy code, ensuring zero complex licensing crossover.
- **Data Channel**: Output is logged to `trivy/logs/trivy_report.json` using the standard Trivy JSON schema.
- **Consumer**: An isolated watcher polling the JSON file size and modified-time.

## 3. Failure Isolation
- The `consumer.go` safely catches unmarshal errors and logs them locally to `STDOUT` without interrupting the primary ingestion pipeline or impacting other running engines.

## 4. Unsupported Scenarios Handled
- If a vulnerability has an unknown severity, SAE defaults it to `Info`.
- If the JSON lacks a `VulnerabilityID`, it skips ingestion seamlessly.
