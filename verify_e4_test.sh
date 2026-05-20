#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik/backend

echo "=== 전체 서비스 테스트 ==="
PKGS=$(GOWORK=off go list ./services/.../service/... 2>/dev/null)
echo "패키지 수: $(echo "$PKGS" | wc -l)"
GOWORK=off go test -count=1 -timeout 180s $PKGS 2>&1 | tail -60
