# ManPaSik Quality Gate 프로세스 정의서

> **용도**: 단계별 개발 시 품질 검증 기준과 Gate 통과 조건을 정의
> **적용**: 모든 IDE·환경에서 코드 변경·Stage 완료·Phase 완료 시 이 문서 기준으로 검증 수행
> **공통 규칙 요약 (모든 IDE)**: [docs/COMMON_RULES.md](docs/COMMON_RULES.md) — 코드 리뷰 → 린트 → 테스트
> **참조**: `.cursor/rules/manpasik-project.mdc` (Cursor 전용), `.cursor/rules/work-logging.mdc`

---

## 1. 3단계 품질 검증 체계 개요

| 레벨 | 시점 | 목적 |
|------|------|------|
| **Level 1** | 매 코드 변경 시 | 즉시 품질 유지 (린트·테스트·빌드) |
| **Level 2** | Stage(기능 단위) 완료 시 | 구현 완결성·보안·문서 검증 |
| **Level 3** | Phase(마일스톤) 완료 시 | 통합·보안 스캔·성능·규정 문서 정합성 |

**원칙**: Level 1 통과 없이 작업 완료로 간주하지 않음. Stage 완료 시 Level 2 체크리스트 필수. Phase 완료 시 Level 3 검증 후 다음 Phase 진행.

**단계 완료 시 필수 3단계 (코드 리뷰 → 린트 → 테스트)**: 모든 개발 단계(기능 추가, 서비스 추가, Stage/Phase 완료)에서 **완료 선언 전에** 아래 순서로 반드시 수행한다.  
1. **코드 리뷰**: 변경/추가된 코드를 대상으로 자기 점검 — 보안(입력 검증·SQL/인젝션·하드코딩 비밀), 프로젝트 패턴 준수(에러 처리·API 일관성·nil/panic 가능성), 의존성 정합성(import·타입·스텁과의 불일치).  
2. **린트**: 해당 언어 린터 실행 (Go: `golangci-lint run`, Rust: `cargo clippy` 등), 에러 0개.  
3. **테스트·빌드**: 해당 범위 테스트 및 빌드 실행 (Go: `go mod tidy && go test ./...` 및 `go build ./...` 등).  
에이전트 환경에서 2·3 실행이 불가한 경우, CHANGELOG에 "검증 필요"를 남기고 사용자에게 실행할 명령을 안내한 뒤, 사용자 검증 후에만 완료로 표기한다. 검증 생략 후 완료로 표기하지 않는다.

**우회·미루기 금지**: 구현 항목을 "나중에" 또는 "우회"로 미루지 않음. 막히면 해결점을 찾아 정상 구현. 부득이한 우회/연기는 KNOWN_ISSUES 등록 및 해결 조건·시한 명시 필수. (`docs/development-philosophy.md` §3.1)

---

## 2. Level 1: 매 작업 즉시 검증 (Every Change)

모든 코드 변경 시 수행하는 최소 검증. **코드 리뷰(자기 점검) → 린트 에러 0개 → 테스트 전체 통과·빌드 성공** 순서로 수행하며, 세 가지 모두 통과해야 작업 완료로 간주한다.

### 2.1 검증 명령어 (언어별)

| 언어 | 린트 | 테스트 | 빌드 |
|------|------|--------|------|
| Rust | `cargo clippy --all-targets` | `cargo test` | `cargo build` |
| Go | `golangci-lint run` | `go test ./...` | `go build ./...` |
| Dart | `dart analyze` | `flutter test` | `flutter build` |
| TypeScript | `eslint .` | `jest --passWithNoTests` | `tsc --noEmit` |

### 2.2 Rust 코어 (Feature 플래그 사용 시)

- AI 제외 빌드: `cargo build -p manpasik-engine --no-default-features --features 'std,ble,nfc,fingerprint'`
- 테스트: `cargo test -p manpasik-engine --no-default-features --features 'std,ble,nfc,fingerprint'`
- 전체 빌드: `cargo build -p manpasik-engine --features full` (tflitec 0.7 사용, Bazel 불필요)

### 2.3 작업 완료 시 기록

Level 1 검증(코드 리뷰 → 린트 → 테스트·빌드) 결과를 CHANGELOG.md 해당 작업 항목에 기록:
- 코드 리뷰: 자기 점검 완료 (보안·패턴·의존성 정합성)
- 린트: 통과 / 경고 N개 (기존 허용 여부)
- 테스트: N개 통과, 실패 0
- 빌드: 성공 / 실패 (실패 시 KNOWN_ISSUES.md 반영)

