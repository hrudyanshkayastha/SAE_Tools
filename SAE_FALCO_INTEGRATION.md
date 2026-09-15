# SAE Falco Integration Architecture

## 1. Architectural Role
Falco is integrated into the SAE product ecosystem under the capability:
**SAE Containers Security**

It is designed to run loosely coupled to the SAE core, avoiding unsafe C++ bindings or blocking logic. 

**Conceptual Data Flow**:
1. Falco Engine (userspace / kernel) captures a runtime security violation.
2. Falco formats the event according to its ruleset and writes structured telemetry to `falco_events.json`.
3. The SAE Go Daemon reads the stream concurrently alongside Zeek, Wazuh, and Suricata.
4. The Falco Adapter (`sae-core/internal/falco`) decodes the JSON, enforces presence of required fields (`rule`, `output`), and parses out security context (`container.id`, `proc.cmdline`, `user.name`).
5. The Adapter maps the Falco priority (e.g. `Emergency`, `Critical`, `Notice`) to SAE's unified Severity models (e.g. `Fatal`, `Critical`, `Low`).
6. The event emerges as a fully normalized OCSF finding inside the unified detection pipeline.

## 2. Event Model & Priorities
Falco's unstructured `output` string and structured `output_fields` dictionary are selectively mapped to OCSF Observables to enable relational correlation:

* **Container Context**: Extracted from `container.id`.
* **Process Context**: Extracted from `proc.cmdline`.
* **User Context**: Extracted from `user.name`.
* **Network/File Context**: Extracted from `fd.name`.

**Priority Mapping Table**:
| Falco Priority | SAE Unified Severity |
|----------------|-----------------------|
| Emergency      | Fatal                 |
| Alert          | Critical              |
| Critical       | Critical              |
| Error          | High                  |
| Warning        | Medium                |
| Notice         | Low                   |
| Informational  | Info                  |
| Debug          | Info                  |

## 3. Failure Isolation
If Falco writes malformed JSON, fails to start, or rotates its log improperly, the `bufio.Scanner` within the SAE consumer will merely pause and re-attempt. The consumer runs within an independent Goroutine and cannot crash or block the ingestion streams for Wazuh, Zeek, or Suricata. All telemetry remains robustly parallelized.
