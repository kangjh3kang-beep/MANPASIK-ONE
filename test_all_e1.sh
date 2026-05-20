#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik/backend

echo "=== Building ALL Go services ==="
FAIL_BUILD=0
for dir in services/*/; do
    svc=$(basename "$dir")
    if GOWORK=off go build "./$dir..." 2>/dev/null; then
        echo "  BUILD OK: $svc"
    else
        echo "  BUILD FAIL: $svc"
        FAIL_BUILD=$((FAIL_BUILD + 1))
    fi
done

echo ""
echo "=== Testing ALL Go services ==="
FAIL_TEST=0
TOTAL_PKG=0
for dir in services/*/; do
    svc=$(basename "$dir")
    output=$(GOWORK=off go test "./$dir..." -count=1 2>&1)
    rc=$?
    pkg_count=$(echo "$output" | grep -c "ok\|FAIL\|no test files")
    TOTAL_PKG=$((TOTAL_PKG + pkg_count))
    if [ $rc -eq 0 ]; then
        echo "  TEST OK: $svc ($pkg_count pkgs)"
    else
        echo "  TEST FAIL: $svc"
        echo "$output" | grep "FAIL"
        FAIL_TEST=$((FAIL_TEST + 1))
    fi
done

# shared packages
for dir in shared/*/; do
    pkg=$(basename "$dir")
    output=$(GOWORK=off go test "./$dir..." -count=1 2>&1)
    rc=$?
    if [ $rc -eq 0 ]; then
        echo "  TEST OK: shared/$pkg"
    else
        echo "  TEST FAIL: shared/$pkg"
        FAIL_TEST=$((FAIL_TEST + 1))
    fi
done

echo ""
echo "=== SUMMARY ==="
echo "Build failures: $FAIL_BUILD"
echo "Test failures: $FAIL_TEST"
echo "Total packages: $TOTAL_PKG"
