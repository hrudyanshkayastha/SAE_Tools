# SAE Integration Architecture

This document describes how external and third-party technologies interface with the SAE Security Fabric.

## Core Integration Boundary

The SAE integration boundary ensures that third-party tools are treated strictly as data providers or delegated execution environments, while SAE maintains the core logic, state, and policy control.

### 1. Ingestion Boundary
Third-party sensors (e.g., Wazuh, Suricata, Zeek) dump native JSON logs into designated directories or APIs. The SAE Ingestion adapters map this proprietary output strictly into standard **OCSF (Open Cybersecurity Schema Framework)** before allowing it into the Event Fabric. 
*   **Result:** The core correlation engine never interacts with tool-specific data.

### 2. Event Fabric Boundary
The SAE Event Fabric is powered by **Redis Streams**. This provides a resilient, decoupled transport mechanism guaranteeing delivery to the SAE storage and correlation nodes.

### 3. AI Reasoning Boundary
SAE utilizes **LangGraph** for orchestration and **Ollama** for local inference. 
*   **Safety Boundary:** The LLM output is strictly treated as *untrusted context*. The SAE Policy Engine intercepts all decisions (e.g., `block_ip`) and evaluates them against human-authorization limits before delegating to the SAE Response Orchestrator.

### 4. Storage Boundary
SAE uses **PostgreSQL** as its Evidence and Operational lake. The database does not store tool-specific tables, but rather SAE standard entities (`sae_telemetry`, `correlations`, `decisions`).

## Technical Matrix

| Capability | Associated Third-Party Provider |
| :--- | :--- |
| Endpoint Intelligence | Wazuh |
| Network Intelligence | Zeek, Suricata |
| Container Security | Falco, KubeArmor |
| Cloud Security | ScoutSuite |
| Vulnerability Intel | Trivy |
| Response | Shuffle |
| Case Management | TheHive |
| Threat Analysis | Cortex |
| Local Execution | Ollama |
| AI Reasoning Workflow | LangGraph |
| GenAI Testing | Garak |