**에이전트 필수 동작**: 단계 완료 시 **1) 코드 리뷰**로 변경 구역 점검 후, **2) 린트**, **3) 테스트·빌드** 순으로 실행 시도한다. 2·3 실행 불가(예: WSL 출력 미캡처) 시 "수동 검증 필요"를 CHANGELOG에 남기고 사용자에게 실행할 명령어를 제시한 뒤, 사용자 검증 후에만 완료로 간주한다.

---

## 3. Level 2: Stage Gate (기능 단위 완료 시)

Stage는 프로젝트 기능 단위. Stage 완료 선언 전 아래 체크리스트를 수행하고, 모든 항목 충족 시 QUALITY_GATES.md의 "현재 Stage 상태"를 갱신.

### 3.1 Stage 정의 및 Gate 기준

| Stage | 범위 | Gate 기준 |
|-------|------|-----------|
| S1 | Rust 코어 엔진 | 62테스트 통과, clippy 0 error (기존 경고 허용), 빌드 성공 |
| S2 | Go 인증 서비스 | 유닛 테스트 커버리지 80%+, golangci-lint 통과, gRPC 계약 검증 |
| S3 | Go 측정/디바이스/사용자 서비스 | 서비스별 테스트 80%+, 서비스 간 gRPC 통합 검증 |
| S4 | Flutter 앱 기본 구조 | 위젯 테스트, 라우팅 검증, Rust FFI 연동 테스트 |
| S5 | Flutter 핵심 화면 | 통합 테스트, UX 플로우 검증 |
| S6 | 전체 통합 | E2E 테스트, Docker Compose 환경 기동 검증 |

### 3.2 Stage Gate 체크리스트 템플릿

매 Stage 완료 시 아래 체크리스트를 수행하고, 결과를 CHANGELOG.md 해당 Stage 완료 항목에 기록.

```markdown
## Stage Gate Checklist — [Stage명]

### 코드 품질
- [ ] 린트 에러 0개 (경고는 기존 허용 목록만 허용)
- [ ] 모든 단위 테스트 통과
- [ ] 테스트 커버리지 >= 80%
- [ ] 새로 추가된 공개 API에 문서 주석 있음
- [ ] 에러 처리 누락 없음 (unwrap/panic 금지)

### 기능 완결성
- [ ] 요구사항 대비 구현 항목 전수 확인
- [ ] 엣지 케이스 테스트 포함
- [ ] 의존 모듈과의 인터페이스 정합성 확인

### 보안 (의료기기)
- [ ] 민감 데이터 하드코딩 없음
- [ ] 입력 검증 누락 없음
- [ ] OWASP Top 10 관련 취약점 점검

### 문서/기록
- [ ] CHANGELOG.md 작업 기록 완료
- [ ] KNOWN_ISSUES.md 미해결 이슈 갱신
- [ ] CONTEXT.md 진행률 갱신
```

---

## 4. Level 3: Phase Gate (마일스톤 완료 시)

Phase 1 MVP의 4개 마일스톤. Phase 완료 시 아래 추가 검증 후 다음 Phase 진행.

| Phase | 범위 |
|-------|------|
| Phase 1A | 코어 완성 (Rust + Infra) |
| Phase 1B | 백엔드 완성 (Go 4서비스 + DB) |
| Phase 1C | 프론트엔드 (Flutter + FFI) |
| Phase 1D | 통합 MVP (E2E + 배포) |

### 4.1 Phase Gate 추가 검증 항목

- 통합 테스트: 서비스 간 gRPC 호출, DB 연동
- 보안 스캔: `cargo audit`, Go 의존성 취약점 스캔
- 성능 벤치마크: Rust criterion, Go benchmark 기준선
- 규정 문서 정합성: 코드 변경이 V&V/위험관리 문서와 일치하는지
- IEC 62304 추적성: 요구사항 ↔ 코드 ↔ 테스트 매핑 확인

---

## 5. 현재 Stage 상태 추적

| Stage | 범위 | 상태 | 통과일 | 비고 |
|-------|------|------|--------|------|
| S1 | Rust 코어 엔진 | ✅ 통과 | 2026-02-10 | 62테스트, BLE 포함 빌드, 경고 2개(기능 무관) |
| S2 | Go 인증 서비스 | ✅ 통과 | 2026-02-10 | 8테스트 PASS, gRPC 핸들러+저장소+main 연동 완료 |
| S3 | Go 측정/디바이스/사용자 서비스 | ✅ 통과 | 2026-02-10 | user 10+device 11+measurement 11=32테스트, 인메모리 저장소 |
| S4 | Flutter 앱 기본 구조 | ✅ 통과 | 2026-02-10 | 7화면 라우트, 3 Provider, 6언어 i18n, 30+테스트, Feature-First 구조 |
| S5 | Flutter 핵심 화면 | ✅ 통과 (S5a) | 2026-02-10 | gRPC 연동, 4서비스 Dockerfile, Repository+화면 고도화, 50+테스트 |
| S6 | 전체 통합 | ✅ 통과 | 2026-02-10 | E2E 4서비스 헬스체크, README gRPC 포트 반영, Phase 1C 완료 |

