#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik/backend

echo "=== Phase E-2 전체 검증 ==="

echo "--- 35/35 빌드 ---"
FAIL=0
for svc in $(ls -d services/*/cmd 2>/dev/null | sed 's|services/||;s|/cmd||'); do
  if GOWORK=off go build -o /dev/null "./services/$svc/cmd" 2>&1; then
    :
  else
    echo "  FAIL: $svc"
    FAIL=$((FAIL+1))
  fi
done
echo "빌드 실패: $FAIL"

echo ""
echo "--- 전체 서비스 테스트 ---"
GOWORK=off go test -count=1 -timeout 180s $(GOWORK=off go list ./services/.../service/... ./services/.../webrtc/... 2>/dev/null) 2>&1 | grep -E '^(ok|FAIL|---)'
echo ""
echo "--- ai-inference 테스트 수 ---"
GOWORK=off go test -count=1 -v ./services/ai-inference-service/internal/service/... 2>&1 | grep -c '--- PASS'
echo ""
echo "--- coaching 테스트 수 ---"
GOWORK=off go test -count=1 -v ./services/coaching-service/internal/service/... 2>&1 | grep -c '--- PASS'
