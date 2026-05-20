#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik/backend/services
total_pass=0
total_fail=0
total_funcs=0
fail_details=""

for svc in */; do
  svc="${svc%/}"
  if [ -f "$svc/go.mod" ]; then
    # Count test functions
    test_count=0
    for tf in $(find "$svc" -name "*_test.go" 2>/dev/null); do
      c=$(grep -c "^func Test" "$tf" 2>/dev/null || echo 0)
      test_count=$((test_count + c))
    done
    total_funcs=$((total_funcs + test_count))

    # Run tests
    output=$(cd "$svc" && GOWORK=off go test ./... 2>&1)
    f=$(echo "$output" | grep -c "^FAIL")
    if [ "$f" -gt 0 ]; then
      total_fail=$((total_fail + 1))
      echo "FAIL $svc (tests:$test_count)"
      fail_details="$fail_details\n--- $svc ---\n$(echo "$output" | grep -E '(FAIL|---FAIL|panic)' | head -5)"
    else
      total_pass=$((total_pass + 1))
      echo "PASS $svc (tests:$test_count)"
    fi
  fi
done

echo ""
echo "=== TEST RESULT: $total_pass services pass, $total_fail services fail ==="
echo "=== TOTAL TEST FUNCTIONS: $total_funcs ==="
if [ -n "$fail_details" ]; then
  echo -e "\n=== FAILURE DETAILS ===$fail_details"
fi
