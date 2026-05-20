#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin
BASE=/home/kangjh3kang/Manpasik/backend/services

PASS=0
FAIL=0
SKIP=0
ERRORS=""

for svc_dir in $BASE/*/; do
  svc=$(basename "$svc_dir")
  test_count=$(find "$svc_dir" -name '*_test.go' 2>/dev/null | wc -l)
  if [ "$test_count" -eq 0 ]; then
    SKIP=$((SKIP+1))
    continue
  fi
  result=$(cd "$svc_dir" && GOWORK=off go test ./... 2>&1)
  if [ $? -eq 0 ]; then
    PASS=$((PASS+1))
    echo "OK: $svc ($test_count test files)"
  else
    FAIL=$((FAIL+1))
    ERRORS="$ERRORS\n--- $svc ---\n$result"
    echo "FAIL: $svc"
  fi
done

echo ""
echo "=== TEST RESULTS ==="
echo "PASS: $PASS | FAIL: $FAIL | SKIP(no tests): $SKIP"
if [ $FAIL -gt 0 ]; then
  echo ""
  echo "=== ERRORS ==="
  echo -e "$ERRORS"
fi
