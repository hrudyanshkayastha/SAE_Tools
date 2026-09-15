curl -s -X GET http://localhost:5001/api/v1/workflows/c68ed44f-efb0-418b-9b24-108b35c189fd \
  -H "Authorization: Bearer 85cb3e87-7126-46b3-b639-2d412052c658" > wf.json
jq '.triggers = [{"app_name":"Builtin","app_version":"1.0.0","id_":"webhook_b7da255a-43d9-43c0-beaa-ab9125da02c7","name":"Webhook","trigger_type":"WEBHOOK","status":"active"}]' wf.json > wf_updated.json
curl -s -X PUT http://localhost:5001/api/v1/workflows/c68ed44f-efb0-418b-9b24-108b35c189fd \
  -H "Authorization: Bearer 85cb3e87-7126-46b3-b639-2d412052c658" \
  -H "Content-Type: application/json" -d @wf_updated.json
