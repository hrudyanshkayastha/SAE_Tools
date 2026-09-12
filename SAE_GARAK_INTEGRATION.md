# SAE Garak Integration

## 1. Unified Event Model
Garak operates as the **SAE AI Security Assessment** scanner. Its outputs are mapped into the Open Cybersecurity Schema Framework (OCSF) standard.
- **ActivityName**: Mapped as `AI Vulnerability Scan`.
- **Severity Mapping**: Calculated dynamically from the Garak `absolute_defcon` metrics. (`2: Critical`, `3: High`, `4: Medium`, `5: Low`). If the model perfectly bypasses the probe (`passed == total`), it emits as `Info`.
- **Product Metadata**: Source mapped as `SAE Tool (Garak Scanner)`.

## 2. Integration Boundary
The integration ensures decoupled operation:
- **Execution Model**: Garak is a batch process. The SAE backend operates as a JSONL log file watcher (similar to Zeek or Suricata), monitoring Garak's output directory for new `.report.jsonl` files.
- **Data Channel**: `garak_runs/*.report.jsonl` (Event Input) -> Go Mapper -> OCSF.

## 3. Failure Isolation
- The `mapper.go` adapter strictly ignores `start_run setup`, `init`, and `module_status` JSON records, silently proceeding until it discovers the terminal `eval` dictionary array.
