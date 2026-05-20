#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik
FAIL=0
PASS=0
for d in backend/services/*/; do
    svc=$(basename "$d")
    result=$(cd "$d" && GOWORK=off go test ./... 2>&1)
    if echo "$result" | grep -q "FAIL"; then
        echo "FAIL: $svc"
        echo "$result" | tail -3
        FAIL=$((FAIL+1))
    else
        PASS=$((PASS+1))
    fi
done
echo ""
echo "=== Go SUMMARY: $PASS PASS, $FAIL FAIL ==="
