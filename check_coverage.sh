#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin
BASE="/home/kangjh3kang/Manpasik/backend/services"

echo "=== Go 서비스 테스트 커버리지 현황 ==="
echo ""
printf "%-30s %8s %8s %s\n" "서비스" "테스트수" "커버리지" "상태"
echo "---------------------------------------------------------------"

TOTAL_TESTS=0
TOTAL_SVCS=0
PASS_SVCS=0

for dir in "$BASE"/*/; do
  svc=$(basename "$dir")
  if [ ! -d "$dir" ]; then continue; fi
  cd "$dir"

  # 테스트 수 카운트
  test_count=$(GOWORK=off go test ./... -v -count=1 2>&1 | grep -c "^--- PASS" || true)

  # 커버리지 측정
  cov_output=$(GOWORK=off go test ./... -cover -count=1 2>&1)

  # 서비스 로직 패키지의 커버리지만 추출
  svc_cov=$(echo "$cov_output" | grep "internal/service" | grep -oP 'coverage: \K[0-9.]+' || echo "0")
  if [ -z "$svc_cov" ]; then svc_cov="0"; fi

  status="OK"
  if echo "$cov_output" | grep -q "FAIL"; then
    status="FAIL"
  else
    PASS_SVCS=$((PASS_SVCS+1))
  fi

  TOTAL_TESTS=$((TOTAL_TESTS + test_count))
  TOTAL_SVCS=$((TOTAL_SVCS+1))

  printf "%-30s %8d %7s%% %s\n" "$svc" "$test_count" "$svc_cov" "$status"
done

echo "---------------------------------------------------------------"
echo "총 서비스: $TOTAL_SVCS | PASS: $PASS_SVCS | 총 테스트: $TOTAL_TESTS"
