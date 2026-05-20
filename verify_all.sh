#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik/backend

echo "=== Go 전체 빌드 ==="
FAIL=0
PASS=0
for svc in $(ls -d services/*/); do
  svcname=$(basename "$svc")
  if GOWORK=off go build "./$svc..." 2>/dev/null; then
    PASS=$((PASS+1))
  else
    echo "FAIL: $svcname"
    FAIL=$((FAIL+1))
  fi
done
echo "빌드: PASS=$PASS FAIL=$FAIL"
echo ""

echo "=== Go 전체 테스트 ==="
GOWORK=off go test ./... 2>&1 | tail -50
echo ""
echo "=== 테스트 카운트 ==="
GOWORK=off go test ./... 2>&1 | grep -E "^ok|FAIL" | wc -l
