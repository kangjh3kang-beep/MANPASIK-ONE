# ManPaSik 테스트 전략 명세서 (Test Strategy Specification)

**문서번호**: MPK-TEST-STRATEGY-v1.0  
**갱신일**: 2026-02-12  
**목적**: Phase별 테스트 범위, 커버리지 목표, 자동화 전략, 도구 체계를 정의  
**참조**: QUALITY_GATES.md, docs/compliance/vnv-master-plan.md

---

## 1. 테스트 피라미드

```
                    ┌──────────────┐
                    │   E2E 테스트   │  ← 적은 수, 높은 가치
                    │  (5~50개)     │
                   ┌┴──────────────┴┐
                   │  통합 테스트     │  ← 서비스 간 연동
                   │  (50~200개)    │
                  ┌┴────────────────┴┐
                  │   단위 테스트      │  ← 많은 수, 빠른 실행
                  │  (500~2000+개)   │
                 ┌┴──────────────────┴┐
                 │  정적 분석 (Lint)    │  ← 매 PR, 자동
                 └────────────────────┘
```

---

## 2. Phase별 테스트 목표

### 2.1 Phase 1 (MVP) — 현재 위치

| 테스트 유형 | 목표 커버리지 | 현재 | 갭 | 자동화 |
|-----------|------------|------|-----|--------|
| **Rust 단위 테스트** | 80% | ~75% (62개 테스트) | +10~15개 | ✅ `cargo test` |
| **Go 단위 테스트** | 80% | ~50% (32개 테스트) | +40~60개 | ✅ `go test` |
| **Flutter 단위 테스트** | 60% | 0% (0개 테스트) | +60~80개 | 🔲 `flutter test` |
| **Go 통합 테스트** | 4서비스 연동 | 있음 (E2E 파일 8개) | 보강 필요 | ⚠️ 부분 |
| **E2E 테스트** | 5개 시나리오 | 있음 (flow_test.go) | 시나리오 확장 | ⚠️ 부분 |
| **Lint (정적 분석)** | 100% 통과 | ✅ CI 적용 | — | ✅ |
| **보안 스캔** | 0 Critical | 미적용 | 도구 도입 | 🔲 |
| **성능 테스트** | 500 RPS | 미실행 | 도구 도입 | 🔲 |

**Phase 1 Gate 통과 조건:**
- [ ] Rust 단위 테스트 75개+ (커버리지 80%)
- [ ] Go 단위 테스트 70개+ (4서비스 각 핸들러/서비스/레포지토리)
- [ ] Flutter 단위 테스트 60개+ (6 Feature 각 Provider/Widget)
- [ ] E2E 5개 시나리오 통과 (Login→Measure→Result→History→Logout)
- [ ] Lint 전체 통과 (Rust clippy + Go golangci-lint + Dart analyze)

### 2.2 Phase 2 (Core)

| 테스트 유형 | 목표 커버리지 | 테스트 수 | 자동화 |
|-----------|------------|---------|--------|
| **Rust 단위 테스트** | 85% | 80개+ | ✅ CI |
| **Go 단위 테스트** | 80% | 150개+ (11서비스) | ✅ CI |
| **Flutter 단위/위젯 테스트** | 70% | 120개+ (10 Feature) | ✅ CI |
| **Go 통합 테스트** | 11서비스 연동 | 30개+ | ✅ CI (Docker) |
| **E2E 테스트** | 15개 시나리오 | 15개+ | ✅ CI |
| **보안 스캔** | 0 Critical, 0 High | 매 PR | ✅ CI |
| **성능 테스트** | 5,000 RPS | 3개 시나리오 | ✅ Nightly |
| **계약 테스트 (Contract)** | Kafka 이벤트 전체 | 18개 토픽 | ✅ CI |

**Phase 2 Gate 통과 조건:**
- [ ] 전체 커버리지 80% 이상
- [ ] E2E 15개 시나리오 전체 PASS
- [ ] 성능: 5,000 RPS에서 P95 < 200ms
- [ ] 보안: Critical/High 취약점 0건
- [ ] Kafka 이벤트 계약 테스트 전체 PASS

