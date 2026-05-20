#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik/backend/services
pass=0
fail=0
fail_list=""
for svc in */; do
  svc="${svc%/}"
  if [ -f "$svc/go.mod" ]; then
    if (cd "$svc" && GOWORK=off go build ./... 2>/dev/null); then
      pass=$((pass+1))
      echo "OK $svc"
    else
      fail=$((fail+1))
      fail_list="$fail_list $svc"
      echo "FAIL $svc"
    fi
  fi
done
echo ""
echo "=== BUILD RESULT: $pass pass, $fail fail ==="
if [ -n "$fail_list" ]; then
  echo "FAILED:$fail_list"
fi
