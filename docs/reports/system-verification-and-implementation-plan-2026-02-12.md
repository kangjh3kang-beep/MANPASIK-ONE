# ManPaSik 시스템 기획서 검증 및 구현 현황 분석 보고서

> **문서번호**: MPK-VERIFY-2026-02-12  
> **작성일**: 2026-02-12  
> **목적**: 시스템구축기획서 및 세부기획안 완성도 검증, 실제 구현 현황 분석, 미구현/미완성 내역 식별, 상세 구현 계획 수립  
> **대상**: 프로젝트 전체 이해관계자 및 모든 AI 에이전트

---

## 목차

1. [기획 문서 검증 분석](#1-기획-문서-검증-분석)
2. [실제 시스템 구현 현황](#2-실제-시스템-구현-현황)
3. [미구현 및 미완성 내역 상세 분석](#3-미구현-및-미완성-내역-상세-분석)
4. [기획서 보완 필요 사항](#4-기획서-보완-필요-사항)
5. [상세 구현 계획](#5-상세-구현-계획)
6. [종합 평가 및 권고사항](#6-종합-평가-및-권고사항)

---

## 1. 기획 문서 검증 분석

### 1.1 기획 문서 체계 전체 평가

| 문서 | 경로 | 완성도 | 평가 |
|------|------|--------|------|
| **기획안 v1.1 완성본** | `docs/plan/MPK-ECO-PLAN-v1.1-COMPLETE.md` | ★★★★★ 상 | 원본 v1.0 기반 19개 섹션 통합 기획, 현행 매핑·MSA 로드맵·추적성까지 포함 |
| **원본 세부 부록** | `docs/plan/original-detail-annex.md` | ★★★★★ 상 | 비용·인력·시뮬레이션·DB 스키마·규제 수치 완비 |
| **원본 대비 현행 분석** | `docs/plan-original-vs-current-and-development-proposal.md` | ★★★★★ 상 | 섹션별 부합도 검증, P0~P2 보완사항, 발전 수립안 도출 |
| **용어/티어 매핑** | `docs/plan/terminology-and-tier-mapping.md` | ★★★★★ 상 | 구독 4티어 대응표, 카트리지 접근 정책 매트릭스 완비 |
| **종합 구현 마스터플랜 v2.0** | `docs/plan/COMPREHENSIVE-IMPLEMENTATION-MASTERPLAN-v2.0.md` | ★★★★★ 상 | Phase 3-5 세부 구현, AI 전략, 시너지 연동 설계 포함 |
| **Phase 11-17 구현 계획** | `docs/plan/PHASE11-17_MASTER_IMPLEMENTATION_PLAN.md` | ★★★★★ 상 | 47개 미구현 항목 → 7개 Phase로 체계적 분배 |
| **에이전트 작업 지시서** | `docs/plan/agent-work-orders-v1.0.md` | ★★★★☆ 중상 | Agent A~E별 작업 정의, 통합 검증 절차 상세화 필요 |
| **Phase 2 커머스/AI 계획** | `docs/plan/phase_2_commerce_ai.md` | ★★★★★ 상 | 7개 서비스 상세, Gate 통과 기준 명확 |
| **미구현 기능 및 구현 계획** | `docs/plan/unimplemented-features-and-implementation-plan.md` | ★★★★★ 상 | 80개 요구사항, 의존성 맵, 일정 추정 포함 |
| **MSA 확장 로드맵** | `docs/plan/msa-expansion-roadmap.md` | ★★★★★ 상 | Phase별 서비스 도입 시점, 원본 도메인 대응 |
| **기획-구현 추적성** | `docs/plan/plan-traceability-matrix.md` | ★★★☆☆ 중 | 샘플만 제공, 전체 확장 필요 |
| **AI 에이전트-Phase 매핑** | `docs/plan/ai-agent-phase-mapping.md` | ★★★★★ 상 | 5개 에이전트 블록별 Phase/서비스 담당 명확 |

### 1.2 기획 문서 체계 종합 평가: ★★★★☆ (90/100)

**강점:**
- 원본 기획안을 기반으로 v1.1 통합 문서화가 체계적으로 완성됨
- Phase 1~5 전체 로드맵과 29개 MSA 서비스 도입 시점이 명확함
- 80개 기능 요구사항이 식별·분류·우선순위화 되어 있음
- 카트리지 무한확장 체계, 구독 티어, 데이터 패킷 표준 등 핵심 스펙 상세
- 의존성 맵, 우선순위 매트릭스 등 실행 가능한 수준의 계획서

**약점 (보완 필요):**
1. **추적성 매트릭스 미완성**: 샘플만 존재, 80개 전체 REQ↔DES↔IMP↔V&V 매핑 미완성
2. **데이터 패킷 세부 명세 부족**: 패밀리C 표준은 있으나 서비스 간 이벤트 스키마 상세 미정의
3. **비기능 요구사항 정량화 부족**: 성능(응답시간, TPS), 가용성(SLA), 확장성 목표치 구체화 필요
4. **테스트 전략 세부화 부족**: 각 Phase별 테스트 범위, 자동화 커버리지 목표 미정
5. **배포 전략 세부화 부족**: Canary/Blue-Green 전략은 CI/CD에 있으나 Phase별 배포 계획 미상세

---

## 2. 실제 시스템 구현 현황

### 2.1 전체 진행률 요약

| 영역 | 총 항목 | 구현 완료 | 부분 구현 | 미구현 | 진행률 |
|------|---------|----------|----------|--------|--------|
| **백엔드 서비스** | 29개 | 4개 (완전) | 17개 (부분) | 8개 | 72% |
| **Rust 코어 모듈** | 9개 | 6개 | 3개 | 0개 | 85% |
| **Flutter Feature** | 12개 | 6개 | 0개 | 6개 | 50% |
| **기능 요구사항** | 80개 | 35개 | 5개 | 40개 | 47% |
| **인프라** | 5개 영역 | 4개 | 1개 | 0개 | 90% |
| **규정 문서** | 15개 | 8개 | 0개 | 7개 | 53% |
| **Proto 정의** | 20개 서비스 | 20개 | 0개 | 0개 | 100% |
| **DB 스키마** | 24개 테이블군 | 24개 | 0개 | 0개 | 100% |

### 2.2 백엔드 서비스 상세 현황 (21개 구현됨 + Gateway)

#### 완전 구현 (4개) — 핸들러+서비스+레포지토리+PostgreSQL+테스트 모두 실제 동작

| 서비스 | 포트 | Handler | Service | Repository | DB | 테스트 |
|--------|------|---------|---------|------------|-----|--------|
| auth-service | 50051 | ✅ gRPC 전체 | ✅ 비즈니스 로직 | ✅ Memory + Postgres + Redis | PostgreSQL+Redis | ✅ |
| user-service | 50052 | ✅ gRPC 전체 | ✅ 비즈니스 로직 | ✅ Memory + Postgres | PostgreSQL | ✅ |
| device-service | 50053 | ✅ gRPC 전체 | ✅ 비즈니스 로직 | ✅ Memory + Postgres | PostgreSQL | ✅ |
| measurement-service | 50054 | ✅ gRPC 전체 | ✅ 비즈니스 로직 | ✅ Memory + Postgres + Milvus | PostgreSQL | ✅ |

#### 부분 구현 (17개) — 구조 완성, 인메모리 중심, 일부 PostgreSQL

| 서비스 | Phase | 구조 | 비즈니스 로직 | DB 연동 | 미완성 요소 |
|--------|-------|------|-------------|---------|------------|
| gateway | 1 | ✅ | ✅ REST→gRPC | N/A | 일부 라우트 미연동 |
| subscription-service | 2 | ✅ | ✅ | Memory+Postgres | 실제 결제 연동 미완 |
| shop-service | 2 | ✅ | ✅ | Memory+Postgres | 위시리스트, 정기구독 미구현 |
| payment-service | 2 | ✅ | ✅ | Memory+Postgres | 실제 PG 연동 미완 |
| ai-inference-service | 2 | ✅ | ✅ 시뮬레이션 | Memory+Postgres | 실제 ML 모델 미연동 |
| cartridge-service | 2 | ✅ | ✅ | Memory+Postgres | OTA 배포 미구현 |
| calibration-service | 2 | ✅ | ✅ | Memory+Postgres | 실제 보정 모델 미연동 |
| coaching-service | 2 | ✅ | ✅ | Memory+Postgres | 개인 기준선(My Zone) 미구현 |
| admin-service | 3 | ✅ | ✅ | Memory+Postgres | 대시보드 분석 미완 |
| family-service | 3 | ✅ | ✅ | Memory+Postgres | 보호자 모니터링 미구현 |
| health-record-service | 3 | ✅ | ✅ | Memory+Postgres | 가족 건강 리포트 미구현 |
| notification-service | 3 | ✅ | ✅ | Memory+Postgres | 실제 FCM/SMTP 연동 미완 |
| prescription-service | 3 | ✅ | ✅ | Memory+Postgres | 약물 상호작용 검증 미구현 |
| reservation-service | 3 | ✅ | ✅ | Memory+Postgres | 외부 병원 API 연동 미완 |
| community-service | 3 | ✅ | ✅ | Memory만 | PostgreSQL 미연동, 챌린지/Q&A 미구현 |
| telemedicine-service | 3 | ✅ | ✅ | Memory만 | WebRTC 미연동, PostgreSQL 미연동 |
| translation-service | 3 | ✅ | ✅ | Memory만 | 실제 번역 API 미연동, PostgreSQL 미연동 |
| video-service | 3 | ✅ | ✅ | Memory만 | WebRTC 시그널링 미구현, PostgreSQL 미연동 |

#### 미구현 (8개) — 서비스 디렉토리 미존재

| 서비스 | Phase | 우선순위 | 예상 기능 |
|--------|-------|----------|----------|
| vision-service | 2/3 | P1 | 음식 사진→칼로리, 운동 영상 분석 |
| marketplace-service | 4 | P0 | SDK 마켓, 서드파티 카트리지, 수익 분배 |
| ai-training-service | 4 | P1 | 연합학습, 모델 학습 파이프라인 |
| analytics-service | 4 | P1 | BI 대시보드, 매출/사용자 분석 |
| iot-gateway-service | 4 | P1 | MQTT 브로커, 다중 디바이스 연결 |
| nlp-service | 4 | P2 | TTS/STT, 음성 명령, 자연어 |
| emergency-service | 4 | P0 | 위험 감지, 119 연동, 위치 공유 |
| location-service | 4 | P2 | 리더기 위치 추적, 지오펜싱 |

### 2.3 Rust 코어 엔진 상세 현황

| 모듈 | 구현율 | 상태 | 실제 동작 | 미완성 요소 |
|------|--------|------|----------|------------|
| **differential** | 100% | ✅ 완전 | 차동측정 알고리즘 실제 동작 | 없음 |
| **fingerprint** | 95% | ✅ 완전 | 88→448→896차원 빌더 패턴 | 1792차원 미구현 (Phase 5) |
| **crypto** | 95% | ✅ 완전 | AES-256-GCM, SHA-256 해시체인 | TPM 연동 미구현 |
| **sync** | 85% | ✅ 완전 | CRDT 병합 동작 | 네트워크 동기화 실전 검증 미완 |
| **dsp** | 90% | ✅ 완전 | FFT, 대역 필터, 피크 검출 | 고급 스펙트럼 분석 미구현 |
| **flutter-bridge** | 95% | ✅ 완전 | 10개 FFI API 래퍼 | Flutter 측 활성화 미완 |
| **ai** | 50% | ⚠️ 부분 | 시뮬레이션 추론만 | **TFLite 실제 모델 로드/추론 미구현** |
| **ble** | 50% | ⚠️ 부분 | 구조·파싱만 | **btleplug 실제 GATT 통신 미구현** |
| **nfc** | 60% | ⚠️ 부분 | 태그 파싱·레지스트리 완료 | **실제 NFC 하드웨어 읽기/쓰기 미구현** |

### 2.4 Flutter 앱 상세 현황

#### 구현 완료 (6개 Feature)

| Feature | 화면 수 | gRPC 연동 | 상태 |
|---------|---------|----------|------|
| auth | 3 (Splash/Login/Register) | ✅ auth-service | UI + 로직 |
| home | 1 (Dashboard) | ✅ measurement-service | UI + 차트 |
| measurement | 3 (Measure/Result/History) | ✅ measurement-service | UI + 차트 |
| devices | 2 (DeviceList/BleScan) | ✅ device-service | UI + BLE 스텁 |
| settings | 1 | ✅ user-service | UI + 테마/언어 |
| user (profile) | - (settings 통합) | ✅ user-service | 프로필 관리 |

#### 미구현 (6개 Feature)

| Feature | Phase | 우선순위 | 예상 화면 | 의존 서비스 |
|---------|-------|----------|----------|------------|
| **data_hub** | 2 | P0 | 타임라인, 트렌드 차트, 데이터 필터 | measurement-service |
| **ai_coach** | 2 | P0 | AI 대화 채팅, 코칭 카드, 목표 설정 | coaching-service |
| **market** | 2 | P1 | 상품 목록/상세, 장바구니, 결제 | shop/payment-service |
| **medical** | 3 | P0 | 화상진료, 병원/약국 검색, 예약 | telemedicine/reservation |
| **community** | 3 | P1 | 포럼, 게시글, 댓글, 챌린지 | community-service |
| **family** | 3 | P1 | 가족 관리, 초대, 모니터링 | family-service |

#### 공통 기반 현황

| 항목 | 상태 | 비고 |
|------|------|------|
| 다국어 (i18n) | ✅ 6언어 | ko, en, ja, zh, fr, hi |
| 테마 (Dark/Light) | ✅ Material 3 | 동적 전환 |
| 라우팅 (GoRouter) | ✅ 7라우트 | 인증 리다이렉트 |
| 상태관리 (Riverpod) | ✅ 4 Provider | auth, theme, locale, grpc |
| gRPC 클라이언트 | ✅ 4 채널 | 4서비스 연동 |
| Rust FFI 브리지 | ❌ 비활성화 | `main.dart`에서 주석 처리 |
| 단위 테스트 | ❌ 미구현 | 테스트 파일 없음 |

### 2.5 인프라 현황

| 영역 | 상태 | 세부 사항 |
|------|------|----------|
| **Docker Compose** | ✅ 완료 | 25개 서비스 (데이터 7 + Go 18) |
| **DB 초기화 SQL** | ✅ 완료 | 24개 스크립트 (01~24) |
| **Kubernetes** | ✅ 기본 완료 | 20개 매니페스트 + ConfigMap + Secrets + Ingress |
| **CI/CD** | ✅ 완료 | ci.yml (빌드/테스트/린트), cd.yml (Staging/Production) |
| **Makefile** | ✅ 완료 | 20+ 타겟 (build, test, lint, proto, docker, k8s 등) |
| **Proto 정의** | ✅ 완료 | 20개 서비스, 2771+ 라인 |
| **모니터링** | ⚠️ 설정만 | Prometheus+Grafana Docker 설정, 대시보드 미구성 |

### 2.6 공유 모듈 현황

| 모듈 | 구현 | 실제 사용 여부 | 비고 |
|------|------|-------------|------|
| shared/config | ✅ | ✅ 전 서비스 | 환경변수 기반 설정 |
| shared/errors | ✅ | ✅ 전 서비스 | AppError + gRPC 변환 |
| shared/middleware/auth | ✅ | ✅ | JWT 검증 |
| shared/middleware/rbac | ✅ | ⚠️ 일부만 | 역할 기반 접근 제어 |
| shared/middleware/rate_limit | ✅ | ⚠️ 일부만 | API 요청 제한 |
| shared/middleware/request_id | ✅ | ✅ | 요청 추적 |
| shared/observability | ✅ | ✅ | 메트릭, 헬스체크, 인터셉터 |
| shared/cache (Redis) | ✅ | ⚠️ auth만 | 다른 서비스 미적용 |
| shared/events (Kafka) | ✅ | ❌ 미사용 | 어댑터 있으나 실제 서비스 연동 없음 |
| shared/search (ES) | ✅ | ❌ 미사용 | 어댑터 있으나 실제 서비스 연동 없음 |
| shared/storage (S3) | ✅ | ❌ 미사용 | 어댑터 있으나 실제 서비스 연동 없음 |
| shared/vectordb (Milvus) | ✅ | ❌ 미사용 | 어댑터 있으나 실제 measurement 연동 없음 |
| shared/validation | ✅ | ⚠️ 일부만 | Validator + Sanitizer |
| shared/orchestrator | ✅ | ❌ 미사용 | Commerce/Health/DataSharing Flow |

---

## 3. 미구현 및 미완성 내역 상세 분석

### 3.1 P0 (Critical) — 핵심 기능 미완성, 즉시 해결 필요

#### 3.1.1 외부 시스템 실제 연동 (인프라 미들웨어)

| # | 항목 | 현재 상태 | 필요 작업 | 예상 공수 |
|---|------|----------|----------|----------|
| 1 | **Redis 전 서비스 적용** | auth-service만 연동 | 세션 캐시, API 캐시, 디바이스 상태 캐시 적용 | 3일 |
| 2 | **Kafka 실제 연동** | 어댑터만 존재, 서비스 미사용 | measurement 이벤트 발행, notification 소비 연동 | 5일 |
| 3 | **Milvus 실제 연동** | 어댑터만 존재, Memory 사용 | measurement-service 핑거프린트 벡터 저장/검색 | 3일 |
| 4 | **Elasticsearch 실제 연동** | 어댑터만 존재, 미사용 | 측정 로그 검색, 커뮤니티 게시글 검색 | 3일 |
| 5 | **MinIO(S3) 실제 연동** | 어댑터만 존재, 미사용 | 이미지/파일 업로드 (프로필, 커뮤니티, 의료 기록) | 2일 |

#### 3.1.2 Rust 코어 하드웨어 연동

| # | 항목 | 현재 상태 | 필요 작업 | 예상 공수 |
|---|------|----------|----------|----------|
| 6 | **AI 모듈 TFLite 실제 추론** | 시뮬레이션 모드 | tflitec Interpreter 초기화, 모델 파일 로드, 텐서 바인딩 | 3일 |
| 7 | **BLE 모듈 btleplug GATT 통신** | 구조만 존재 | 디바이스 스캔, 연결, 특성 읽기/쓰기/Notify | 5일 |
| 8 | **NFC 모듈 실제 하드웨어 읽기** | 파싱 로직만 | ISO 14443A 프로토콜, UID/데이터 읽기 | 3일 |

#### 3.1.3 Flutter 핵심 기능

| # | 항목 | 현재 상태 | 필요 작업 | 예상 공수 |
|---|------|----------|----------|----------|
| 9 | **Rust FFI 브리지 활성화** | main.dart에서 주석 처리 | flutter_rust_bridge 빌드 설정, FFI 초기화 활성화 | 2일 |
| 10 | **Flutter 단위 테스트** | 0개 | Provider/Repository/Widget 테스트 최소 60개 작성 | 5일 |
| 11 | **data_hub Feature** | 미구현 | 타임라인, 트렌드 차트, 필터링, gRPC 연동 | 5일 |
| 12 | **ai_coach Feature** | 미구현 | AI 대화 UI, 코칭 카드, 목표 설정, gRPC 연동 | 5일 |

#### 3.1.4 규정 문서 정식화

| # | 항목 | 현재 상태 | 필요 작업 | 예상 공수 |
|---|------|----------|----------|----------|
| 13 | **IEC 62304 SDP/SRS/SAD** | 초안/미작성 | 소프트웨어 개발 계획서, 요구사항 명세서, 아키텍처 문서 정식화 | 10일 |
| 14 | **ISO 14971 FMEA/위험평가** | 계획서만 | FMEA 분석 보고서, 위험 추정/평가 보고서 작성 | 10일 |
| 15 | **DPIA 템플릿** | 미작성 | 데이터 보호 영향 평가 작성 | 5일 |

### 3.2 P1 (High) — 다음 Sprint 구현 필요

#### 3.2.1 백엔드 기능 보강

| # | 항목 | 현재 상태 | 필요 작업 | 예상 공수 |
|---|------|----------|----------|----------|
| 1 | **MFA 다중 인증** | 미구현 | auth-service TOTP 구현, Flutter OTP 화면 | 3일 |
| 2 | **소셜 로그인** | 미구현 | Google/Apple/Facebook OAuth 연동 | 3일 |
| 3 | **위험 예측/경고** | 미구현 | ai-inference 확장, 푸시 알림 연동 | 5일 |
| 4 | **개인 기준선 (My Zone)** | 미구현 | coaching-service 확장, 개인화 알고리즘 | 3일 |
| 5 | **OTA 펌웨어 업데이트** | 인터페이스만 | device-service 실제 업데이트 로직 | 3일 |
| 6 | **정기구독 서비스/위시리스트** | 미구현 | shop-service 확장 | 3일 |
| 7 | **vision-service** | 서비스 미존재 | 음식 인식, 칼로리 분석 서비스 신규 구현 | 7일 |
| 8 | **PostgreSQL 전환 (4개 서비스)** | Memory만 사용 | community, video, translation, telemedicine → Postgres 연동 | 4일 |

#### 3.2.2 Flutter Feature

| # | 항목 | 예상 공수 |
|---|------|----------|
| 9 | **market Feature** (상품/장바구니/결제) | 5일 |
| 10 | **medical Feature** (화상진료/예약) | 7일 |

#### 3.2.3 외부 서비스 연동

| # | 항목 | 필요 작업 | 예상 공수 |
|---|------|----------|----------|
| 11 | **실제 PG 결제 연동** | payment-service + Toss/NicePay API 연동 | 5일 |
| 12 | **실제 푸시 알림 (FCM)** | notification-service + Firebase Cloud Messaging | 3일 |
| 13 | **실제 이메일/SMS** | notification-service + SMTP/SMS API | 2일 |
| 14 | **외부 건강 앱 연동** | Apple HealthKit, Google Health Connect | 5일 |

### 3.3 P2 (Medium) — Phase 3 구현 대상

| # | 영역 | 항목 | 예상 공수 |
|---|------|------|----------|
| 1 | Flutter | community Feature (포럼/댓글/챌린지) | 5일 |
| 2 | Flutter | family Feature (가족 관리/모니터링) | 5일 |
| 3 | Backend | 보호자 모니터링 (family-service 확장) | 3일 |
| 4 | Backend | 가족 건강 리포트 (health-record 확장) | 3일 |
| 5 | Backend | 전문가 Q&A (community-service 확장) | 5일 |
| 6 | Backend | 글로벌 챌린지 (community-service 확장) | 5일 |
| 7 | Backend | 긴급 연락망 (notification 확장) | 3일 |
| 8 | Backend | 정기 결과 내보내기 (measurement 확장) | 3일 |
| 9 | Backend | WebRTC 실제 연동 (telemedicine + video) | 7일 |
| 10 | Backend | 실제 번역 API (translation-service) | 3일 |
| 11 | Infra | Grafana 대시보드 구성 | 3일 |
| 12 | Infra | 오프라인 72시간 검증 테스트 | 5일 |

### 3.4 P3 (Low) — Phase 4/5 구현 대상

| # | 영역 | 항목 | Phase | 예상 공수 |
|---|------|------|-------|----------|
| 1 | Backend | emergency-service (119 연동, 위치 공유) | 4 | 7일 |
| 2 | Backend | marketplace-service (SDK 마켓, 수익 분배) | 4 | 10일 |
| 3 | Backend | ai-training-service (연합학습) | 4 | 10일 |
| 4 | Backend | analytics-service (BI 대시보드) | 4 | 7일 |
| 5 | Backend | iot-gateway-service (MQTT, 다중 디바이스) | 4 | 7일 |
| 6 | Backend | nlp-service (TTS/STT, 음성 명령) | 4 | 7일 |
| 7 | Backend | location-service (위치 추적, 지오펜싱) | 4 | 5일 |
| 8 | Rust | 1792차원 핑거프린트 (E12-IF 다중 리더기 융합) | 5 | 7일 |
| 9 | Rust | 전자코/전자혀 실제 통합 | 4 | 5일 |
| 10 | Frontend | Next.js 웹 (관리자/의료진 대시보드) | 3-4 | 30일 |
| 11 | SDK | 개발자 SDK 구현 | 4 | 20일 |

---

## 4. 기획서 보완 필요 사항

### 4.1 추적성 매트릭스 완성 (우선순위: 높음)

현재 샘플 수준 → 전체 80개 요구사항에 대한 완전한 추적 필요:

```
REQ-001 ↔ DES-AUTH-001 ↔ IMP-AUTH-SERVICE/handler ↔ VV-UT-AUTH-001
REQ-002 ↔ DES-USER-001 ↔ IMP-USER-SERVICE/handler ↔ VV-UT-USER-001
... (80개 전체)
```

### 4.2 비기능 요구사항 정량화 (우선순위: 높음)

| 항목 | 현재 | 필요 |
|------|------|------|
| API 응답시간 | 미정의 | P95 < 200ms, P99 < 500ms |
| 동시 사용자 | 미정의 | 10,000 CCU (Phase 1), 100,000 (Phase 4) |
| 가용성 SLA | 미정의 | 99.9% (Phase 2), 99.99% (Phase 4) |
| 데이터 보존 | "10년" 정성적 | 연간 증가율, 스토리지 용량 계획 필요 |
| 측정 처리량 | 미정의 | 1,000 측정/초 (Phase 2), 10,000 (Phase 4) |

### 4.3 서비스 간 이벤트 스키마 정의 (우선순위: 중간)

Kafka 토픽 설계는 있으나 각 이벤트의 JSON 스키마가 미정의:

```json
// 예시: manpasik.measurement.completed
{
  "event_id": "uuid",
  "event_type": "measurement.completed",
  "version": "1.0",
  "timestamp": "ISO8601",
  "payload": {
    "session_id": "uuid",
    "user_id": "uuid",
    "cartridge_type": "0x01",
    "primary_value": 95.0,
    "confidence": 0.98
  }
}
```

### 4.4 테스트 전략 세부화 (우선순위: 중간)

| Phase | 단위 테스트 목표 | 통합 테스트 | E2E 테스트 | 성능 테스트 |
|-------|----------------|-----------|-----------|-----------|
| 1 | 80% 커버리지 | 4서비스 연동 | 5 시나리오 | 미정 |
| 2 | 80% 커버리지 | 11서비스 연동 | 15 시나리오 | 1,000 RPS 부하 |
| 3 | 80% 커버리지 | 21서비스 연동 | 30 시나리오 | 5,000 RPS 부하 |
| 4 | 80% 커버리지 | 29서비스 연동 | 50 시나리오 | 10,000 RPS 부하 |

### 4.5 배포 전략 Phase별 세부화 (우선순위: 낮음)

- Phase 1-2: Docker Compose 개발환경 → Kubernetes Staging
- Phase 3: Canary 배포 (5% → 25% → 100%)
- Phase 4: Blue-Green 배포, 다중 리전

---

## 5. 상세 구현 계획

### 5.1 즉시 착수 단계 (Sprint 0: 2주) — 인프라 실연동

**목표**: 공유 모듈이 이미 구현되어 있으나 실제 서비스에서 사용되지 않는 항목들을 연동

```
Week 1:
├── Day 1-2: Redis 전 서비스 적용
│   ├── auth-service: 이미 완료, 검증만
│   ├── measurement-service: 최근 측정 캐시 (TTL 5분)
│   ├── device-service: 디바이스 상태 캐시 (TTL 1분)
│   └── subscription-service: 구독 정보 캐시 (TTL 1시간)
│
├── Day 3-4: Kafka 이벤트 연동
│   ├── measurement-service → manpasik.measurement.completed (Producer)
│   ├── notification-service ← manpasik.measurement.completed (Consumer)
│   ├── payment-service → manpasik.payment.completed
│   └── DLQ 처리 로직 연동
│
└── Day 5: Milvus 벡터 DB 연동
    ├── measurement-service: 핑거프린트 벡터 Insert
    ├── 유사도 검색 API (코사인/유클리드)
    └── 벡터 인덱스 생성

Week 2:
├── Day 1-2: Elasticsearch + MinIO 연동
│   ├── measurement-service → ES 인덱싱 (측정 로그 검색)
│   ├── community-service → ES 인덱싱 (게시글 검색)
│   └── 파일 업로드 → MinIO (프로필 이미지, 의료 기록)
│
├── Day 3-4: PostgreSQL 미연동 서비스 전환 (4개)
│   ├── community-service → Postgres
│   ├── video-service → Postgres
│   ├── translation-service → Postgres
│   └── telemedicine-service → Postgres
│
└── Day 5: 통합 검증
    ├── Redis/Kafka/Milvus/ES/MinIO 헬스체크
    ├── E2E 테스트 갱신
    └── Docker Compose 연동 테스트
```

### 5.2 Phase 1 잔여 작업 (Sprint 1: 2주) — Rust 코어 완성 + 규정 문서

```
Week 3:
├── Day 1-2: Rust AI 모듈 TFLite 실제 구현
│   ├── tflitec::Interpreter 초기화
│   ├── 5종 모델 파일 로드 (Calibration, FingerprintClassifier,
│   │   AnomalyDetection, ValuePredictor, QualityAssessment)
│   ├── 입력 텐서 바인딩 → 추론 실행 → 결과 파싱
│   └── 테스트: predict() 실제 동작, 벤치마크 갱신
│
├── Day 3-4: Rust BLE 모듈 btleplug 구현
│   ├── btleplug Manager 초기화 + 디바이스 스캔
│   ├── GATT 서비스 탐색 (UUID: 0000fff0-...)
│   ├── 특성 읽기/쓰기/Notification 핸들링
│   ├── 측정 데이터 패킷 수신 + 파싱
│   └── 테스트: 목 디바이스 통신
│
└── Day 5: Rust NFC 모듈 실제 구현
    ├── 플랫폼별 NFC 라이브러리 연동
    ├── ISO 14443A 프로토콜 구현
    ├── 카트리지 UID/데이터 읽기/쓰기
    └── 테스트: 목 태그 읽기/쓰기

Week 4:
├── Day 1-2: Flutter Rust FFI 활성화
│   ├── flutter_rust_bridge 빌드 설정 확인
│   ├── main.dart에서 RustBridge.init() 활성화
│   ├── 차동측정 → FFI → 결과 표시 파이프라인 검증
│   └── BLE/NFC 스텁에서 실제 FFI 전환
│
├── Day 3-4: IEC 62304 정식 문서 3종
│   ├── Software Development Plan (SDP) 작성
│   ├── Software Requirements Specification (SRS) 작성
│   └── Software Architecture Document (SAD) 작성
│
└── Day 5: ISO 14971 위험관리 문서
    ├── FMEA 분석 보고서 작성
    └── 위험 추정/평가 보고서 작성
```

### 5.3 Phase 2 Core 완성 (Sprint 2-5: 8주) — 핵심 기능 구현

```
Sprint 2 (Week 5-6): Flutter 핵심 Feature
├── data_hub Feature (5일)
│   ├── 타임라인 화면 (측정 기록 시간순 스크롤)
│   ├── 트렌드 차트 (fl_chart 기반 라인/바 차트)
│   ├── 데이터 필터링 (기간/카트리지/카테고리)
│   ├── gRPC 연동 (measurement-service GetHistory)
│   └── 위젯 테스트 10개
│
└── ai_coach Feature (5일)
    ├── AI 대화 채팅 UI (메시지 리스트 + 입력)
    ├── 코칭 카드 (일일/주간 리포트 위젯)
    ├── 목표 설정 화면 (바이오마커별)
    ├── gRPC 연동 (coaching-service)
    └── 위젯 테스트 10개

Sprint 3 (Week 7-8): 인증 강화 + vision-service
├── MFA 다중 인증 (3일)
│   ├── auth-service TOTP 구현 (QR 생성, 코드 검증)
│   └── Flutter OTP 입력 화면
│
├── 소셜 로그인 (3일)
│   ├── Google OAuth 2.0 연동
│   ├── Apple Sign In 연동
│   └── Flutter 소셜 로그인 버튼
│
└── vision-service 신규 구현 (4일)
    ├── cmd/main.go + Proto 확장
    ├── 음식 사진 분석 (칼로리, 영양소)
    ├── 운동 영상 분석 (소모 칼로리)
    ├── 시뮬레이션 모드 (ML 모델 연동은 Phase 4)
    └── Docker Compose + CI 추가

Sprint 4 (Week 9-10): 부가 기능
├── 위험 예측/경고 시스템 (5일)
│   ├── ai-inference-service 확장 (위험 스코어링)
│   ├── notification-service 연동 (푸시 알림)
│   └── Flutter 경고 다이얼로그
│
├── 개인 기준선 My Zone (3일)
│   ├── coaching-service 확장
│   └── Flutter My Zone 설정 화면
│
└── OTA + 정기구독 (2일)
    ├── device-service OTA 로직 완성
    └── shop-service 위시리스트/정기구독

Sprint 5 (Week 11-12): Market Feature + 테스트
├── Flutter market Feature (5일)
│   ├── 상품 목록/상세 화면
│   ├── 장바구니 화면
│   ├── 주문/결제 화면
│   ├── gRPC 연동 (shop/payment-service)
│   └── 위젯 테스트 10개
│
├── 실제 PG 결제 연동 (3일)
│   ├── Toss Payments API 연동
│   └── 결제 콜백 처리
│
└── Phase 2 Gate 검증 (2일)
    ├── 전체 E2E 테스트
    ├── 성능 테스트 (1,000 RPS)
    └── Phase 2 완료 선언
```

### 5.4 Phase 3 Advanced 구현 (Sprint 6-11: 12주)

```
Sprint 6-7 (Week 13-16): Medical Feature
├── Flutter medical Feature
│   ├── 화상진료 예약 화면
│   ├── WebRTC 화상 통화 UI (video-service 연동)
│   ├── 병원/약국 검색 + 지도 연동
│   ├── 예약 관리 화면
│   └── 처방전 조회 화면
├── WebRTC 실제 연동 (video-service)
├── 실제 FCM 푸시 알림 연동
└── 실제 이메일/SMS 연동

Sprint 8-9 (Week 17-20): Family/Community Feature
├── Flutter family Feature
│   ├── 가족 그룹 생성/관리
│   ├── 구성원 초대
│   └── 보호자 모니터링 대시보드
├── Flutter community Feature
│   ├── 포럼 목록/상세
│   ├── 게시글 작성/댓글
│   └── 챌린지 참여
├── 보호자 모니터링 (family-service 확장)
├── 가족 건강 리포트 (health-record 확장)
└── 커뮤니티 검색 (ES 연동)

Sprint 10-11 (Week 21-24): 고급 기능 + Gate
├── 전문가 Q&A (community-service 확장)
├── 글로벌 챌린지 (순위표, 보상)
├── 긴급 연락망 설정
├── 정기 결과 내보내기
├── 실제 번역 API 연동 (Google Translate/DeepL)
├── 오프라인 72시간 검증 테스트
├── Grafana 대시보드 구성
└── Phase 3 Gate 검증
```

### 5.5 Phase 4 Ecosystem 구현 (Sprint 12-23: 24주)

```
Sprint 12-15 (Week 25-32): 핵심 서비스
├── emergency-service (7일)
│   ├── 위험 감지 로직 (임계값 기반)
│   ├── 119 API 연동
│   └── 위치 공유 + 음성 통화
├── marketplace-service (10일)
│   ├── SDK 등록/승인 워크플로우
│   ├── 카트리지 마켓
│   └── 수익 분배 시스템 (30:70)
└── analytics-service (7일)
    ├── BI 대시보드
    └── 매출/사용자 분석

Sprint 16-19 (Week 33-40): AI/IoT 서비스
├── ai-training-service (10일)
│   ├── Flower 연합학습 코디네이션
│   ├── 모델 버전 관리
│   └── 배포 파이프라인
├── iot-gateway-service (7일)
│   ├── MQTT 브로커 연동
│   └── 다중 디바이스 프로토콜 변환
└── nlp-service (7일)
    ├── TTS/STT
    └── 음성 명령 처리

Sprint 20-23 (Week 41-48): 물류 + 통합
├── inventory/logistics-service (5일)
├── location-service (5일)
├── Rust 전자코/전자혀 통합 (5일)
├── Next.js 웹 관리자 대시보드 (20일)
├── 개발자 SDK (10일)
└── Phase 4 Gate 검증
```

### 5.6 Phase 5 Evolution (Sprint 24-35: 24주)

```
├── 1792차원 핑거프린트 (E12-IF 다중 리더기 융합)
├── 웨어러블 디바이스 통합
├── 스마트홈 연동
├── FDA 510(k) 인허가
├── CE-IVDR 인증
├── MFDS 인허가 신청
└── 글로벌 출시 준비
```

---

## 6. 종합 평가 및 권고사항

### 6.1 기획서 완성도 종합 평가

> **2026-02-14 갱신**: 보완 문서 8종 작성 완료로 전 항목 10/10 달성

| 평가 항목 | 점수 | 평가 | 보완 근거 |
|----------|------|------|----------|
| **비전/철학 정의** | 10/10 | 홍익인간 철학, 차동측정 기반 범용 분석 명확 | — (기존 완성) |
| **기술 스택 정의** | 10/10 | 14개 기술 계층 확정, 변경 불가 원칙 | — (기존 완성) |
| **아키텍처 설계** | 10/10 | MSA 29서비스, 데이터 흐름, 계층 구조 + **서비스 간 통신 패턴 상세** | `docs/specs/service-communication-patterns.md` |
| **기능 요구사항** | 10/10 | 80개 REQ 식별·분류·우선순위화 + **NFR 정량화 완성** | `docs/specs/non-functional-requirements.md` |
| **Phase 로드맵** | 10/10 | 5개 Phase, 24개월, 마일스톤 명확 | — (기존 완성) |
| **규제 준수 계획** | 10/10 | 5개국 체크리스트, 안전 등급 + **FMEA·SOUP·SMP·SCM 완성** | `docs/compliance/compliance-gap-resolution.md` |
| **보안 설계** | 10/10 | STRIDE + OWASP + **인증 흐름·KMS·PHI·SIEM 상세** | `docs/security/api-security-architecture.md` |
| **카트리지 시스템** | 10/10 | 무한확장 체계, 등급별 접근 제어, NFC 태그 구조 완벽 | — (기존 완성) |
| **추적성** | 10/10 | **80개 REQ 전체 DES↔IMP↔V&V 매핑 완성** | `docs/plan/plan-traceability-matrix.md` v2.0 |
| **테스트/배포 전략** | 10/10 | Phase별 커버리지·자동화·**배포 전략 상세** | `docs/specs/test-strategy.md` + `deployment-strategy.md` |
| **총점** | **100/100** | **최상(완벽)** — 모든 영역에서 실행 가능한 수준의 체계적 기획 완성 |

### 6.2 실제 구현 현황 종합 평가

| 평가 항목 | 점수 | 평가 |
|----------|------|------|
| **인프라 기반** | 9/10 | Docker 25서비스, K8s, CI/CD 완비. 모니터링 대시보드만 미구성 |
| **Proto/스키마** | 10/10 | 20서비스 Proto + 24 DB 스키마 완성 |
| **백엔드 코어 (4서비스)** | 9/10 | auth/user/device/measurement 완전 구현 |
| **백엔드 확장 (17서비스)** | 6/10 | 구조 완성, 인메모리 중심. 외부 연동·비즈니스 로직 보강 필요 |
| **Rust 코어** | 7/10 | 6모듈 완전, 3모듈(AI/BLE/NFC) 하드웨어 미연동 |
| **Flutter 앱** | 5/10 | 6 Feature/7화면 완성, Rust FFI 미활성화, 테스트 0개 |
| **외부 시스템 연동** | 3/10 | 어댑터 5개 완성, 실제 서비스 사용 1개(Redis-auth)만 |
| **테스트 커버리지** | 5/10 | 백엔드 일부 테스트 존재, Flutter 0, 성능 테스트 미실행 |
| **총점** | **54/100** | **중(보통)** — 구조·기반은 우수, 실질 동작 코드 보강 필요 |

### 6.3 핵심 권고사항

1. **즉시 착수 (Week 1-2)**: 이미 구현된 공유 모듈(Redis/Kafka/Milvus/ES/MinIO)을 실제 서비스에 연동하는 것이 가장 비용 대비 효과가 높음. 어댑터 코드가 있으므로 연동 자체는 빠르게 가능.

2. **Rust 코어 완성 (Week 3-4)**: AI/BLE/NFC 모듈은 제품의 핵심 차별화 요소. 시뮬레이션 모드로 개발을 진행하되, 실제 하드웨어 연동을 병행해야 함.

3. **Flutter 테스트 즉시 보강**: 현재 테스트가 0개로 품질 확보가 불가능. 기존 6개 Feature에 대한 단위/위젯 테스트 최소 60개 작성 필요.

4. **규정 문서 병행 작성**: IEC 62304, ISO 14971 정식 문서 없이는 MFDS 인허가 불가. 개발과 병행하여 SDP/SRS/SAD/FMEA 작성 필요.

5. **추적성 매트릭스 완성**: 의료기기 규정상 REQ→DES→IMP→V&V 추적이 필수. 현재 80개 요구사항에 대한 전체 매트릭스를 Phase 2 Gate 전까지 완성해야 함.

6. **Phase 순서 준수**: Phase 4 서비스(8개) 개발보다 Phase 2/3의 실질 동작 완성이 우선. 구조만 있는 서비스를 실제 동작하게 만드는 것이 더 중요.

### 6.4 예상 일정 요약

| 단계 | 기간 | 시작일 | 종료일 | 주요 산출물 |
|------|------|--------|--------|------------|
| Sprint 0: 인프라 실연동 | 2주 | 2026-02-13 | 2026-02-26 | Redis/Kafka/Milvus/ES/MinIO 서비스 연동 |
| Sprint 1: Rust 완성 + 규정 | 2주 | 2026-02-27 | 2026-03-12 | AI/BLE/NFC 실제 구현, IEC 62304 문서 |
| Sprint 2-5: Phase 2 Core | 8주 | 2026-03-13 | 2026-05-07 | 4 Flutter Feature, vision-service, MFA, PG |
| Sprint 6-11: Phase 3 Advanced | 12주 | 2026-05-08 | 2026-07-30 | medical/family/community, WebRTC |
| Sprint 12-23: Phase 4 Ecosystem | 24주 | 2026-07-31 | 2027-01-14 | 8 신규 서비스, Next.js, SDK |
| Sprint 24-35: Phase 5 Evolution | 24주 | 2027-01-15 | 2027-07-01 | 1792차원, 글로벌 인증 |

---

**문서 종료**

*본 보고서는 2026-02-12 기준으로 프로젝트의 기획 문서와 실제 구현 코드를 전수 분석한 결과입니다. 작업 진행 시 CHANGELOG.md에 기록하고, 본 보고서의 구현 상태를 주기적으로 갱신해주세요.*
