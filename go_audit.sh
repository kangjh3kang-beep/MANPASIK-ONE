#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik/backend

echo "=== Go 서비스별 빌드 검증 ==="
build_pass=0
build_fail=0
build_fail_list=""

for svc in services/*/; do
  svc_name=$(basename "$svc")
  if [ -d "$svc/cmd" ]; then
    if GOWORK=off go build "./$svc/..." 2>/dev/null; then
      build_pass=$((build_pass+1))
    else
      build_fail=$((build_fail+1))
      build_fail_list="$build_fail_list $svc_name"
      echo "BUILD FAIL: $svc_name"
    fi
  fi
done
echo "BUILD: $build_pass pass, $build_fail fail"
if [ -n "$build_fail_list" ]; then
  echo "BUILD FAILURES:$build_fail_list"
fi

echo ""
echo "=== Go 서비스별 테스트 ==="
test_pass=0
test_fail=0
test_fail_list=""
total_funcs=0

for svc in services/*/; do
  svc_name=$(basename "$svc")
  # Count test functions
  tc=0
  while IFS= read -r tf; do
    c=$(grep -c '^func Test' "$tf" 2>/dev/null || echo 0)
    tc=$((tc + c))
  done < <(find "$svc" -name '*_test.go' 2>/dev/null)
  total_funcs=$((total_funcs + tc))

  if [ "$tc" -gt 0 ]; then
    output=$(GOWORK=off go test "./$svc/..." 2>&1)
    if echo "$output" | grep -q '^FAIL'; then
      test_fail=$((test_fail+1))
      test_fail_list="$test_fail_list $svc_name"
      echo "TEST FAIL: $svc_name ($tc funcs)"
      echo "$output" | grep -E '(FAIL|---\s*FAIL)' | head -3
    else
      test_pass=$((test_pass+1))
      echo "TEST PASS: $svc_name ($tc funcs)"
    fi
  else
    echo "NO TESTS: $svc_name"
  fi
done

# shared tests
shared_tc=0
while IFS= read -r tf; do
  c=$(grep -c '^func Test' "$tf" 2>/dev/null || echo 0)
  shared_tc=$((shared_tc + c))
done < <(find shared -name '*_test.go' 2>/dev/null)
total_funcs=$((total_funcs + shared_tc))

if [ "$shared_tc" -gt 0 ]; then
  output=$(GOWORK=off go test ./shared/... 2>&1)
  if echo "$output" | grep -q '^FAIL'; then
    test_fail=$((test_fail+1))
    echo "TEST FAIL: shared ($shared_tc funcs)"
  else
    test_pass=$((test_pass+1))
    echo "TEST PASS: shared ($shared_tc funcs)"
  fi
fi

echo ""
echo "=== SUMMARY ==="
echo "BUILD: $build_pass pass, $build_fail fail"
echo "TEST: $test_pass pass, $test_fail fail"
echo "TOTAL TEST FUNCTIONS: $total_funcs"
