#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin
BASE="/home/kangjh3kang/Manpasik/backend/services"

PASS=0
FAIL=0
FAIL_LIST=""
TOTAL=0

for dir in "$BASE"/*/; do
  svc=$(basename "$dir")
  if [ ! -f "$dir/go.mod" ] && [ ! -d "$dir/cmd" ]; then
    continue
  fi
  TOTAL=$((TOTAL + 1))
  echo -n "Building $svc... "
  cd "$dir"
  if GOWORK=off go build ./... 2>/dev/null; then
    echo "PASS"
    PASS=$((PASS + 1))
  else
    echo "FAIL"
    FAIL=$((FAIL + 1))
    FAIL_LIST="$FAIL_LIST $svc"
  fi
done

echo ""
echo "=== Full Build Results ==="
echo "PASS: $PASS / $TOTAL"
echo "FAIL: $FAIL"
if [ -n "$FAIL_LIST" ]; then
  echo "Failed:$FAIL_LIST"
fi
