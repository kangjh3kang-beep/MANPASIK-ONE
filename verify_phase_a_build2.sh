#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin
BASE=/home/kangjh3kang/Manpasik/backend/services

PASS=0
FAIL=0
ERRORS=""

for svc_dir in $BASE/*/; do
  svc=$(basename "$svc_dir")
  if [ ! -f "$svc_dir/cmd/main.go" ]; then
    continue
  fi
  result=$(cd "$svc_dir" && GOWORK=off go build ./... 2>&1)
  if [ $? -eq 0 ]; then
    PASS=$((PASS+1))
    echo "OK: $svc"
  else
    FAIL=$((FAIL+1))
    ERRORS="$ERRORS\n--- $svc ---\n$result"
    echo "FAIL: $svc"
  fi
done

echo ""
echo "=== BUILD RESULTS ==="
echo "PASS: $PASS"
echo "FAIL: $FAIL"
if [ $FAIL -gt 0 ]; then
  echo ""
  echo "=== ERRORS ==="
  echo -e "$ERRORS"
fi
