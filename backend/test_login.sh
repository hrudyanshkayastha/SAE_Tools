#!/bin/bash
curl -s -X POST http://localhost:5001/api/v1/login -H 'Content-Type: application/json' -d '{"username":"admin@example.com","password":"Password123!"}'
