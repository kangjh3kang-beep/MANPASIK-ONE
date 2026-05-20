#!/bin/bash
export PATH=/home/kangjh3kang/flutter/bin:/usr/local/go/bin:/usr/local/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik/frontend/flutter-app
flutter test 2>&1 | tee /tmp/ft_result.txt | tail -3
echo "---EXIT---"
echo "TOTAL_TESTS=$(grep -c '+[0-9]' /tmp/ft_result.txt)"
echo "FAILURES=$(grep -c 'FAILED\|Some tests failed' /tmp/ft_result.txt)"