### 2.3 Phase 3 (Advanced)

| 테스트 유형 | 목표 커버리지 | 테스트 수 | 자동화 |
|-----------|------------|---------|--------|
| **Go 단위 테스트** | 80% | 300개+ (21서비스) | ✅ CI |
| **Flutter 단위/위젯/통합 테스트** | 75% | 200개+ (12 Feature) | ✅ CI |
| **Go 통합 테스트** | 21서비스 연동 | 60개+ | ✅ CI (Docker) |
| **E2E 테스트** | 30개 시나리오 | 30개+ | ✅ CI |
| **성능 테스트** | 50,000 RPS | 5개 시나리오 | ✅ Nightly |
| **카오스 테스트** | 5개 장애 시나리오 | 5개 | ✅ Weekly |
| **오프라인 검증** | 72시간 연속 | 1개 (REQ-065) | 수동 + 자동화 |
| **접근성 테스트** | WCAG AA | 12개 화면 | ✅ CI |

**Phase 3 Gate 통과 조건:**
- [ ] 전체 커버리지 80% 이상
- [ ] E2E 30개 시나리오 전체 PASS
- [ ] 성능: 50,000 RPS에서 P95 < 200ms
- [ ] 카오스: 5개 장애 시나리오에서 RTO < 5분
- [ ] 오프라인 72시간 연속 동작 검증 PASS

### 2.4 Phase 4 (Ecosystem)

| 테스트 유형 | 목표 커버리지 | 테스트 수 |
|-----------|------------|---------|
| Go 단위 테스트 | 80% | 500개+ (29서비스) |
| Flutter 테스트 | 80% | 300개+ |
| E2E 테스트 | 50개 시나리오 | 50개+ |
| 성능 테스트 | 200,000 RPS | 10개 시나리오 |
| 카오스 테스트 | 10개 장애 시나리오 | 10개 |
| 보안 침투 테스트 | OWASP Top 10 | 연 1회 |
| 규정 적합성 | IEC 62304 V&V | 전체 REQ 추적 |

---

## 3. 테스트 유형별 상세

### 3.1 단위 테스트 (Unit Test)

| 언어 | 프레임워크 | 명명 규칙 | 실행 명령 |
|------|----------|---------|----------|
| Rust | `#[cfg(test)]` + criterion | `test_기능_시나리오_기대결과` | `cargo test` |
| Go | testing + testify | `Test서비스명_메서드명_시나리오` | `go test ./...` |
| Dart | flutter_test + mocktail | `test('기능 설명', ...)` | `flutter test` |

**필수 테스트 대상:**
- 모든 public 함수/메서드
- 모든 gRPC 핸들러 (성공 + 에러 케이스)
- 모든 서비스 레이어 비즈니스 로직
- 모든 레포지토리 CRUD
- 모든 Provider (Riverpod StateNotifier)
- 모든 Widget (렌더링 + 상호작용)

### 3.2 통합 테스트 (Integration Test)

**범위:** 서비스 ↔ DB, 서비스 ↔ 서비스, 서비스 ↔ 외부 시스템

| 테스트 대상 | 도구 | 환경 |
|-----------|------|------|
| Go 서비스 ↔ PostgreSQL | testcontainers-go | Docker |
| Go 서비스 ↔ Redis | testcontainers-go | Docker |
| Go 서비스 ↔ Kafka | testcontainers-go | Docker |
| Go 서비스 ↔ Milvus | testcontainers-go | Docker |
| Flutter ↔ gRPC | mock gRPC server | 로컬 |

### 3.3 E2E 테스트 (End-to-End)

**시나리오 목록 (Phase별 누적):**

