#!/bin/bash
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
  echo -n "Testing $svc... "
  cd "$BASE/$svc"
  OUTPUT=$(GOWORK=off go test ./... 2>&1)
  EXIT=$?
  if [ $EXIT -eq 0 ]; then
    # Count actual test runs
    TEST_COUNT=$(echo "$OUTPUT" | grep -c "^ok")
    echo "PASS ($TEST_COUNT packages)"
    PASS=$((PASS + 1))
  else
    echo "FAIL (exit=$EXIT)"
    echo "$OUTPUT" | grep -E "FAIL|Error|cannot" | head -5
    FAIL=$((FAIL + 1))
    FAIL_LIST="$FAIL_LIST $svc"
  fi
done

echo ""
echo "=== Test Results ==="
echo "PASS: $PASS / ${#SERVICES[@]}"
echo "FAIL: $FAIL"
if [ -n "$FAIL_LIST" ]; then
  echo "Failed:$FAIL_LIST"
fi
