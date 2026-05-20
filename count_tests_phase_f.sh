#!/bin/bash
echo "=== Go Services Test Count ==="
for d in /home/kangjh3kang/Manpasik/backend/services/*/; do
    svc=$(basename "$d")
    cnt=$(grep -rc 'func Test' "$d" --include='*_test.go' 2>/dev/null | awk -F: '{s+=$2} END{print s+0}')
    echo "$cnt $svc"
done | sort -rn

echo ""
echo "=== Flutter Test Count ==="
total=0
for f in $(find /home/kangjh3kang/Manpasik/frontend/flutter-app/test -name '*_test.dart' 2>/dev/null); do
    cnt=$(grep -c 'test(\|testWidgets(' "$f" 2>/dev/null)
    total=$((total + cnt))
    rel=${f#/home/kangjh3kang/Manpasik/frontend/flutter-app/test/}
    echo "$cnt $rel"
done | sort -rn
echo ""
echo "Flutter TOTAL: $total"

echo ""
echo "=== Features without tests ==="
for d in /home/kangjh3kang/Manpasik/frontend/flutter-app/lib/features/*/; do
    feat=$(basename "$d")
    if [ ! -d "/home/kangjh3kang/Manpasik/frontend/flutter-app/test/features/$feat" ]; then
        echo "NO TEST: $feat"
    fi
done
