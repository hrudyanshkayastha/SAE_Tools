# SAE System API Documentation

The SAE unified core platform exposes a RESTful interface on `:8080` for platform health monitoring and event retrieval.

## 1. Platform Health
### `GET /health`
Returns the status of the HTTP server.
- **200 OK**: `{"status":"ok"}`

### `GET /ready`
Returns the operational readiness of the platform, testing upstream dependencies like PostgreSQL.
- **200 OK**: `{"status":"ready"}`
- **503 Service Unavailable**: `{"status":"postgres unavailable"}` (Triggered if DB connection drops).

## 2. Core Entities
### `GET /events`
Retrieves the most recent 50 OCSF events natively correlated and persisted in PostgreSQL.
- **Response**: `[ { "event_id": "...", "correlation_id": "...", "activity_name": "...", "severity": "...", "message": "..." } ]`

### `GET /incidents`
Retrieve correlated incidents grouped by Target or Attack Vector.
- *Status*: Implemented API shell, underlying DB schema mapping pending.

### `GET /investigations`
Retrieve AI-generated investigation timelines and hypotheses.
- *Status*: Implemented API shell, underlying DB schema mapping pending.

### `GET /decisions`
Retrieve policy decisions and response actions.
- *Status*: Implemented API shell, underlying DB schema mapping pending.
