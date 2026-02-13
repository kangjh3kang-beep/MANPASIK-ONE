# Phase 1D: 통합 MVP (E2E + 배포)

> **전제**: Phase 1C 완료 (Flutter + Go 4서비스 + E2E 4서비스 헬스체크)
> **목표**: E2E 비즈니스 플로우 검증, CI/배포 정리, Phase 1D Gate 통과

---

## 1. 범위

| Phase 1D 항목 | 내용 | Gate |
|---------------|------|------|
| **D1** | E2E 비즈니스 플로우 | Login → StartSession → EndSession → GetHistory (Go gRPC) |
| **D2** | CI/실행 경로 정리 | E2E를 backend 모듈에서 실행, Makefile/CI 일치 |
| **D3** | Phase 1D Gate·문서 | QUALITY_GATES, CHANGELOG, CONTEXT 갱신 |

---

## 2. 현재 상태

- **E2E**: `tests/e2e/service_test.go` — 4서비스 헬스체크만 있음. `TestMeasurementFlow`는 Skip.
- **CI**: integration-test job이 `cd tests/e2e` 후 `go test` 실행(루트에 go.mod 없으면 실패 가능).
- **Makefile**: `test-integration` → `cd tests/e2e && go test`.
- **gRPC 클라이언트**: backend shared gen에는 Server만 있고 Client 미생성 → E2E에서 `conn.Invoke` 또는 backend 내 e2e로 통합.

---

## 3. 구현 방침

- **E2E 위치**: backend 모듈 내 `backend/tests/e2e`에 플로우 테스트 추가. 기존 `tests/e2e` 헬스체크는 backend로 이전하거나, E2E 실행을 backend 기준으로 통일.
- **플로우 테스트**: Auth Login → (user_id 추출) → Measurement StartSession → EndSession → GetMeasurementHistory. 서비스 미기동 시 스킵 가능(환경변수 또는 build tag).
- **CI**: `working-directory: backend`, `go test ./tests/e2e/...` (필요 시 `-short`로 헬스/플로우 스킵 옵션).

---

## 4. 구현 순서

### Step 1: backend/tests/e2e 구성

- `backend/tests/e2e/health_test.go`: 기존 4서비스 헬스체크 이전(또는 복사).
- `backend/tests/e2e/flow_test.go`: Login → StartSession → EndSession → GetHistory, 서비스 연결 실패 시 `t.Skip`.

### Step 2: CI·Makefile 정리

- Makefile `test-integration`: `cd backend && go test -v -tags=integration ./tests/e2e/...`
- CI integration-test: `working-directory: backend`, `go test -v ./tests/e2e/...` (필요 시 `-short` 또는 서비스 기동 후 실행).

### Step 3: Phase 1D Gate·문서

- QUALITY_GATES: Phase 1D 행·Stage D1 체크리스트.
- CHANGELOG, CONTEXT 갱신.

---

## 5. Phase 1D Gate 통과 기준

- [ ] E2E 플로우 테스트 1개: Login → StartSession → EndSession → GetHistory (서비스 기동 시 통과 또는 스킵 정책 명시)
- [ ] E2E 실행 경로: backend 모듈 기준으로 통일, Makefile/CI 반영
- [ ] go test ./... (backend) 통과
- [ ] CHANGELOG, CONTEXT, QUALITY_GATES 갱신

---

## 6. 참조

- E2E 기존: `tests/e2e/service_test.go`
- Proto: `backend/shared/proto/manpasik.proto`, `backend/shared/gen/go/v1`
- CI: `.github/workflows/ci.yml`

---

**문서 버전**: 1.0  
**최종 업데이트**: 2026-02-10
