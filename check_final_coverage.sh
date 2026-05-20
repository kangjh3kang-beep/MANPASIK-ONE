#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin
BASE="/home/kangjh3kang/Manpasik/backend/services"

PASS_COUNT=0
FAIL_COUNT=0

for svc_dir in "$BASE"/*/; do
    svc_name=$(basename "$svc_dir")
    cd "/home/kangjh3kang/Manpasik/backend" || continue

    # 서비스 테스트 패키지 경로
    test_files=$(find "$svc_dir" -name "*_test.go" 2>/dev/null | head -1)

    if [ -z "$test_files" ]; then
        result=$(GOWORK=off go build "./services/$svc_name/..." 2>&1)
        if [ $? -eq 0 ]; then
            PASS_COUNT=$((PASS_COUNT + 1))
            echo "BUILD $svc_name [no tests]"
        else
            FAIL_COUNT=$((FAIL_COUNT + 1))
            echo "FAIL  $svc_name [build]"
        fi
    else
        output=$(GOWORK=off go test "./services/$svc_name/..." -cover -count=1 2>&1)
        status=$?
        if [ $status -eq 0 ]; then
            PASS_COUNT=$((PASS_COUNT + 1))
            svc_cov=$(echo "$output" | grep "internal/service" | grep -oP '[\d.]+%' | head -1)
            if [ -z "$svc_cov" ]; then
                svc_cov="N/A"
            fi
            echo "PASS  $svc_name  service=$svc_cov"
        else
            FAIL_COUNT=$((FAIL_COUNT + 1))
            echo "FAIL  $svc_name"
            echo "$output" | grep "FAIL" | head -3
        fi
    fi
done

echo ""
echo "=== 총: $((PASS_COUNT + FAIL_COUNT)), PASS: $PASS_COUNT, FAIL: $FAIL_COUNT ==="
