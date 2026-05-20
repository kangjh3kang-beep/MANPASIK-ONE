#!/bin/bash
# Phase E-2: 전체 Go 서비스 테스트 검증
export PATH=/usr/local/go/bin:/usr/bin:/bin
BASE=/home/kangjh3kang/Manpasik/backend

PASS=0
FAIL=0
SKIP=0
TOTAL_TESTS=0
FAILED_SVCS=""

for svc_dir in "$BASE"/services/*/; do
  svc=$(basename "$svc_dir")

  # 테스트 파일이 있는지 확인
  test_count=$(find "$svc_dir" -name "*_test.go" 2>/dev/null | wc -l)
  if [ "$test_count" -eq 0 ]; then
    SKIP=$((SKIP+1))
    continue
  fi

  cd "$svc_dir"
  output=$(GOWORK=off go test ./... 2>&1)
  exit_code=$?

  # 테스트 수 카운트
  tests_run=$(echo "$output" | grep -c "^ok")
  TOTAL_TESTS=$((TOTAL_TESTS+tests_run))

  if [ $exit_code -eq 0 ]; then
    PASS=$((PASS+1))
    echo "  [PASS] $svc ($tests_run packages)"
  else
    FAIL=$((FAIL+1))
    FAILED_SVCS="$FAILED_SVCS $svc"
    echo "  [FAIL] $svc"
    echo "$output" | grep -E "FAIL|Error" | head -3
  fi
done

echo ""
echo "=== TEST RESULT ==="
echo "Services: PASS=$PASS FAIL=$FAIL SKIP=$SKIP"
echo "Total test packages: $TOTAL_TESTS"
if [ -n "$FAILED_SVCS" ]; then
  echo "FAILED:$FAILED_SVCS"
else
  echo "ALL PASS"
fi
