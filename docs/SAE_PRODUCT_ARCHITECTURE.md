# SAE Product Architecture

SAE represents a unified cybersecurity platform, composed of proprietary intelligence modules acting on standardized telemetry.

## The SAE Security Fabric

The product is built on a sequence of proprietary abstraction layers:

```text
SAE SENSORS & AGENTS (External / Third-Party Providers)
       ↓
SAE INGESTION
       ↓
SAE NORMALIZATION (OCSF)
       ↓
SAE EVENT FABRIC (Redis Backed)
       ↓
SAE DETECTION ENGINE
       ↓
SAE CORRELATION ENGINE (Deterministic IP/Asset Grouping)
       ↓
SAE THREAT & RISK INTELLIGENCE
       ↓
SAE INCIDENT MANAGEMENT
       ↓
LangGraph (LangGraph/Ollama Supported)
       ↓
SAE GUARDIAN
       ↓
SAE POLICY ENGINE (Safety Bounding)
       ↓
Shuffle
       ↓
SAE INDEPENDENT VERIFICATION
       ↓
SAE EVIDENCE & AUDIT (PostgreSQL Backed)
```

## Data Model (Canonical Entities)

SAE uses a standardized, product-first data model:
*   **Asset / Agent**: The endpoint or target in question.
*   **Event**: A normalized OCSF log (e.g., from Endpoint or Network Intelligence).
*   **Detection**: A prioritized signal.
*   **Threat / Risk**: A scored evaluation of an Incident.
*   **Incident**: An aggregated Correlation grouping.
*   **Policy / Guardian Run**: The safety evaluation against AI-generated actions.
*   **Response / Case**: The orchestration result.
*   **Evidence / Audit**: Immutable logs stored in the operational data lake.

Third-party tools generate or enrich these entities, but they do not become the primary product identity. SAE is the operational platform above the raw tooling layer.
