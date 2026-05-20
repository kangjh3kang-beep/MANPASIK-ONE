#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik/backend

echo "=== E-4 전체 백엔드 빌드 ==="
FAIL=0
PASS=0
for dir in services/*/; do
  svc=$(basename $dir)
  if GOWORK=off go build ./$dir... 2>&1 | grep -q "error"; then
    echo "  FAIL $svc"
    FAIL=$((FAIL+1))
  else
    PASS=$((PASS+1))
  fi
done
echo "빌드 PASS: $PASS, FAIL: $FAIL"

echo ""
echo "=== E-4 전체 서비스 테스트 ==="
GOWORK=off go test -count=1 -timeout 180s $(GOWORK=off go list ./services/*/internal/service/... 2>/dev/null) 2>&1 | tail -80
