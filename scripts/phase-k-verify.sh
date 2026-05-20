#!/bin/bash
# Phase K 8개 서비스 일괄 검증
export PATH=/usr/local/go/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik

SERVICES=(nlp-service vision-service marketplace-service iot-gateway-service ai-inference-service assistant-service concept-service shop-service)

PASS=0
FAIL=0
RESULTS=()

for svc in "${SERVICES[@]}"; do
  pushd /home/kangjh3kang/Manpasik/backend/services/$svc > /dev/null
  TEST_OUT=$(GOWORK=off go test ./... 2>&1)
  EXIT=$?
  if [ $EXIT -eq 0 ]; then
    echo "[PASS] $svc"
    PASS=$((PASS+1))
  else
    echo "[FAIL] $svc"
    echo "$TEST_OUT" | tail -10
    FAIL=$((FAIL+1))
  fi
  popd > /dev/null
done

echo ""
echo "=== Phase K 종합 결과 ==="
echo "PASS: $PASS / 8"
echo "FAIL: $FAIL / 8"
