#!/bin/bash
export PATH=/home/kangjh3kang/flutter/bin:/usr/local/go/bin:/usr/local/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik/frontend/flutter-app

flutter test 2>&1 | tee /tmp/flutter_test_out.txt
echo "---EXIT_CODE: $?---"
echo "---LAST_LINES---"
tail -5 /tmp/flutter_test_out.txt
echo "---GREP_RESULT---"
grep -oE '[0-9]+:[0-9]+ \+[0-9]+[^:]*' /tmp/flutter_test_out.txt | tail -1
