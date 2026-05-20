#!/bin/bash
export PATH=/home/kangjh3kang/flutter/bin:/usr/local/go/bin:/usr/local/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik/frontend/flutter-app

result=$(flutter test 2>&1)
last_line=$(echo "$result" | grep -E '^\d+:\d+ \+' | tail -1)
echo "FLUTTER_TEST_RESULT: $last_line"

fail_count=$(echo "$result" | grep -c "FAILED\|Some tests failed")
echo "FAIL_COUNT: $fail_count"