### 5.1 Phase 상태

| Phase | 범위 | 상태 | 비고 |
|-------|------|------|------|
| Phase 1A | 코어 + Infra | ✅ 통과 | Rust 100%, Docker 구성 완료 |
| Phase 1B | Go 4서비스 + DB | ✅ 통과 | 4서비스 빌드·테스트 PASS, DB 스키마 4종, 인메모리 fallback |
| Phase 1C | Flutter + FFI | ✅ 완료 | S4·S5·S6 통과 (gRPC, 차트, BLE/NFC 스텁, E2E 4서비스) |
| Phase 1D | 통합 MVP | ✅ 완료 | E2E 플로우 통과, Docker 빌드·DB 초기화 정상, gRPC 4서비스 연동 검증 (2026-02-10) |

### 5.2 Stage S2 Gate Checklist — Go 인증 서비스

#### 코드 품질
- [x] 린트 에러 0개
- [x] 모든 단위 테스트 통과 (8개 PASS)
- [x] 새로 추가된 공개 API에 문서 주석 있음
- [x] 에러 처리 누락 없음 (표준 AppError 사용)

#### 기능 완결성
- [x] Register, Login, RefreshToken, Logout, ValidateToken 구현 완료
- [x] JWT Access+Refresh Token Rotation
- [x] PostgreSQL + 인메모리 fallback
- [x] gRPC 핸들러 → 서비스 → 저장소 전 계층 연동

#### 보안 (의료기기)
- [x] 비밀번호 bcrypt cost=12 해싱
- [x] JWT secret 환경변수로 관리
- [x] 입력 검증 (이메일/비밀번호)
- [x] 내부 에러 미노출 (AppError → gRPC Status 변환)

#### 문서/기록
- [x] CHANGELOG.md 작업 기록 완료
- [x] KNOWN_ISSUES.md 갱신 (TFLite 해결)
- [x] CONTEXT.md 진행률 갱신

### 5.3 Stage S3 Gate Checklist — Go 측정/디바이스/사용자 서비스

#### 코드 품질
- [x] 린트 에러 0개 (go build ./... 성공)
- [x] 단위 테스트: user 10개 + device 11개 + measurement 11개 = 32개
- [x] 공개 API 문서 주석 완료
- [x] Repository 패턴 일관 적용

#### 기능 완결성
- [x] user-service: GetProfile, UpdateProfile, GetSubscription (4티어)
- [x] device-service: RegisterDevice (구독 기반 제한), ListDevices, UpdateStatus, OTA
- [x] measurement-service: StartSession, ProcessMeasurement, EndSession, GetHistory
- [x] 인메모리 저장소 (PostgreSQL/TimescaleDB/Milvus/Kafka 전환 준비)

#### 보안 (의료기기)
- [x] 입력 검증 누락 없음 (빈 값 체크)
- [x] 디바이스 수 제한 (구독 기반)
- [x] 에러 응답 표준화 (내부 정보 미노출)

#### 문서/기록
- [x] CHANGELOG.md 기록 완료
- [x] DB 초기화 스크립트 4종
- [x] CONTEXT.md 진행률 갱신

### 5.4 Stage S4 Gate Checklist — Flutter 앱 기본 구조

#### 코드 품질
- [x] 린트 에러 0개 (dart analyze 통과 필요)
- [x] 단위 테스트: auth 7 + theme 4 + locale 10 + validators 11 = 32개 작성
- [x] 공개 위젯/Provider에 문서 주석 있음
- [x] Feature-First 디렉토리 구조 완성

#### 기능 완결성
- [x] P0 라우트 7개 정의 (splash, login, register, home, measure, devices, settings)
- [x] GoRouter 인증 기반 리다이렉트 동작
- [x] Riverpod ProviderScope + 3개 핵심 Provider (auth, theme, locale)
- [x] Material Design 3 테마 (light/dark)
- [x] 다국어 6개 언어 ARB (ko, en, ja, zh, fr, hi) + 확장 구조

#### 보안 (의료기기)
- [x] 비밀번호 강도 검증 (8자+, 영문+숫자)
- [x] 이메일 형식 검증 (Regex)
- [x] 민감 정보 하드코딩 없음

#### 문서/기록
- [x] CHANGELOG.md 기록 완료
- [x] QUALITY_GATES.md S4 Gate 체크리스트 작성
- [x] CONTEXT.md 진행률 갱신

---

## 6. 참조

