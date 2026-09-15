curl -s -X POST http://localhost:5001/api/v1/hooks \
  -H "Authorization: Bearer 85cb3e87-7126-46b3-b639-2d412052c658" \
  -H "Content-Type: application/json" \
  -d '{"workflow_id": "c68ed44f-efb0-418b-9b24-108b35c189fd", "name": "sae"}'
