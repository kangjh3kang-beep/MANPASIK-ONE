# Phase 1C Stage S6: 전체 통합

> **전제**: S5a·S5b 완료 (gRPC, Flutter 화면, 차트, BLE/NFC 스텁)
> **목표**: Docker Compose 기동 검증, E2E/통합 테스트, Phase 1C 완료 선언

---

## 1. 현재 상태

- **S5a**: Go 4서비스 Docker, Flutter gRPC·Repository·화면, 60+ 테스트
- **S5b**: fl_chart 결과·트렌드 화면, BLE/NFC UI(스텁), /measurement/result
- **E2E**: `tests/e2e/service_test.go` — auth/measurement 헬스체크 존재, user/device 미포함
- **README**: Go gRPC 포트(50051–50054) 미기재

---

## 2. S6 범위

| 순서 | 항목 | 내용 | Gate |
|------|------|------|------|
| 1 | **Docker Compose 검증** | Go 4서비스 + Postgres 기동 확인, README·서비스 목록에 gRPC 포트 반영 | 문서·실행 가능 |
| 2 | **통합/E2E 테스트** | 4서비스 헬스체크 E2E, (선택) Flutter 측정 플로우 위젯 테스트 | 테스트 통과 |
| 3 | **S6 Gate·Phase 1C 완료** | QUALITY_GATES S6 체크리스트, Phase 1C 완료 선언, CHANGELOG/CONTEXT 갱신 | 문서 갱신 |

---

## 3. 구현 순서

### Step 1: Docker Compose 기동 검증 + README 갱신

- **README.md**: 서비스 접속 표에 Go gRPC 4서비스 추가 (auth 50051, user 50052, device 50053, measurement 50054)
- **(선택)** `infrastructure/docker/README.md` 또는 Makefile 타깃: `make dev-go` — postgres + 4 Go 서비스만 up (빠른 검증용)
- 검증: `docker compose -f docker-compose.dev.yml up -d postgres auth-service user-service device-service measurement-service` 후 포트 리스닝 확인

### Step 2: 통합/E2E 테스트

- **tests/e2e/service_test.go**: user-service(50052), device-service(50053) 헬스체크 추가
- **Flutter**: 기존 위젯 테스트 유지, (선택) 로그인 → 홈 → 측정 결과 화면 진입 플로우 1개 추가
- CI에서 E2E는 서비스 기동 후 실행 또는 `-short`로 스킵 가능하도록 유지

### Step 3: S6 Gate·Phase 1C 완료

- **QUALITY_GATES.md**: S6 행 "통과" 처리, S6 Gate 체크리스트 추가, Phase 1C "완료" 반영
- **CHANGELOG.md**: S6 작업 로그
- **CONTEXT.md**: Phase 1C 완료, 다음 Phase 1D 준비

---

## 4. S6 Gate 통과 기준

- [ ] Docker Compose로 Postgres + Go 4서비스 기동 가능
- [ ] E2E 테스트: 4서비스 헬스체크 통과 (또는 서비스 미기동 시 스킵)
- [ ] README에 gRPC 포트 50051–50054 명시
- [ ] flutter analyze 0, flutter test 통과
- [ ] go test ./... 통과
- [ ] CHANGELOG·CONTEXT·QUALITY_GATES 갱신

---

## 5. 핵심 파일 참조

- E2E: `tests/e2e/service_test.go`
- Docker Compose: `infrastructure/docker/docker-compose.dev.yml`
- README: `README.md`

---

**문서 버전**: 1.0  
**최종 업데이트**: 2026-02-10
