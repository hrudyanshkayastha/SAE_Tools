#!/bin/bash
curl -s -X GET http://localhost:5001/api/v1/users -b 'session_token=0afbf9a2-2596-4cdd-9c5c-80d2cbf4757b' | jq '.'
