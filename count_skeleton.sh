#!/bin/bash
for svc in ai-inference analytics emergency marketplace nlp iot-gateway vision assistant concept voice-profile data-platform; do
  f="/home/kangjh3kang/Manpasik/backend/services/${svc}-service/internal/service"
  if [ -d "$f" ]; then
    lines=$(wc -l $f/*.go 2>/dev/null | tail -1 | awk '{print $1}')
    files=$(ls $f/*.go 2>/dev/null | wc -l)
    echo "$svc: $lines LOC ($files files)"
  fi
done
