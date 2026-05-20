#!/bin/bash
export PATH=/home/kangjh3kang/flutter/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik/frontend/flutter-app
flutter test > /tmp/ft_result.txt 2>&1
echo "EXIT_CODE=$?"
tail -3 /tmp/ft_result.txt
