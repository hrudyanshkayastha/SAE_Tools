#!/bin/bash
curl -s -X POST http://localhost:5001/api/v1/login -H 'Content-Type: application/json' -d '{"username":"admin@example.com","password":"Password123!"}' > response.json
API_KEY=$(jq -r '.active_org.users[] | select(.username=="admin@example.com") | .apikey' response.json)
if [ "$API_KEY" = "null" ] || [ -z "$API_KEY" ]; then
    API_KEY=$(jq -r '.active_org.users[0].apikey' response.json)
fi
echo "API_KEY=$API_KEY"
