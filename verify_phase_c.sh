#!/bin/bash
# Phase C 전체 빌드+테스트 검증 스크립트
export PATH=/usr/local/go/bin:/usr/bin:/bin

BASE=/home/kangjh3kang/Manpasik/backend
cd "$BASE"

TOTAL=0
BUILD_PASS=0
BUILD_FAIL=0
TEST_PASS=0
TEST_FAIL=0
FAIL_LIST=""

for svc_dir in services/*/; do
  svc=$(basename "$svc_dir")
  # cmd/main.go 또는 internal/ 가 있는 서비스만 대상
  if [ ! -d "$svc_dir/cmd" ] && [ ! -d "$svc_dir/internal" ]; then
    continue
  fi
  TOTAL=$((TOTAL + 1))

  # Build
  if GOWORK=off go build "./$svc_dir..." 2>/dev/null; then
    BUILD_PASS=$((BUILD_PASS + 1))
  else
    BUILD_FAIL=$((BUILD_FAIL + 1))
    FAIL_LIST="$FAIL_LIST BUILD:$svc"
    echo "FAIL BUILD: $svc"
    continue
  fi

  # Test
  TEST_OUT=$(GOWORK=off go test "./$svc_dir..." 2>&1)
  if [ $? -eq 0 ]; then
    TEST_PASS=$((TEST_PASS + 1))
  else
    TEST_FAIL=$((TEST_FAIL + 1))
    FAIL_LIST="$FAIL_LIST TEST:$svc"
    echo "FAIL TEST: $svc"
    echo "$TEST_OUT" | tail -5
  fi
done

echo ""
echo "========================================="
echo "Phase C 검증 결과"
echo "========================================="
echo "총 서비스: $TOTAL"
echo "빌드 PASS: $BUILD_PASS / $TOTAL"
echo "빌드 FAIL: $BUILD_FAIL"
echo "테스트 PASS: $TEST_PASS / $TOTAL"
echo "테스트 FAIL: $TEST_FAIL"
if [ -n "$FAIL_LIST" ]; then
  echo "실패 목록: $FAIL_LIST"
fi
echo "========================================="
