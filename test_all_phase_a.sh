#!/bin/bash
# Phase A: 13개 서비스 테스트 검증
export PATH=/usr/local/go/bin:/usr/bin:/bin
BASE="/home/kangjh3kang/Manpasik/backend/services"

PASS=0
FAIL=0
SKIP=0
FAIL_LIST=""

SERVICES=(
  analytics-service
  emergency-service
  audit-service
  digital-twin-service
  iot-gateway-service
  vision-service
  assistant-service
  marketplace-service
  concept-service
  cartridge-store-service
  nlp-service
  data-platform-service
  voice-profile-service
)

for svc in "${SERVICES[@]}"; do
  echo -n "Testing $svc... "
  cd "$BASE/$svc"
  OUTPUT=$(GOWORK=off go test ./... 2>&1)
  if echo "$OUTPUT" | grep -q "FAIL"; then
    echo "FAIL"
    echo "$OUTPUT" | tail -5
    FAIL=$((FAIL + 1))
    FAIL_LIST="$FAIL_LIST $svc"
  elif echo "$OUTPUT" | grep -q "no test files"; then
    echo "SKIP (no tests)"
    SKIP=$((SKIP + 1))
  else
    echo "PASS"
    PASS=$((PASS + 1))
  fi
done

echo ""
echo "=== Test Results ==="
echo "PASS: $PASS"
echo "SKIP: $SKIP"
echo "FAIL: $FAIL"
if [ -n "$FAIL_LIST" ]; then
  echo "Failed:$FAIL_LIST"
fi
