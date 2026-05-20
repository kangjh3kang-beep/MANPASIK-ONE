#!/bin/bash
export PATH=/usr/local/go/bin:/usr/bin:/bin
cd /home/kangjh3kang/Manpasik/backend

echo "=== Go 전체 테스트 ==="
GOWORK=off go test -count=1 -timeout 180s $(GOWORK=off go list ./services/.../service/... ./services/.../webrtc/... 2>/dev/null) 2>&1 | tail -80

echo ""
echo "=== 종합 ==="
GOWORK=off go test -count=1 -timeout 180s $(GOWORK=off go list ./... 2>/dev/null | grep -v 'fhir-adapter\|hl7-parser') 2>&1 | tail -80
