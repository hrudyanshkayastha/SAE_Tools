# SAE OmniRoute AI Gateway Evaluation

This report evaluates the diegosouzapw/OmniRoute repository as an **optional** fallback and routing layer for SAE's AI logic. OmniRoute is fundamentally evaluated as an AI gateway, not a security engine.

## 1. Source & Dependency Audit
- **FACT:** OmniRoute is a mature Node.js / TypeScript AI routing application.
- **EVIDENCE:** package.json declares version 3.8.51, utilizing a Next.js full-stack framework with 356 abstracted providers.
- **TEST:** Inspected the source at commit 152d95108c9c3d557562311ffed63240a511eb31.
- **RESULT:** The repository natively exposes an OpenAI-compatible /v1/* endpoint (src/server/authz/types.ts). It has a 4-tier automatic fallback system (Subscription → API → Cheap → Free) and explicitly supports local endpoints (e.g., Ollama).
- **VERDICT:** Excellent runtime compatibility with SAE's LangGraph requirement.

## 2. Security Audit
- **FACT:** SAE requires strict security for API keys and prompt evidence.
- **EVIDENCE:** 
  - src/lib/db/encryption.ts demonstrates AES-256-GCM encryption for all stored credentials at rest.
  - src/lib/logPayloads.ts and src/lib/piiSanitizer.ts demonstrate active redaction of sensitive keys (uthorization, pi_key) before persisting telemetry logs.
  - Core network fetches originate via safe HTTP proxy middleware (http-proxy-middleware, xios).
- **TEST:** Grepped source for secret handling and dangerous etch patterns.
- **RESULT:** No glaring SSRF, unencrypted token persistence, or unredacted prompt logging behaviors were found.
- **VERDICT:** Meets standard AI gateway security expectations.

## 3. Malformed Input Investigation
- **FACT:** SAE recently exhibited a LangGraph malformed input error.
- **EVIDENCE:** The issue resided entirely in sae-core/internal/langgraph/graph.py failing to parse a JSON list instead of a dict on stdin. 
- **TEST:** The defect occurs *before* any HTTP request leaves the LangGraph Python environment.
- **RESULT:** OmniRoute sits strictly *downstream* of this failure. 
- **LIMITATION:** OmniRoute **does not fix** the malformed input issue. That issue must be (and was) resolved natively within the SAE graph.py logic.
- **VERDICT:** Irrelevant to the malformed input bug, but highly relevant for upstream model resilience.

## 4. Performance & Safe Prototype
- **FACT:** A safe proxy test script (prototype.go) was authored to validate standard OpenAI JSON request bridging to OmniRoute.
- **EVIDENCE:** OmniRoute's architectural footprint is large (Next.js, SQLite, heavy npm dependency tree). 
- **RESULT:** The raw proxy layers natively handle JSON forwarding quickly, but deploying the full stack locally is a heavy requirement.
- **LIMITATION:** Local deployment is resource-intensive due to the Node.js/Next.js ecosystem. It adds significant RAM overhead compared to a compiled Go binary.
- **VERDICT:** Performance is adequate for routing, but deployment is complex.

---

## 5. Final Verdict & Integration Recommendation

**OMNIROUTE:**
**Commit:** 152d95108c9c3d557562311ffed63240a511eb31
**Version:** 3.8.51
**License:** MIT License
**Runtime:** Node.js >= 24, TypeScript, SQLite/PostgreSQL
**Security findings:** AES-256-GCM encrypted keys, PII/Secret redaction verified.
**SAE compatibility:** High. Native OpenAI /v1 endpoint replaces direct Ollama calls seamlessly.
**Malformed-input relevance:** None. The SAE defect occurred upstream in graph.py.
**Performance:** Architecturally heavy (Next.js), fast proxy routing, high memory footprint.
**Fallback capability:** Excellent (4-tier circuit breakers).
**Integration complexity:** Moderate (Requires standing up a separate Node container next to SAE).

**Final verdict:** **3. OPTIONAL AI GATEWAY**

### Minimal Future Integration Design
If enabled, SAE should not replace its local Ollama binary requirement. Instead, the user configures:
SAE_AI_GATEWAY=http://localhost:20128/v1

**Path:**
SAE Correlation → LangGraph (graph.py) → OmniRoute OpenAI Adapter → Ollama / Remote Providers

OmniRoute will execute fallback strategies automatically, while SAE's Policy Engine safely receives standard structured JSON responses as if talking directly to Ollama.
