#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik/backend

echo "=== Go 35/35 서비스 빌드 ==="
FAIL=0
for svc in $(ls -d services/*/cmd 2>/dev/null | sed 's|services/||;s|/cmd||'); do
  if GOWORK=off go build -o /dev/null "./services/$svc/cmd" 2>&1; then
    echo "  OK: $svc"
  else
    echo "  FAIL: $svc"
    FAIL=$((FAIL+1))
  fi
done

echo ""
echo "=== Go 테스트 (서비스 + shared) ==="
GOWORK=off go test ./services/*/internal/service/... ./services/*/internal/webrtc/... ./shared/... -count=1 -timeout 120s 2>&1 | tail -80

echo ""
echo "빌드 실패: $FAIL"
