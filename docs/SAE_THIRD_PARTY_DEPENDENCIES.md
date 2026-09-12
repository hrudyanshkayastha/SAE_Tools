# SAE Third-Party Dependencies

This document maps SAE product capabilities to their underlying third-party implementations.

## 1. Product to Provider Mapping

| SAE Capability | Third-Party Provider / Integration | State |
| :--- | :--- | :--- |
| **Endpoint Intelligence** | Wazuh | CONNECTED |
| **Network Intelligence** | Zeek | CONNECTED |
| **Network Detection** | Suricata | CONNECTED |
| **Container Intelligence** | Falco | PARTIAL (Awaiting Kernel Modules) |
| **CloudGuard** | KubeArmor | PARTIAL (Awaiting eBPF) |
| **Vulnerability Intelligence** | Trivy | CONNECTED |
| **Cloud Assessment** | ScoutSuite | PARTIAL (Awaiting Cloud Keys) |
| **Response Orchestration** | Shuffle | BLOCKED (Orchestrator constraints) |
| **Case Management** | TheHive | BLOCKED (Cassandra/ES constraints) |
| **Threat Analysis** | Cortex | BLOCKED (ES constraints) |
| **Local AI Runtime** | Ollama | CONNECTED |
| **Investigation Orchestration** | LangGraph | CONNECTED |
| **AI Security Testing** | Garak | CONNECTED |

## 2. Infrastructure Layer

| Infrastructure Need | Third-Party Provider | State |
| :--- | :--- | :--- |
| **Event Fabric** | Redis Streams | CONNECTED |
| **Operational State (Lake)** | PostgreSQL | CONNECTED |
| **Telemetry Lake** | ClickHouse | BLOCKED (Volume/Auth constraints) |

## 3. License Observance

All underlying providers run in standard configurations and emit data through custom Go adapters. We do NOT modify the upstream sources of the third-party providers. We simply consume their JSON output and normalize it into OCSF via the SAE Normalization Engine.
