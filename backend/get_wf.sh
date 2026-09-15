#!/bin/bash
TOKEN="0afbf9a2-2596-4cdd-9c5c-80d2cbf4757b"
WF_ID=$(curl -s -X POST http://localhost:5001/api/v1/workflows -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' -d '{"name":"Runtime_Test","description":"test"}' | jq -r '.id')
if [ "$WF_ID" = "null" ]; then
    WF_ID=$(curl -s -X GET http://localhost:5001/api/v1/workflows -H "Authorization: Bearer $TOKEN" | jq -r '.[0].id')
fi
echo "WF_ID=$WF_ID"
