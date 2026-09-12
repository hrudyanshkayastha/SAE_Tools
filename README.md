# SAE (Security AI Engine)

**Status: MVP / Prototype**

SAE is a privacy-first security telemetry, correlation, and AI-assisted investigation prototype. It demonstrates the architectural feasibility of combining open-source detection tools with a local LLM via an OCSF-standardized event fabric.

## Architecture
- **Sensors**: Wazuh, Zeek, Suricata, Trivy, Garak
- **Fabric**: OCSF, Redis, PostgreSQL
- **AI Triage**: LangGraph, Ollama (Llama3.2 3B)

## Limitations and Constraints (Not Production Ready)
- **NOT an Enterprise XDR/SOAR/CNAPP replacement**: The platform cannot physically execute response actions autonomously in current compute environments.
- **NO Zero-Day Prediction**: Relies explicitly on signature and rule-based detection from integrated sensors.
- **NO Autonomous SOC**: SAE acts as an AI-assisted triage engine, but human analysts remain required.
- **NO Native UEBA**: Does not calculate behavioral baselines.
- **Cloud/Container**: Falco, KubeArmor, and ScoutSuite integrations are currently blocked from runtime execution.

For a full breakdown of the integrations, blockers, and technical validation, see `SAE_FINAL_VERDICT.md` and `SAE_FINAL_AUDIT.md`.
