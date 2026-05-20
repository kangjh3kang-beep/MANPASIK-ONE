#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik
PASS=0
FAIL=0
for svc in backend/services/*/; do
  svc_name=$(basename "$svc")
  cd /home/kangjh3kang/Manpasik/"$svc"
  result=$(GOWORK=off go test ./... 2>&1)
  if echo "$result" | grep -q "FAIL"; then
    FAIL=$((FAIL+1))
    echo "FAIL: $svc_name"
    echo "$result" | grep -E "FAIL|Error"
  else
    PASS=$((PASS+1))
  fi
done
echo "=== Go Test Summary ==="
echo "PASS: $PASS services"
echo "FAIL: $FAIL services"
