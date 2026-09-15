# SAE LangGraph Integration

## 1. Unified Event Model
LangGraph operates as the orchestrator of the **SAE Reasoning Pipeline**. The final output state of the graph is mapped into the Open Cybersecurity Schema Framework (OCSF) standard.
- **ActivityName**: Mapped as `SAE Reasoning Workflow`.
- **Severity Mapping**: Calculated dynamically from the graph's `RiskScore` state (`<30: Info`, `<50: Low`, `<80: Medium`, `>=80: High`).
- **Product Metadata**: Source mapped as `SAE Tool (LangGraph Engine)`.

## 2. Integration Boundary
The integration ensures decoupled operation and strict mitigation of LLM hallucination risks:
- **Execution Model**: The core engine is Go-based. To avoid creating a separate Python microservice daemon, the Go adapter uses `os/exec` to invoke `python graph.py` and pipes the event over `stdin`. This fulfills the requirement to "Follow the existing SAE language/runtime architecture rather than introducing an unnecessary second backend."
- **Data Channel**: `stdin` (Event Input) -> `stdout` (Structured JSON Result).
- **Security Posture**: The graph is explicitly structured to run evidence validation *after* the LLM node.

## 3. Failure Isolation
- The `client.go` adapter traps any errors, panics, or infinite loops caused by the python execution, preventing the SAE ingestion daemon from failing if the graph crashes.
