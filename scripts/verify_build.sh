#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin
BASE="/home/kangjh3kang/Manpasik/backend/services"
PASS=0
FAIL=0
for svc_dir in "$BASE"/*/; do
  svc=$(basename "$svc_dir")
  if [ -f "$svc_dir/go.mod" ]; then
    if (cd "$svc_dir" && GOWORK=off go build ./... 2>/dev/null); then
      echo "BUILD OK: $svc"
      PASS=$((PASS + 1))
    else
      echo "BUILD FAIL: $svc"
      FAIL=$((FAIL + 1))
    fi
  fi
done
echo ""
echo "BUILD TOTAL: $PASS pass, $FAIL fail"
