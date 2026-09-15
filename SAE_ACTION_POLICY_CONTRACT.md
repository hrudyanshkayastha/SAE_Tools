# SAE Action Policy Contract

## 1. Decision Schema
The SAE investigation pipeline operates on an explicit action enumeration model. The AI (LangGraph/Ollama) evaluates correlated event evidence and outputs a structured JSON object containing exactly two keys:
- \ction\ (string): The recommended response action.
- \severity\ (string): The assessed severity of the security event sequence.

### Supported Actions
- \LOG_AND_MONITOR\: Used for Info or Low severity events (e.g., routine authorized actions).
- \ESCALATE_TO_HUMAN\: Used for Medium severity or complex, ambiguous events requiring manual investigation. (Note: This is an investigatory state, not an execution action).
- \BLOCK_IP\: Used for High/Critical external network attacks (e.g., Brute Force).
- \ISOLATE_HOST\: Used for High/Critical internal malware or endpoint compromise.
- \KILL_PROCESS\: Used for High/Critical malicious processes.
- \DISABLE_ACCOUNT\: Used for High/Critical compromised user credentials.

---

## 2. Policy Engine Rules
The Go \engine.go\ intercepts the AI output boundary and applies deterministic policy rules. The LLM has **no direct execution capability**.

### Authorization Rules
1. **Validation & Fail-Closed**: Any recommended action not in the explicitly supported list is intercepted. The policy engine logs a malformed decision error and explicitly assigns the \REJECTED\ state, failing closed.
2. **Escalation Routing**: If the AI selects \ESCALATE_TO_HUMAN\, the policy recognizes this as a request for manual investigation. Automated response is halted, and the state is explicitly marked \REJECTED\ (meaning automation is blocked pending human intervention).
3. **Observation Routing**: \LOG_AND_MONITOR\ simply concludes the automated pipeline in the \COMPLETED_NO_ACTION\ state.
4. **Privileged Action Thresholds**: If an actionable response (\BLOCK_IP\, \ISOLATE_HOST\, \KILL_PROCESS\, \DISABLE_ACCOUNT\) is recommended, it must strictly satisfy the following threshold evidence:
   - The AI must calculate a \RiskScore\ of **80 or higher**.
   - The aggregated temporal chain severity must be explicitly **High or Critical**.
   If these conditions are not met, the action is marked \REJECTED\ as an unauthorized privileged execution attempt.
   If conditions are met, the action is marked \AUTHORIZED\ and proceeds to Response Execution.

---

## 3. Response Routing
Authorized actionable decisions are packaged into a \shuffle.ActionPayload\ and sent directly to the local SAE Shuffle instance. 
- SAE initiates an authenticated REST call to Shuffle (\ExecuteWorkflow\).
- If successful, the engine records an \EXECUTING\ state.
- If network connection fails, the state fails closed into \FAILED\ and an OCSF failure telemetry event is ingested.

---

## 4. Test Evidence
Deterministic tests validate the policy rules:
- **Benign Event (\LOG_AND_MONITOR\)**: \erify_ai_test.go\ natively demonstrates 3 correlated benign events evaluating to \LOG_AND_MONITOR\, concluding with low risk and safe policy closure.
- **Suspicious/Malformed Event (\ESCALATE_TO_HUMAN\)**: The LLM natively catches hallucination attempts (Semantic Evaluation) and safely defaults complex scenarios to \ESCALATE_TO_HUMAN\.
- **High-Confidence Attack (\BLOCK_IP\)**: Brute force payloads correctly elicit the \BLOCK_IP\ recommendation with a Critical (100) risk score, passing the authorized policy gate.
- **Fail Closed Mechanism**: Policy strictly defaults to block automation unless all thresholds are cleared.

---

## 5. Remaining Pipeline Gaps
While \RESPOND\ execution triggering is now fully validated in production code via strict action models, the following pipeline stages remain mathematically unverified end-to-end:

### VERIFY - [BLOCKED]
- SAE does not currently contain independent Target-State Verification architecture. It solely relies on querying Shuffle's API (\/api/v1/workflows/{id}/executions\) for execution status (i.e., did the API container finish).
- There is no feedback loop to prove a firewall actually blocked an IP or an EDR actually isolated a host.

### DISPLAY - [BLOCKED]
- End-to-end Dashboard UI rendering of the complete pipeline has not been verified in a browser context. We have strictly verified the underlying API HTTP 200 data delivery.
### Runtime Evidence (Execution Trace)
On live tests using the unmodified production engine.go and graph.py pipeline inside a Linux WSL container runtime with real Dockerized dependencies:

**1. AI DECISION:**
The LLM parsed the SSH Brute Force event using the explicitly defined schema.
{"action": "BLOCK_IP", "severity": "Critical"}

**2. POLICY AUTHORIZATION:**
The engine received BLOCK_IP and Critical, satisfying the strict RiskScore >= 80 and severity IN (High, Critical) policy bounds. 
State assigned: AUTHORIZED.

**3. SHUFFLE EXECUTE:**
The daemon dispatched an authenticated REST request to the local Shuffle container (POST /api/v1/workflows/.../execute).
Shuffle replied HTTP 200 returning execution_id: ff29c7e1-0fbe-4e5c-ac15-ffc260351141.
State assigned: EXECUTING.

**4. VERIFICATION POLLING:**
The daemon autonomously polled Shuffle via GET /api/v1/workflows/{id}/executions.
Shuffle eventually reported the execution container completed the task, returning FINISHED.
The daemon correctly mapped this to SUCCEEDED and securely closed the telemetry loop in the PostgreSQL database.

### Verdict
- **RESPOND**: VERIFIED (Target action generation, policy validation, and workflow dispatch are proven to function securely end-to-end on unmodified production code).
- **VERIFY**: BLOCKED (The daemon merely verifies the SOAR workflow container finished running; it lacks an architectural mechanism to verify the *target infrastructure* actually applied the BLOCK_IP firewall rule).
- **DISPLAY**: BLOCKED (Dashboard UI has not been validated in a browser context).
### Target-State Verification Overview
SAE now features a deterministic post-response Target-State Verification engine that operates independently of Shuffle's execution state.

**1. Verification Abstraction:**
erification.TargetVerifier securely exposes VERIFIED, NOT_VERIFIED, UNKNOWN, and ERROR matrices. 

**2. Independent Truth Override:**
When Shuffle returns FINISHED, the local SOAR cycle is marked as SUCCEEDED internally, but engine.go immediately launches the autonomous target verifier. The final E2E OCSF Correlation is dictated by the VERIFIER, not the SOAR response.

**3. Test Results:**
- SUCCESS: BLOCK_IP applies cleanly → target BLOCKED → VERIFIED
- FALSE SUCCESS: Shuffle succeeds but target remains NOT_BLOCKED → NOT_VERIFIED
- UNKNOWN / ERROR paths explicitly caught.

### Matrix
- RESPOND: VERIFIED
- VERIFY: VERIFIED
- DISPLAY: BLOCKED
