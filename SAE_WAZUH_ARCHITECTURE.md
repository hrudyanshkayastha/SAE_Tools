# SAE-Wazuh Architecture

## System Boundary

Wazuh is deployed as an isolated sidecar container alongside the SAE Core. 

```mermaid
graph TD
    %% External Endpoints
    E[Endpoints / Hosts] -->|OSSEC Protocol: 1514| WA[Wazuh Manager Sidecar]

    %% Wazuh Internals
    subgraph "Isolated Wazuh Engine (GPLv2)"
        WA -->|Rule Matching| AE[ossec-analysisd]
        AE -->|JSON Alerts| AL(alerts.json)
        API[Wazuh REST API: 55000]
    end

    %% SAE Core Internals
    subgraph "SAE Core Backend (Proprietary / Go Native)"
        WC[Wazuh API Client]
        TC[Telemetry Consumer]
        OM[OCSF Mapper]
        CE[SAE Correlation Engine]
        
        WC <-->|JWT Auth & Health| API
        TC -->|Tail/Socket Stream| AL
        TC -->|Raw JSON| OM
        OM -->|OCSF Model| CE
    end
```

## Security & Reliability Design
1. **Authentication:** SAE maintains an ephemeral JWT token via the `/security/user/authenticate` endpoint. Token rotation and expiration handling are built into the Go client.
2. **Timeouts:** The Go HTTP Client enforces a strict 10-second timeout to prevent the SAE Core from hanging if the Wazuh Manager fails to respond.
3. **Data Loss Prevention:** The telemetry consumer seeks to EOF and tails. In production, this will utilize an inode-tracking mechanism to ensure no alerts are lost during daemon restarts.
