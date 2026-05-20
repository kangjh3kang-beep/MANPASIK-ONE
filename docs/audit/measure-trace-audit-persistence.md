# Measure Trace Audit Persistence Audit

**Date**: 2026-05-01  
**Owner**: Codex  
**Scope**: Gateway Measure trace intake -> audit-service persisted audit entry

## 목적

이전 단계에서 Flutter Measure golden path trace event는 Gateway의 `POST /api/v1/measurements/trace-events`까지 도달했다. 다만 Gateway 응답의 `audit_status`는 audit-service write RPC 부재 때문에 deferred 상태였다. 이번 단계는 proto/gRPC 대규모 변경 없이 audit-service의 기존 `AuditRepository.Store` 경로를 재사용하는 전용 HTTP write intake를 추가해 trace event를 실제 감사 저장소에 남긴다.

## 구현 계약

### audit-service

- `POST /audit/events`
  - `action`, `resource_type` 필수
  - `occurred_at`은 RFC3339/RFC3339Nano 형식
  - 정상 처리 시 `201 Created`
  - 응답: `accepted=true`, `status=persisted`, `entry_id`
- `AuditService.RecordActionWithMetadata`
  - `admin_id`, `action`, `resource_type`, `resource_id`, `description`, `ip_address`, `user_agent`, `metadata`, `timestamp`를 받아 기존 repository에 저장한다.
  - 기존 `RecordAction`은 동일 저장 경로를 호출하므로 기존 호출부와 호환된다.

### Gateway

- `AUDIT_INTAKE_URL` 환경변수가 설정되면 `HTTPAuditEventRecorder`를 등록한다.
- Measure trace event 수신 시 audit event로 변환한다.
  - `action`: `measure.trace.<phase>`
  - `resource_type`: `measurement_trace`
  - `resource_id`: `session_id`, 없으면 Gateway `event_id`
  - `metadata`: schema/source/route/phase/elapsed/session/cartridge/engine/unit/confidence/has_primary_value/failure/diagnostic
  - `primary_value`는 여전히 저장 및 전달 금지
- audit-service 저장 성공 시 Gateway 응답:
  - `audit_status=persisted`
  - `audit_entry_id=<entry_id>`
- audit intake가 없거나 실패하면 측정 trace 수신 자체는 `202 Accepted`로 유지하고 `audit_status`로 상태를 노출한다.

## 변경 파일

- `backend/services/audit-service/internal/service/audit.go`
- `backend/services/audit-service/internal/handler/http.go`
- `backend/services/audit-service/internal/handler/http_test.go`
- `backend/services/audit-service/cmd/main.go`
- `backend/services/gateway/internal/handler/audit_recorder.go`
- `backend/services/gateway/internal/handler/audit_recorder_test.go`
- `backend/services/gateway/internal/handler/rest_handler.go`
- `backend/services/gateway/internal/handler/measurement_routes.go`
- `backend/services/gateway/internal/handler/e2e_test.go`
- `backend/services/gateway/cmd/main.go`
- `docker-compose.yml`
- `infrastructure/docker/docker-compose.dev.yml`
- `infrastructure/kubernetes/base/config/configmap.yaml`

## 품질 게이트

```bash
cd /home/kangjh3kang/Manpasik
export PATH=/home/kangjh3kang/sdk/go-go1.26.2/bin:/usr/local/bin:/usr/bin:/bin
gofmt -w backend/services/audit-service/internal/service/audit.go \
  backend/services/audit-service/internal/handler/http.go \
  backend/services/audit-service/internal/handler/http_test.go \
  backend/services/audit-service/cmd/main.go \
  backend/services/gateway/internal/handler/audit_recorder.go \
  backend/services/gateway/internal/handler/audit_recorder_test.go \
  backend/services/gateway/internal/handler/rest_handler.go \
  backend/services/gateway/internal/handler/measurement_routes.go \
  backend/services/gateway/internal/handler/e2e_test.go \
  backend/services/gateway/cmd/main.go
go test ./backend/services/audit-service/... ./backend/services/measurement-service/... ./backend/services/gateway/...
docker compose -f docker-compose.yml config --quiet
docker compose -f infrastructure/docker/docker-compose.dev.yml config --quiet
```

결과:

- Go 테스트: PASS
- root Docker Compose config: PASS
- dev Docker Compose config: PASS
- `kubectl kustomize infrastructure/kubernetes/base`: 기존 namespace 변환 충돌로 BLOCKED
- `kubectl apply --dry-run=client` ConfigMap 단독 검증: 로컬 Kubernetes API 미연결로 BLOCKED

## 다음 단계

- 배포 설정 반영: root Docker Compose, dev Docker Compose, K8s base ConfigMap에 `AUDIT_INTAKE_URL=http://audit-service:9100/audit/events` 계열 값을 명시했다.
- dev Docker Compose에는 audit-service 컨테이너와 `39-audit.sql` init 마운트를 추가해 clean DB에서 `audit_entries` 테이블이 생성되도록 했다.
- H2에서 `docs/audit/mock-retirement-register.md`의 Measure 관련 mock/stub 항목을 native Rust/DB 연결 상태와 대조한다.
