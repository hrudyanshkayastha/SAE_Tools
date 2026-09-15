sleep 5
curl -s -X POST http://localhost:5001/api/v1/register -H 'Content-Type: application/json' -d '{"username":"admin@example.com","password":"Password123!"}' > /dev/null
TOKEN=\
echo "TOKEN=\"
WF_ID=\
echo "WF_ID=\"
