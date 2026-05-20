#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin

echo "=== admin-service 테스트 ==="
cd /home/kangjh3kang/Manpasik/backend/services/admin-service
GOWORK=off go test ./... -v -count=1 2>&1 | grep -E "^(--- |ok |FAIL)" | tail -40
echo ""
echo "=== assistant-service 테스트 ==="
cd /home/kangjh3kang/Manpasik/backend/services/assistant-service
GOWORK=off go test ./... -v -count=1 2>&1 | grep -E "^(--- |ok |FAIL)" | tail -40
echo ""
echo "=== coaching-service 테스트 ==="
cd /home/kangjh3kang/Manpasik/backend/services/coaching-service
GOWORK=off go test ./... -v -count=1 2>&1 | grep -E "^(--- |ok |FAIL)" | tail -40
