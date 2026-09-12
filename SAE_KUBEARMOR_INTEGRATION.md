# SAE KubeArmor Integration

## 1. Unified Event Model
KubeArmor logs are mapped seamlessly to the OCSF representation used by SAE. 
- **Severity Mapping**: KubeArmor severities (1-10) are mapped linearly to SAE severities (`Fatal`, `Critical`, `High`, `Medium`, `Low`, `Info`).
- **Telemetry Type**: KubeArmor alerts manifest as `Runtime Security Event`.

## 2. Integration Boundary
The boundary enforces strict separation between the KubeArmor agent and the SAE core engine.
- **Data Source**: Standard structured JSON telemetry (`kubearmor_alerts.json`).
- **Consumer**: An isolated `bufio.Reader` tail loop managed by a dedicated goroutine.
- **Parsing**: Schema parsing and missing-field failure logic live strictly within `internal/kubearmor`.

## 3. Policy and Action Semantics
KubeArmor uniquely provides native enforcement mechanisms (`Action`: `Block`, `Audit`).
- **SAE Autonomous Response**: SAE *ingests* KubeArmor's enforcement and detection telemetry. SAE does *not* natively manage KubeArmor's AppArmor/eBPF LSM policies or command its actions directly.
- **Semantic Definition**: "Block" actions are parsed as fatal runtime threats that were resolved at the edge by the KubeArmor daemon.

## 4. Failure Isolation
If KubeArmor halts, loops, or outputs corrupted JSON, the dedicated Go channel handles unmarshal failures by safely logging mapper errors. The unified SAE architecture guarantees that Wazuh, Zeek, Suricata, and Falco continue undisturbed.
