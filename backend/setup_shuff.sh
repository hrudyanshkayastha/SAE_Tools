#!/bin/bash
curl -s -X POST http://localhost:5001/api/v1/register -H 'Content-Type: application/json' -d '{"username":"admin@example.com","password":"Password123!"}' > /dev/null
TOKEN=$(curl -s -X POST http://localhost:5001/api/v1/login -H 'Content-Type: application/json' -d '{"username":"admin@example.com","password":"Password123!"}' | grep -oP '"success":"\K[^"]+')
echo "TOKEN=$TOKEN"
WF_ID=$(curl -s -X POST http://localhost:5001/api/v1/workflows -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' -d '{"name":"Runtime_Test","description":"test"}' | grep -oP '"id":"\K[^"]+')
echo "WF_ID=$WF_ID"
