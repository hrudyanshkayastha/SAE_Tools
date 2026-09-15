# SAE Suricata Integration Architecture

## 1. Architectural Boundary
Suricata natively logs high-speed network insights to a structured JSON file called `eve.json`. Rather than wrapping Suricata within Go processes (which risks instability at gigabit speeds), SAE interfaces with Suricata strictly via its output sink. 
Suricata operates continuously in C/Rust, outputting `eve.json`, while the SAE Go Daemon tails this file concurrently with `alerts.json` (Wazuh) and `conn.log` (Zeek).

## 2. Event Model Mapping (OCSF)
The `sae-core/internal/suricata` package translates Suricata EVE events into the SAE Open Cybersecurity Schema Framework (OCSF) format. 

### Supported Suricata Event Types
| EVE Event Type | SAE Activity Name | OCSF Mapping Behavior |
|---|---|---|
| `alert` | Intrusion Detection | Maps `alert.signature` and `alert.category` into the main message. Maps Suricata 1-4 severities directly to SAE standard severities. |
| `flow` | Network Activity | Maps flow endpoints, duration, and protocol (Netflow equivalent). |
| `dns` | DNS Activity | Extracts `rrname` (queried domain) and `rrtype`. |
| `http` | HTTP Activity | Extracts Hostname, URL path, Method, and User-Agent. |
| `tls` | TLS Activity | Extracts SNI (Server Name Indication) and TLS version. |

## 3. Negative Handling
- Malformed JSON records are gracefully rejected by `MapEVEToOCSF`, returning safe errors without crashing the main SAE event loop.
- Unrecognized or generic event types (e.g., `smb`, `dcerpc` before specific adapters are written) default to `Network Telemetry` with an `Info` severity, ensuring no data is blindly dropped.

## 4. Unification
In `main.go`, Suricata is spawned as a concurrent goroutine alongside Wazuh and Zeek. All three streams merge into a single `OCSFFinding` struct definition, proving that a single relational DB/Clickhouse sink can handle all three natively.
