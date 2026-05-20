#!/bin/bash
# Phase A: 13개 서비스 빌드 검증
export PATH=/usr/local/go/bin:/usr/bin:/bin
BASE="/home/kangjh3kang/Manpasik/backend/services"

PASS=0
FAIL=0
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
  echo -n "Building $svc... "
  cd "$BASE/$svc"
  if GOWORK=off go build ./... 2>&1; then
    echo "PASS"
    PASS=$((PASS + 1))
  else
    echo "FAIL"
    FAIL=$((FAIL + 1))
    FAIL_LIST="$FAIL_LIST $svc"
  fi
done

echo ""
echo "=== Build Results ==="
echo "PASS: $PASS / ${#SERVICES[@]}"
echo "FAIL: $FAIL"
if [ -n "$FAIL_LIST" ]; then
  echo "Failed:$FAIL_LIST"
fi
