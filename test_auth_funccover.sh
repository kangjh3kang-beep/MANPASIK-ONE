#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik/backend/services/auth-service
GOWORK=off go test ./internal/service/ -coverprofile=cover.out -count=1 2>&1 | tail -3
GOWORK=off go tool cover -func=cover.out 2>&1 | grep -v "0.0%" | head -30
echo "---"
echo "=== 0% 함수 ==="
GOWORK=off go tool cover -func=cover.out 2>&1 | grep "0.0%"
rm -f cover.out
