#!/bin/bash
cd /mnt/e/'New folder'/SAE_Tools/SAE/backend
export SAE_SHUFFLE_WEBHOOK="http://localhost:5001/api/v1/workflows/4f455a65-2ac8-40fc-8b8b-0efa5fc6c4ce/execute"
export SAE_SHUFFLE_AUTH_TOKEN="5abcac45-85f4-4c95-ad1a-27e3b35ea13e"
export REDIS_ADDR="localhost:6379"
export POSTGRES_DSN="postgres://postgres:postgres@localhost:5432/sae?sslmode=disable"
export OLLAMA_URL="http://localhost:11434"
go build -o sae_daemon main.go
./sae_daemon > daemon.log 2>&1 &
echo "$!" > daemon.pid
