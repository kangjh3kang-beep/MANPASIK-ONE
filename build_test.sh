#!/bin/bash
cd /home/kangjh3kang/Manpasik/backend

echo "=== Building subscription-service ==="
go build ./services/subscription-service/cmd/ 2>&1
echo "Build exit: $?"

echo "=== Running subscription-service tests ==="
go test -v -count=1 ./services/subscription-service/internal/service/ 2>&1
echo "Test exit: $?"

echo "=== Building all services ==="
go build ./... 2>&1
echo "All build exit: $?"

echo "=== DONE ==="
