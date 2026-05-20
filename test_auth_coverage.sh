#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik/backend/services/auth-service
echo "=== auth-service 테스트 ==="
GOWORK=off go test ./... -v -count=1 2>&1 | grep -E "^(--- |ok |FAIL|=== RUN)" | tail -50
echo ""
echo "=== auth-service 커버리지 ==="
GOWORK=off go test ./... -cover -count=1 2>&1 | grep -E "(coverage|FAIL)"
