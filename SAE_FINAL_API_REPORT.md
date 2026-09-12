# SAE Final API Report

## Real Implementation (Not a Mock)
The HTTP server running on `:8080` interfaces natively with the `storage` package executing SQL queries against the PostgreSQL operational data lake.

## Endpoints

1. `GET /health`: Basic listener sanity.
2. `GET /ready`: Pings `localhost:5432`. Returns `503 Service Unavailable` if PostgreSQL connection is broken.
3. `GET /events`: Uses `SELECT raw FROM sae_telemetry ORDER BY timestamp DESC LIMIT 50`. Exposes full underlying stringified OCSF JSON seamlessly.
4. `GET /incidents`: Uses `json_build_object` inside PostgreSQL to query aggregated counts across active threat streams.
5. `GET /decisions`: Retrieves LangGraph's generated context, validation, risk score, and explicit orchestrator action.
6. `GET /investigations`: Aliased to decisions based on AI telemetry.
