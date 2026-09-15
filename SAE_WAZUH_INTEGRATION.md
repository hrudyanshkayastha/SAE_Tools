# SAE-Wazuh Integration Strategy

## 1. Overview
Wazuh (SAE Tool) is a powerful, GPLv2-licensed endpoint protection platform. To utilize its capabilities within the proprietary SAE Core without triggering viral licensing, we must integrate it via a strict adapter pattern over network protocols (API/Sockets) rather than direct code inclusion.

## 2. Components Inspected
- **Wazuh Agent:** C-based, monitors system logs, FIM, and rootkits.
- **Wazuh Manager (`ossec-analysisd`):** The monolithic rules engine.
- **Wazuh API:** Python/Node.js based REST API running on port 55000.
- **Alert Output:** JSON format output to `logs/alerts/alerts.json`.

## 3. SAE Integration Boundary
We have chosen the **ADAPTER/IPC INTEGRATION** method.
- **Data Plane (Telemetry):** We do not use Filebeat. Instead, the SAE Ingestion Daemon (Go) implements a file-tailing consumer (`internal/wazuh/consumer.go`) that mounts the Wazuh `alerts.json` directory. This provides zero-network-latency ingestion while maintaining process isolation.
- **Control Plane (Management):** The SAE Go Core uses an HTTP API Client (`internal/wazuh/client.go`) to authenticate via JWT to the Wazuh Manager, enabling SAE to health-check the engine and query agent status dynamically.

## 4. Unified Event Model Mapping
Wazuh's proprietary JSON alerts are intercepted by the `MapAlertToOCSF()` function. This abstracts away the Wazuh-specific taxonomy (e.g., `rule.level`) and normalizes it to the standard OCSF (Open Cybersecurity Schema Framework) format before passing it to the SAE Correlation Engine.
