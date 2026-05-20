#!/bin/bash
set -e
export PATH=/usr/local/go/bin:/usr/bin:/bin

BASE="/home/kangjh3kang/Manpasik/backend/services"
PASS=0
FAIL=0
FAIL_LIST=""

SERVICES=(
  auth-service user-service subscription-service shop-service device-service
  measurement-service cartridge-service calibration-service coaching-service
  health-record-service notification-service family-service payment-service
  reservation-service prescription-service community-service video-service
  telemedicine-service translation-service admin-service
  ai-inference-service analytics-service emergency-service marketplace-service
  nlp-service iot-gateway-service vision-service assistant-service
  concept-service voice-profile-service data-platform-service
  cartridge-store-service gateway
)

echo "=== Phase B3 빌드 검증 ==="
for svc in "${SERVICES[@]}"; do
  DIR="$BASE/$svc"
  if [ ! -d "$DIR" ]; then continue; fi
  cd "$DIR"
  if GOWORK=off go build ./... 2>/dev/null; then
    PASS=$((PASS+1))
    echo "PASS $svc"
  else
    FAIL=$((FAIL+1))
    FAIL_LIST="$FAIL_LIST $svc"
    echo "FAIL $svc"
  fi
done
echo "=== 빌드: $PASS PASS / $FAIL FAIL ==="
if [ -n "$FAIL_LIST" ]; then echo "실패: $FAIL_LIST"; fi

echo ""
echo "=== B3 대상 서비스 테스트 ==="
B3_SERVICES=(vision-service assistant-service ai-inference-service)
T_PASS=0
T_FAIL=0
for svc in "${B3_SERVICES[@]}"; do
  DIR="$BASE/$svc"
  cd "$DIR"
  if GOWORK=off go test ./... -count=1 2>&1 | grep -q "^ok"; then
    count=$(GOWORK=off go test ./... -v -count=1 2>&1 | grep -c "^--- PASS" || true)
    T_PASS=$((T_PASS+1))
    echo "PASS $svc ($count 테스트)"
  else
    T_FAIL=$((T_FAIL+1))
    echo "FAIL $svc"
    GOWORK=off go test ./... -v -count=1 2>&1 | tail -10
  fi
done
echo "=== B3 테스트: $T_PASS PASS / $T_FAIL FAIL ==="
