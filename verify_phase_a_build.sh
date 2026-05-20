#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik/backend/services

PASS=0
FAIL=0
ERRORS=""

for svc in */; do
  svc=${svc%/}
  if [ ! -f "$svc/cmd/main.go" ]; then
    continue
  fi
  cd "/home/kangjh3kang/Manpasik/backend/services/$svc"
  result=$(GOWORK=off go build ./... 2>&1)
  if [ $? -eq 0 ]; then
    PASS=$((PASS+1))
  else
    FAIL=$((FAIL+1))
    ERRORS="$ERRORS\n--- $svc ---\n$result"
  fi
done

echo "=== BUILD RESULTS ==="
echo "PASS: $PASS"
echo "FAIL: $FAIL"
if [ $FAIL -gt 0 ]; then
  echo ""
  echo "=== ERRORS ==="
  echo -e "$ERRORS"
fi
