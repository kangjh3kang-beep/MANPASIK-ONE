#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik/backend

echo "=== Phase C-1 전체 검증 ==="
echo ""

# 1. 전체 빌드
echo "--- [1/3] 전체 빌드 검증 ---"
FAIL_BUILD=0
for svc_dir in services/*/; do
    svc_name=$(basename "$svc_dir")
    if GOWORK=off go build "./$svc_dir..." 2>/dev/null; then
        echo "  OK  $svc_name"
    else
        echo "  FAIL $svc_name"
        FAIL_BUILD=$((FAIL_BUILD + 1))
    fi
done

echo ""
echo "--- [2/3] 전체 테스트 검증 ---"
FAIL_TEST=0
PASS_TEST=0
for svc_dir in services/*/; do
    svc_name=$(basename "$svc_dir")
    result=$(GOWORK=off go test "./$svc_dir..." -count=1 2>&1)
    if echo "$result" | grep -q "FAIL"; then
        echo "  FAIL $svc_name"
        echo "$result" | grep "FAIL" | head -3
        FAIL_TEST=$((FAIL_TEST + 1))
    else
        PASS_TEST=$((PASS_TEST + 1))
        echo "  OK  $svc_name"
    fi
done

echo ""
echo "--- [3/3] shared 패키지 ---"
shared_result=$(GOWORK=off go test ./shared/... -count=1 2>&1)
if echo "$shared_result" | grep -q "FAIL"; then
    echo "  FAIL shared"
    FAIL_TEST=$((FAIL_TEST + 1))
else
    echo "  OK  shared"
    PASS_TEST=$((PASS_TEST + 1))
fi

echo ""
echo "=== 결과 ==="
echo "빌드 실패: $FAIL_BUILD"
echo "테스트 성공: $PASS_TEST"
echo "테스트 실패: $FAIL_TEST"
if [ $FAIL_BUILD -eq 0 ] && [ $FAIL_TEST -eq 0 ]; then
    echo "ALL PASS ✓"
else
    echo "FAIL ✗"
fi
