#!/bin/bash
curl -s -X POST http://localhost:5001/api/v1/workflows -b 'session_token=0afbf9a2-2596-4cdd-9c5c-80d2cbf4757b' -H 'Content-Type: application/json' -d '{"name":"Runtime_Test","description":"test"}' | jq -r '.id'
