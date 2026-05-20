#!/bin/bash
# 모든 백엔드 서비스 빌드+테스트 검증
export PATH=/usr/local/go/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik/backend/services

PASS=0
FAIL=0
FAILED=()

for svc_dir in */; do
  svc=${svc_dir%/}
  pushd /home/kangjh3kang/Manpasik/backend/services/$svc > /dev/null
  if GOWORK=off go test ./... > /tmp/all-test-$svc.log 2>&1; then
    PASS=$((PASS+1))
  else
    FAIL=$((FAIL+1))
    FAILED+=("$svc")
  fi
  popd > /dev/null
done

echo "=== 전체 서비스 검증 결과 ==="
echo "PASS: $PASS"
echo "FAIL: $FAIL"
if [ ${#FAILED[@]} -gt 0 ]; then
  echo "실패 서비스:"
  for s in "${FAILED[@]}"; do
    echo "  - $s"
    tail -3 /tmp/all-test-$s.log | sed 's/^/    /'
  done
fi
