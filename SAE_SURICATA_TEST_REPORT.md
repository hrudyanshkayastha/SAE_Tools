# SAE Suricata Adapter Test Report

## 1. Test Scope
The `sae-core/internal/suricata` package normalizes raw EVE JSON output from the C/Rust Suricata engine. Tests assert that valid events map to OCSF correctly and malformed data fails safely.

## 2. Regression Gate
Before writing new tests, existing coverage for Wazuh and Zeek was run:
- Wazuh Engine Tests: 4/4 PASS
- Zeek Engine Tests: 4/4 PASS

## 3. Suricata Unit Test Results
Command: `go test -v ./...` in the `backend` directory.

| Target | Test Name | Input Description | Expected Behavior | Status |
|---|---|---|---|---|
| Suricata Engine | `TestMapEVEToOCSF_Alert` | Valid EVE alert | Maps to Intrusion Detection | **PASS** |
| Suricata Engine | `TestMapEVEToOCSF_Flow` | Valid EVE flow | Maps to Network Activity | **PASS** |
| Suricata Engine | `TestMapEVEToOCSF_DNS` | Valid EVE dns | Maps to DNS Activity | **PASS** |
| Suricata Engine | `TestMapEVEToOCSF_HTTP` | Valid EVE http | Maps to HTTP Activity | **PASS** |
| Suricata Engine | `TestMapEVEToOCSF_TLS` | Valid EVE tls | Maps to TLS Activity | **PASS** |
| Suricata Engine | `TestMapEVEToOCSF_Malformed` | Broken EVE JSON | Safe rejection, no crash | **PASS** |
| Suricata Engine | `TestMapEVEToOCSF_MissingFields`| Partial EVE data | Partial safe mapping | **PASS** |

## 4. Findings
- Suricata's `event_type` successfully drives the unified OCSF Activity category.
- Negative tests proved that malformed EVE lines do not panic the daemon.
- Zero regressions were introduced into the Wazuh or Zeek ingestion pathways.
