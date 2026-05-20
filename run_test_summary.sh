#!/bin/bash
export PATH=/home/kangjh3kang/flutter/bin:/usr/local/go/bin:/usr/local/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik/frontend/flutter-app
flutter test > /tmp/ft_out.txt 2>&1
EXITCODE=$?
tail -3 /tmp/ft_out.txt > /tmp/ft_summary.txt
echo "EXIT_CODE=$EXITCODE" >> /tmp/ft_summary.txt
