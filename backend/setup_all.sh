#!/bin/bash
TOKEN=$(curl -s -X POST http://localhost:5001/api/v1/login -H 'Content-Type: application/json' -d '{"username":"admin@example.com","password":"Password123!"}' | jq -r '.cookies[] | select(.key=="session_token") | .value')
WF_ID=$(curl -s -X POST http://localhost:5001/api/v1/workflows -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' -d '{"name":"Runtime_Test","description":"test"}' | jq -r '.id')
if [ "$WF_ID" = "null" ]; then
    WF_ID=$(curl -s -X GET http://localhost:5001/api/v1/workflows -H "Authorization: Bearer $TOKEN" | jq -r '.[0].id')
fi
echo "TOKEN=$TOKEN"
echo "WF_ID=$WF_ID"
