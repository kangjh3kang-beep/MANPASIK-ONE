#!/bin/bash
# Phase E-1: 전체 Go 서비스 빌드 검증
export PATH=/usr/local/go/bin:/usr/bin:/bin
BASE=/home/kangjh3kang/Manpasik/backend

PASS=0
FAIL=0
SKIP=0
FAILED_SVCS=""

for svc_dir in "$BASE"/services/*/; do
  svc=$(basename "$svc_dir")
  if [ ! -d "$svc_dir/cmd" ]; then
    SKIP=$((SKIP+1))
    continue
  fi
  cd "$svc_dir"
  if GOWORK=off go build ./... 2>&1 | tail -3; then
    PASS=$((PASS+1))
    echo "  [PASS] $svc"
  else
    FAIL=$((FAIL+1))
    FAILED_SVCS="$FAILED_SVCS $svc"
    echo "  [FAIL] $svc"
  fi
done

echo ""
echo "=== BUILD RESULT ==="
echo "PASS: $PASS | FAIL: $FAIL | SKIP: $SKIP"
if [ -n "$FAILED_SVCS" ]; then
  echo "FAILED:$FAILED_SVCS"
else
  echo "ALL PASS"
fi
