#!/bin/bash
for svc in ai-inference analytics emergency marketplace nlp iot-gateway vision assistant concept voice-profile data-platform; do
  f="/home/kangjh3kang/Manpasik/backend/services/${svc}-service/internal/service"
  test_count=$(grep -c "func Test" $f/*_test.go 2>/dev/null || echo 0)
  echo "$svc: $test_count tests"
done