| # | 시나리오 | Phase | 서비스 범위 |
|---|---------|-------|-----------|
| E2E-001 | 회원가입 → 로그인 → 프로필 조회 | 1 | auth, user |
| E2E-002 | 디바이스 등록 → 목록 조회 → 상태 확인 | 1 | device |
| E2E-003 | 측정 세션 시작 → 데이터 스트림 → 종료 → 이력 조회 | 1 | measurement |
| E2E-004 | 전체 플로우: Login→StartSession→EndSession→GetHistory | 1 | auth, measurement |
| E2E-005 | 서비스 헬스체크 (4서비스) | 1 | 전체 |
| E2E-006 | 구독 생성 → 업그레이드 → 카트리지 접근 확인 | 2 | subscription, cartridge |
| E2E-007 | 상품 조회 → 장바구니 → 주문 → 결제 | 2 | shop, payment |
| E2E-008 | 측정 완료 → AI 추론 → 코칭 리포트 | 2 | measurement, ai-inference, coaching |
| E2E-009 | 카트리지 인증 → 보정 → 측정 | 2 | cartridge, calibration, measurement |
| E2E-010 | 위험 감지 → 알림 발송 | 2 | ai-inference, notification |
| E2E-011 | 가족 그룹 생성 → 구성원 초대 → 모니터링 | 3 | family |
| E2E-012 | 화상진료 예약 → 상담 → 처방 | 3 | reservation, telemedicine, prescription |
| E2E-013 | 커뮤니티 게시글 → 댓글 → 번역 | 3 | community, translation |
| E2E-014 | 건강 기록 생성 → FHIR 내보내기 | 3 | health-record |
| E2E-015 | 관리자 대시보드 → 사용자 관리 → 권한 설정 | 3 | admin |

### 3.4 성능 테스트 (Performance Test)

| 도구 | 용도 | 시나리오 |
|------|------|---------|
| **k6** | HTTP/gRPC 부하 테스트 | 동시 사용자 시뮬레이션 |
| **ghz** | gRPC 전용 벤치마크 | RPC별 응답시간/처리량 |
| **criterion** | Rust 마이크로벤치마크 | 차동측정, 핑거프린트, 암호화 |

**Phase 2 성능 테스트 시나리오:**
1. 로그인 부하: 1,000 CCU, 5분 지속 → P95 < 300ms
2. 측정 스트림: 500 동시 세션, 10분 지속 → 메시지 지연 < 50ms
3. 벡터 검색: 100 동시 쿼리, 100만 벡터 → P95 < 200ms

### 3.5 보안 테스트 (Security Test)

| 도구 | 대상 | 빈도 |
|------|------|------|
| **gosec** | Go 소스 코드 | 매 PR |
| **cargo-audit** | Rust 의존성 | 매 PR |
| **trivy** | Docker 이미지 | 매 빌드 |
| **OWASP ZAP** | API 엔드포인트 | 주 1회 |
| **Bandit** | Python (ai-ml) | 매 PR (Phase 4+) |

---

## 4. 테스트 자동화 파이프라인

```
PR 생성
  ├── Lint (clippy, golangci-lint, dart analyze)
  ├── 단위 테스트 (cargo test, go test, flutter test)
  ├── 보안 스캔 (gosec, cargo-audit, trivy)
  └── 커버리지 리포트 → PR 코멘트

Merge to main
  ├── 통합 테스트 (testcontainers)
  ├── Docker 이미지 빌드
  └── E2E 테스트 (Staging 환경)

Nightly (매일 새벽)
  ├── 성능 테스트 (k6)
  └── 보안 전체 스캔 (OWASP ZAP)

Weekly
  └── 카오스 테스트 (Phase 3+)
```

---

## 5. 커버리지 측정 및 리포팅

| 언어 | 커버리지 도구 | 리포팅 |
|------|------------|--------|
| Rust | `cargo tarpaulin` | CI 아티팩트 |
| Go | `go test -cover` + `go tool cover` | CI 아티팩트 |
| Dart | `flutter test --coverage` + `lcov` | CI 아티팩트 |

**커버리지 게이트 (CI 차단 조건):**
- Phase 1: 전체 60% 미만 시 PR 차단
- Phase 2+: 전체 75% 미만 시 PR 차단
- Phase 3+: 전체 80% 미만 시 PR 차단

---

**참조**: QUALITY_GATES.md, .github/workflows/ci.yml, docs/compliance/vnv-master-plan.md
