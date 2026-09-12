# SAE Final E2E Report

## Execution Profile
**Command**: `go run main.go`
**Infrastructure**: `docker-compose up -d` (Redis, PostgreSQL)
**Test Input Script**: `add_events.ps1` (Dynamically streams to `/logs/alerts/alerts.json`, `/logs/suricata/eve.json`, `/logs/zeek/conn.log`)

## The Flow (Verified by Output)
1. **Wazuh** `5710 SSH Brute Force` injected.
2. **Suricata** `ET SCAN Suspicious inbound` injected.
3. **Zeek** `conn.log` connection trace injected.
4. The watcher routines intercepted all three, normalized them, and pushed them to Redis Streams with distinct UUIDs.
5. The `sae_group` consumer read the stream and pushed telemetry entries to PostgreSQL.
6. The `correlate` engine detected IP `192.168.1.99` across all three OCSF events.
7. Context array reached critical threshold (`3+ events`) -> Triggered `LangGraph`.
8. Python execution passed via standard input payload -> Sent to `llama3.2:3b`.
9. Extracted risk payload: `{"risk_score": 10, "action": "LOG_AND_MONITOR", "validation": "Evidence validated"}`
10. System stored decision natively inside PostgreSQL `decisions` table.
11. Tested REST API using `curl.exe -s http://localhost:8080/incidents` and recovered correlation payloads.

## System Robustness
The system explicitly prevented the LLM from taking actions that could harm the operating system, while preserving the raw audit telemetry inside PostgreSQL (`SELECT raw FROM sae_telemetry`).