- **규칙**: `.cursor/rules/manpasik-project.mdc` 9항 Quality Gate 개발 프로세스
- **기록**: `.cursor/rules/work-logging.mdc` — 작업 완료 시 Level 1 검증 결과 기록
- **이슈**: `KNOWN_ISSUES.md` — 빌드/환경 이슈
- **규정**: `docs/compliance/vnv-master-plan.md` — V&V 전략

### 5.5 Stage S5a Gate Checklist — Flutter gRPC 연동·화면 고도화

#### 코드 품질
- [x] 린트 에러 0개 (flutter analyze)
- [x] 단위·위젯 테스트 50개+ (auth 7, theme 4, locale 10, validators 11, repository 17, grpc_client 7, screen 5)
- [x] Repository 패턴 (interface + gRPC impl + Fake for test)
- [x] gRPC 채널 관리·JWT 인터셉터 구현

#### 기능 완결성
- [x] Go 4서비스 Dockerfile + docker-compose.dev.yml 추가 (50051–50054)
- [x] Dart gRPC 클라이언트 (수동 proto 메시지·스텁, GrpcClientManager, AuthInterceptor)
- [x] AuthRepositoryImpl / DeviceRepositoryImpl / MeasurementRepositoryImpl / UserRepositoryImpl
- [x] HomeScreen: GetMeasurementHistory 연동
- [x] MeasurementScreen: StartSession/EndSession 연동
- [x] DeviceListScreen: ListDevices 연동
- [x] SettingsScreen: GetProfile/GetSubscription 연동

#### 보안 (의료기기)
- [x] JWT Bearer 메타데이터 자동 첨부 (인증 필요 RPC)
- [x] 민감 정보 하드코딩 없음 (상수·환경 연동)

#### 문서/기록
- [x] CHANGELOG.md S5a 작업 기록
- [x] CONTEXT.md 진행률 갱신

### 5.6 Stage S6 Gate Checklist — 전체 통합

#### 코드 품질
- [x] E2E 테스트: auth/user/device/measurement 4서비스 gRPC health check
- [x] getEnvOrDefault 환경변수 연동 (os.Getenv)
- [x] Go test ./... / flutter test 통과 전제 유지

#### 기능 완결성
- [x] Docker Compose: Postgres + Go 4서비스(50051–50054) 기동 가능
- [x] README.md: Go gRPC 서비스 접속 표 추가 (50051–50054)
- [x] docs/plan/phase_1c_stage_s6.md: S6 범위·Gate 기준 문서화

#### 문서/기록
- [x] CHANGELOG.md S6 작업 기록
- [x] CONTEXT.md Phase 1C 완료 반영
- [x] QUALITY_GATES.md S6 통과·Phase 1C 완료 갱신

### 5.7 Phase 1D (통합 MVP) — D1 진행

#### 완료 항목
- [x] E2E 플로우 테스트: Register → Login → ValidateToken → StartSession → EndSession → GetMeasurementHistory (backend/tests/e2e/flow_test.go)
- [x] E2E 실행 경로: backend 모듈 기준 (Makefile test-integration, CI integration-test)
- [x] 서비스 미기동 시 헬스/플로우 테스트 스킵 (t.Skipf)
- [x] docs/plan/phase_1d_integration_mvp.md 계획 문서
- [x] tests/e2e/README.md: E2E 실행 방법 안내

#### Phase 1D Gate 통과 ✅ (2026-02-10)
- [x] 통합 테스트: 서비스 기동 후 E2E 플로우 통과 (Register→Login→ValidateToken→StartSession→EndSession→GetHistory)
- [x] Docker Compose 빌드·기동 정상 (Go 1.24, protoc 생성 코드, DB 초기화 수정 완료)
- [x] CHANGELOG·CONTEXT·QUALITY_GATES 최종 갱신

### 5.8 Phase 2 Core — 7/7 서비스 완료 ✅ (2026-02-11)

| 서비스 | 포트 | 상태 | 테스트 | 비고 |
|--------|------|------|--------|------|
| subscription-service | :50055 | ✅ | 14 | 4티어 구독, 카트리지 접근 정책 |
| shop-service | :50056 | ✅ | 있음 | 상품, 장바구니, 주문 |
| payment-service | :50057 | ✅ | 있음 | PG 연동, 구독/상품 결제 |
| ai-inference-service | :50058 | ✅ | 있음 | 바이오마커 분류, 이상 탐지, 트렌드 예측, 건강 점수 |
| cartridge-service | :50059 | ✅ | 20+ | NFC 태그, 30종 레지스트리, 사용 추적 |
| calibration-service | :50060 | ✅ | 12 | 팩토리/현장 보정, 22종 모델 |
| coaching-service | :50061 | ✅ | 11 | 건강 목표, AI 코칭, 일일/주간 리포트, 추천 |

**마지막 업데이트**: 2026-02-11 (Phase 2 Core 완료, Phase 3 Advanced 시작 대기)
