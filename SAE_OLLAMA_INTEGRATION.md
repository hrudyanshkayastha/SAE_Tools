# SAE Ollama Integration

## 1. Unified Event Model
Ollama acts as the **Local AI Runtime** for SAE. It evaluates contextual evidence and produces structured JSON responses mapped into the Open Cybersecurity Schema Framework (OCSF) standard.
- **ActivityName**: Mapped as `AI Threat Investigation`.
- **Severity Mapping**: The LLM explicitly defines the severity (`Info`, `Low`, `Medium`, `High`, `Critical`), reverting to `Info` if it hallucinates or provides invalid output.
- **Product Metadata**: Source mapped as `SAE Tool (Ollama Runtime)`.

## 2. Integration Boundary
The integration ensures decoupled operation and strict mitigation of LLM hallucination risks:
- **Execution Model**: Ollama operates as a separate HTTP API daemon. It executes synchronously when queried by the SAE backend (`client.go`).
- **Data Channel**: SAE queries `/api/generate` with a carefully engineered prompt locking the LLM to output only JSON. 
- **Security Posture**: The LLM output is entirely isolated to *recommendation generation*. It possesses absolutely no authority to perform writes, trigger playbooks, or alter system state. Its output is merely ingested into the event bus for downstream policy evaluation.

## 3. Failure Isolation
- The Go adapter uses strong typing and failure fallbacks. If the LLM generates raw markdown, malformed JSON, or fails to adhere to the schema, the adapter safely degrades to treating the entire response as a raw text recommendation string while defaulting severity and risk scores to lowest-impact baselines.
