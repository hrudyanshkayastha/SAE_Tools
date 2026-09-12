# SAE Wazuh License Audit

## Policy Review
- **Proprietary Code:** `SAE/backend` (Go)
- **Third-Party Code:** `SAE/engine` (GPLv2 Wazuh Source)

## Audit Objectives
1. **Verify Licensing Attribution Files:** `SAE/licenses/WAZUH_GPLv2.txt` exists and matches the upstream `LICENSE` file.
2. **Verify Copyright Notices Intact:** After the forensic Git restoration, all original Wazuh Inc. copyright notices, GPL headers, and authorship markers inside the C/C++ files are fully intact. The unsafe find-and-replace algorithm that scrubbed these identifiers has been reverted.
3. **Verify No GPLv2 Source is Mixed into Proprietary Components:** The `SAE/backend` directory relies strictly on decoupled filesystem tailing (`consumer.go`) and standard HTTP calls (`client.go`). No cgo linking occurred.
4. **Verify No Unsupported Claims:** SAE branding applies only to the Go backend wrappers. The engine itself correctly retains its identity as Wazuh, complying with open-source reuse standards.

## Findings
**Status: AUDIT PASSED**

The integration respects the structural boundaries required by the GPLv2 license. All source-level attributions have been forensically restored to compliance.
