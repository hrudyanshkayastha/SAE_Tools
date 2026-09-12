# SAE Final Dependency Report

## Go Modules
- `github.com/google/uuid`: OCSF UID Generation
- `github.com/redis/go-redis/v9`: Event Bus Publish/Subscribe
- `github.com/lib/pq`: Operational Data Lake & State Engine
- `github.com/ClickHouse/clickhouse-go/v2`: Telemetry Lake Driver (Present, but backend runtime constraints blocked integration natively).

## Runtimes
- **Docker Compose (Windows natively)**
  - Container: Postgres (Port 5432)
  - Container: Redis (Port 6379)
  - Container: ClickHouse (Failed network volume bind due to OS virtualization bugs; disabled gracefully).
- **Ollama**: (Port 11434). Model: `llama3.2:3b`.
- **Python (Native Windows Path)**: Used by LangGraph for StateGraph resolution.
