# SAE Garak License Audit

## 1. Upstream Identity
- **Project**: NVIDIA Garak (NVIDIA/garak)
- **License**: Apache License 2.0

## 2. SAE Distribution Posture
- **Integration Boundary**: Garak is consumed as an unmodified Python binary. The SAE application interfaces with it by parsing its generated `.report.jsonl` log files.
- **Legal Status**: Because Garak is Apache 2.0 licensed, consuming its output files carries zero legal risk. It poses no viral contamination risk to the Go backend logic.
