#!/bin/bash
export PATH=/home/kangjh3kang/flutter/bin:/usr/local/go/bin:/usr/local/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik/frontend/flutter-app

flutter test --reporter expanded 2>&1 > /tmp/flutter_expanded.txt
exit_code=$?
echo "EXIT_CODE=$exit_code"

pass_count=$(grep -cE '^\s*\+[0-9]+' /tmp/flutter_expanded.txt)
fail_count=$(grep -cE '^\s*-[0-9]+' /tmp/flutter_expanded.txt)
total_lines=$(wc -l < /tmp/flutter_expanded.txt)
echo "PASS=$pass_count FAIL=$fail_count TOTAL_LINES=$total_lines"

# Show last 20 lines
echo "---LAST20---"
tail -20 /tmp/flutter_expanded.txt
