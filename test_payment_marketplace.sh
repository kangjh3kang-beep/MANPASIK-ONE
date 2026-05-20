#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin
echo "=== payment-service 테스트 ==="
cd /home/kangjh3kang/Manpasik/backend/services/payment-service
GOWORK=off go test ./... -v -count=1 2>&1 | grep -E "^(--- |ok |FAIL|=== RUN)" | tail -30
echo ""
echo "=== marketplace-service 테스트 ==="
cd /home/kangjh3kang/Manpasik/backend/services/marketplace-service
GOWORK=off go test ./... -v -count=1 2>&1 | grep -E "^(--- |ok |FAIL|=== RUN)" | tail -40
