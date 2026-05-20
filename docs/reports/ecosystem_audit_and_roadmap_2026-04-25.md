# 만파식(萬波息) 글로벌 생태계 종합 실사 보고서 및 상세 구축계획

**문서번호**: MPS-AUDIT-ROADMAP-2026-04-25
**작성일**: 2026-04-25
**기준**: 코드베이스 전수 실사 + 기획문서 7종 대조분석
**범위**: 6-Layer 전계층 (HW→Rust→Flutter→Go→AI→Infra)
**선행문서**: MPS-SYS-MASTER-v2.0, MPS-TECH-SPEC-v2.4.3, 100-percent-implementation-plan-v1.0

---

## 제1부: 현재 구축 현황 전수 실사

### 1.1 자산 인벤토리 (2026-04-25 실측)

| 계층 | 항목 | 2026-02-23 기준 | **2026-04-25 실측** | 변화 |
|------|------|:---:|:---:|:---:|
| **Proto** | 서비스 / RPC / 줄수 | 33 / 233 / 4,274 | **33 / 237 / 4,380** | +4 RPC |
| **Go 백엔드** | 서비스 수 | 33 + 1 GW | **35 + 1 GW** | +2 (audit, digital-twin) |
| **Go 테스트** | 테스트 파일 | ~49 | **54** | +5 |
| **Go PostgreSQL** | 영속화 저장소 | **0** | **39 파일, 전 서비스** | 0%→구조완성 |
| **Gateway** | 라우트 파일 / 엔드포인트 | 16 / 240+ | **16 / 240+** | 유지 |
| **DB 스키마** | SQL 파일 / 테이블 | 34 / 112 | **42 / 130+** | +8 SQL |
| **Rust 엔진** | 줄수 / 테스트 | 6,508 / 62 | **6,569 / 90** | +61줄, +28테스트 |
| **Flutter 화면** | Screen 파일 | 97 | **72** | 리팩토링 통합 |
| **Flutter REST** | 메서드 수 | 230 | **230** | 유지 |
| **Flutter 위젯** | 커스텀 위젯 | 49 | **49** | 유지 |
| **Flutter 저장소** | REST Repository | 16 | **16** | 유지 |
| **Flutter 라우트** | 라우트 정의 | 90 | **90** | 유지 |
| **Flutter 테스트** | 테스트 파일 | 30 | **51** | +21 (+70%) |
| **K8s 서비스** | YAML 정의 | — | **36** | 신규 |
| **Docker** | 컴포즈 서비스 | — | **33** | 신규 |
| **AI/ML** | Python 파일 | — | **3** | 골격 |
| **기획 문서** | 마크다운 | 135 | **135+** | 유지 |

### 1.2 Go 서비스 로직 깊이 재분류 (2026-04-25)

이전 분류 대비 **대폭 성장**:

| 분류 | 서비스 (줄수 기준) | 개수 | 변화 |
|------|-----|:---:|:---:|
| **REAL (500+줄)** | auth(570), health-record(782), family(490), notification(654), admin(1280), payment(365), **ai-inference(1285)**, **coaching(1157)**, **data-platform(760)**, **cartridge(701)**, **subscription(669)**, **community(632)**, **prescription(605)**, **measurement(587)**, **calibration(552)** | **15** | 6→15 |
| **CRUD (300~499줄)** | nlp(412), assistant(442), iot-gateway(394), vision(398), telemedicine(379), translation(362), emergency(357), marketplace(333), video(391), device(304), analytics(303), shop(291) | **12** | 15→12 |
| **SKELETON (<300줄)** | concept(280), audit(280), cartridge-store(225), user(208), voice-profile(208), digital-twin(189) | **6** | 12→6 |
| **GATEWAY** | gateway (REST→gRPC 프록시, 별도 구조) | **1** | 유지 |

### 1.3 Go 서비스별 전수 감사

| 서비스 | main.go | gRPC | Memory | **Postgres** | 테스트 | 로직 줄수 | 등급 |
|--------|:---:|:---:|:---:|:---:|:---:|:---:|:---:|
| admin-service | Y | Y | Y | **Y** | 5 | 1,280 | REAL |
| ai-inference-service | Y | Y | Y | **Y** | 2 | 1,285 | REAL |
| analytics-service | Y | Y | Y | **Y** | 1 | 303 | CRUD |
| assistant-service | Y | Y | Y | **Y** | 2 | 442 | CRUD |
| audit-service | Y | Y | Y | **Y** | 1 | 280 | SKELETON |
| auth-service | Y | Y | Y | **Y** | 6 | 570 | REAL |
| calibration-service | Y | Y | Y | **Y** | 1 | 552 | REAL |
| cartridge-service | Y | Y | Y | **Y** | 1 | 701 | REAL |
| cartridge-store-service | Y | Y | Y | **Y** | 1 | 225 | SKELETON |
| coaching-service | Y | Y | Y | **Y** | 1 | 1,157 | REAL |
| community-service | Y | Y | Y | **Y** | 1 | 632 | REAL |
| concept-service | Y | Y | Y | **Y** | 1 | 280 | SKELETON |
| data-platform-service | Y | Y | Y | **Y** | 1 | 760 | REAL |
| device-service | Y | Y | Y | **Y** | 1 | 304 | CRUD |
| digital-twin-service | Y | Y | Y | **Y** | 1 | 189 | SKELETON |
| emergency-service | Y | Y | Y | **Y** | 3 | 357 | CRUD |
| family-service | Y | Y | Y | **Y** | 1 | 490 | REAL |
| gateway | Y | — | — | — | 2 | — | GW |
| health-record-service | Y | Y | Y | **Y** | 1 | 782 | REAL |
| iot-gateway-service | Y | Y | Y | **Y** | 1 | 394 | CRUD |
| marketplace-service | Y | Y | Y | **Y** | 1 | 333 | CRUD |
| measurement-service | Y | Y | Y | **Y** | 1 | 587 | REAL |
| nlp-service | Y | Y | Y | **Y** | 2 | 412 | CRUD |
| notification-service | Y | Y | Y | **Y** | 1 | 654 | REAL |
| payment-service | Y | Y | Y | **Y** | 1 | 365 | CRUD |
| prescription-service | Y | Y | Y | **Y** | 1 | 605 | REAL |
| reservation-service | Y | Y | Y | **Y** | 1 | 481 | REAL |
| shop-service | Y | Y | Y | **Y** | 1 | 291 | CRUD |
| subscription-service | Y | Y | Y | **Y** | 1 | 669 | REAL |
| telemedicine-service | Y | Y | Y | **Y** | 3 | 379 | CRUD |
| translation-service | Y | Y | Y | **Y** | 2 | 362 | CRUD |
| user-service | Y | Y | Y | **Y** | 1 | 208 | SKELETON |
| video-service | Y | Y | Y | **Y** | 1 | 391 | CRUD |
| vision-service | Y | Y | Y | **Y** | 2 | 398 | CRUD |
| voice-profile-service | Y | Y | Y | **Y** | 1 | 208 | SKELETON |

### 1.4 Flutter 프론트엔드 상세 (feature별)

| Feature | 화면수 | 테스트 | REST 바인딩 |
|---------|:---:|:---:|:---:|
| settings | 13 | 2 | O |
| market | 12 | 1 | O |
| admin | 9 | 1 | O |
| medical | 6 | 2 | O |
| family | 6 | 2 | O |
| community | 6 | 2 | O |
| measurement | 4 | 2 | O |
| auth | 4 | 2 | O |
| ai_coach | 3 | 2 | O |
| data_hub | 2 | 2 | O |
| devices | 2 | 2 | O |
| coach | 1 | — | O |
| chat | 1 | 1 | O |
| home | 1 | 1 | O |
| notification | 1 | 2 | O |
| onboarding | 1 | 1 | O |
| **합계** | **72** | **51** | **16 저장소** |

### 1.5 Rust 코어 엔진 상세

| 모듈 | 파일 | 역할 | 테스트 |
|------|------|------|:---:|
| `lib.rs` | SSOT 상수 + 모듈 선언 | 진입점 | — |
| `api.rs` | FFI API 인터페이스 | 외부 노출 | — |
| `differential/` | 차분식 Sdiff = S - αR | 핵심 알고리즘 | O |
| `fingerprint/` | 88→448→896→1792차원 | 특성 추출 | O |
| `ble/` | BLE GATT 프로토콜 | 통신 | O |
| `nfc/` | NFC 매니페스트 파싱 | 카트리지 인식 | O |
| `crypto/` | AES-256-GCM, SHA-256 | 보안 | O |
| `dsp/` | FFT, 신호처리 | 전처리 | O |
| `ai/` | TFLite/ONNX 추론 | AI 엣지 | O |
| `sync/` | CRDT 동기화 | 오프라인 | O |
| `flutter-bridge/` | 16 FFI 함수 | 브릿지 | O |
| **총합** | **12 파일, 6,569줄** | **9 모듈** | **90 테스트** |

### 1.6 인프라 현황

| 항목 | 상태 | 상세 |
|------|------|------|
| **DB 스키마** | 42 SQL, 130+ 테이블 | 설계 완료, 영속화 연결 진행중 |
| **K8s** | 36 서비스 YAML | 배포 미실행 (로컬 docker-compose) |
| **Prometheus** | prometheus.yml + alert_rules.yml | 구성 완료, 실행 미검증 |
| **Docker Compose** | 33 서비스 정의 | 로컬 개발 환경 |
| **CI** | ci.yml (832B) | 기본 린트 수준 |
| **CD** | cd.yml (19KB) + 생성 스크립트 | 파이프라인 설계 완료 |
| **AI/ML** | 3 Python 파일 | 골격 수준 |

---

## 제2부: 기획 vs 현실 갭 분석

### 2.1 완성도 매트릭스 업데이트

| 차원 | 2026-02-23 | **2026-04-25** | 목표 | 잔여 갭 |
|------|:---:|:---:|:---:|:---:|
| 구조적 완성도 | 95% | **98%** | 100% | 2% |
| Go 비즈니스 로직 | 40% | **62%** | 100% | 38% |
| Flutter 데이터 바인딩 | 87% | **90%** | 100% | 10% |
| **DB 영속화** | **0%** | **25%** | 100% | **75%** |
| 외부 연동 | 25% | **30%** | 100% | 70% |
| 테스트 커버리지 | 35% | **42%** | 80% | 38% |
| 네이티브 기능 (Rust FFI) | 10% | **25%** | 80% | 55% |
| 운영/배포 | 15% | **22%** | 90% | 68% |
| 규제 문서 | 80% | **82%** | 100% | 18% |
| **종합** | **~44%** | **~53%** | **100%** | **~47%** |

### 2.2 핵심 변화 (2개월간 +9%p 성장)

| 성과 | 상세 |
|------|------|
| PostgreSQL 저장소 전 서비스 구조 생성 | 39 파일 — Phase A 구조적 기반 완료 |
| Go 서비스 로직 대폭 고도화 | REAL 6→15, SKELETON 12→6 |
| Flutter 테스트 +70% 증가 | 30→51 파일 |
| Rust 테스트 +45% 증가 | 62→90 테스트 |
| Proto +4 RPC 확장 | 237 RPC로 성장 |
| 신규 서비스 2개 추가 | audit-service, digital-twin-service |
| DB 스키마 +8 SQL 추가 | 42 SQL, 130+ 테이블 |
| K8s 36 서비스 YAML 생성 | 배포 인프라 기반 |

### 2.3 기획 문서 대비 갭 분석

| 기획서 항목 | 기획 목표 | 현실 | 갭 수준 |
|------------|----------|------|:---:|
| **마스터플랜 v2.0**: 6-Layer 아키텍처 | 전 계층 구현 | L1-L2(HW) 미구현, L3-L6 진행중 | 중 |
| **Tech Spec v2.4.3**: HAL(harness) + 9블록 AFE | Rust SensorTrait 구현 | harness/ 디렉토리 미생성, 기본 모듈만 존재 | **대** |
| **Tech Spec v2.4.3**: 디지털 트윈/자가치유 | §6.4~§6.16 모듈 | digital_twin/, checkup/, context/, diagnostics/, healing/ 미구현 | **대** |
| **UX 사이트맵 v1.4.2**: 120+ 화면 | 전체 화면 구현 | 72 화면 구현 (60%) | 중 |
| **100% 계획 v1.0**: Phase A~H 32 스프린트 | 순차 실행 | Phase A 25%, Phase B 부분, 나머지 미착수 | 중 |
| **에이전트 작업지시서**: 6 에이전트 병렬 | 도메인별 고도화 | 부분 실행 (Agent A~E 일부) | 중 |

---

## 제3부: 향후 완벽 구현 · 구축 계획

### 3.1 수정된 로드맵 (현실 반영)

기존 8 Phase 32 Sprint에서 이미 달성된 부분을 차감하고, 잔여 작업을 재정렬:

```
Phase A': 영속화 완성        ─── 3 Sprint (기존 4 → 구조완성 차감)
Phase B': 로직 고도화        ─── 4 Sprint (기존 6 → REAL 확대 차감)
Phase C : 외부 연동          ─── 6 Sprint (변경 없음)
Phase D': 네이티브 통합      ─── 3 Sprint (기존 4 → Rust Phase 0 차감)
Phase E : AI/ML 파이프라인   ─── 4 Sprint (변경 없음)
Phase F': 테스트/품질        ─── 3 Sprint (기존 4 → 테스트 증가 차감)
Phase G : 운영/배포          ─── 2 Sprint (변경 없음)
Phase H : 규제/인증          ─── 2 Sprint (변경 없음)
합계: 27 Sprint (기존 32 → 5 Sprint 절감)
```

### 3.2 Phase A': 영속화 완성 (3 Sprint)

PostgreSQL 저장소 구조는 전 서비스 생성 완료. **실제 SQL 쿼리 구현 + DB 연결 전환**이 핵심.

| Sprint | 작업 | 대상 서비스 | 산출물 |
|--------|------|-----------|--------|
| **A'-1** | 공통 DB 래퍼 + REAL 15서비스 쿼리 구현 | auth, admin, health-record, family, notification, payment, ai-inference, coaching, data-platform, cartridge, subscription, community, prescription, measurement, calibration | `postgres.go` 내 CRUD 쿼리, 마이그레이션 스크립트 |
| **A'-2** | CRUD 12서비스 쿼리 구현 | nlp, assistant, iot-gateway, vision, telemedicine, translation, emergency, marketplace, video, device, analytics, shop | 쿼리 구현 + 트랜잭션 처리 |
| **A'-3** | SKELETON 6서비스 쿼리 + 통합 마이그레이션 검증 | concept, audit, cartridge-store, user, voice-profile, digital-twin | docker-compose 통합 테스트, `DB_DSN` 전환 로직 |

**완료 조건**: `docker-compose up` → 전 서비스 PostgreSQL 연결 → 재시작 후 데이터 영속 확인

### 3.3 Phase B': 비즈니스 로직 고도화 (4 Sprint)

SKELETON 6개를 CRUD 이상으로, 기존 CRUD를 REAL 수준으로 끌어올림.

| Sprint | 도메인 | 핵심 작업 |
|--------|--------|----------|
| **B'-1** | AI/의료 | concept→목적별 그룹핑, audit→감사 체인, digital-twin→잔차 모니터링, vision→이미지 분류기 |
| **B'-2** | 커머스 | cartridge-store→수익 분배(70/30), shop→재고관리+프로모션, user→프로필 고도화, voice-profile→TTS 파이프라인 |
| **B'-3** | 소통/분석 | nlp→의도분류+엔티티 인식, assistant→14도메인 액션매핑, translation→의학용어 특화, analytics→코호트 분석 |
| **B'-4** | IoT/응급 | iot-gateway→MQTT+텔레메트리, emergency→에스컬레이션 체인, telemedicine→WebRTC 세션 관리, device→OTA 전체 흐름 |

### 3.4 Phase C: 외부 서비스 연동 (6 Sprint)

| Sprint | 연동 | 서비스 | 상세 |
|--------|------|--------|------|
| **C-1** | 소셜 로그인 + 푸시 | auth, notification | Google/Apple/Naver OAuth, FCM/APNS |
| **C-2** | 결제 게이트웨이 | payment, subscription | Toss 실거래, 빌링키, 세금계산서 |
| **C-3** | 화상진료 + 실시간 | telemedicine, gateway | Agora/Twilio WebRTC, WebSocket |
| **C-4** | AI/Vision/번역 API | vision, assistant, translation, nlp | OpenAI, Google Translate, Azure Cognitive |
| **C-5** | 건강 플랫폼 | Flutter(iOS/Android), health-record | HealthKit, Google Health Connect, FHIR R4, LOINC |
| **C-6** | 인프라 연동 | video, analytics, ai-inference, emergency, auth | MinIO, Elasticsearch, Milvus, 119 API, PASS 본인인증 |

### 3.5 Phase D': 네이티브 통합 (3 Sprint)

| Sprint | 작업 | 상세 |
|--------|------|------|
| **D'-1** | Rust FFI 활성화 | flutter_rust_bridge codegen 실행, Android NDK/iOS 크로스컴파일, `rust_ffi_stub.dart` → 실제 바인딩 전환 |
| **D'-2** | BLE/NFC 실연동 | flutter_blue_plus→Rust BLE 파서, nfc_manager→Rust NFC 파서, 카트리지 자동인식, 측정 데이터 스트림 |
| **D'-3** | 오프라인 + 센서 | CRDT 동기화 엔진 활성화, Hive 로컬 DB, 카메라(음식분석), 72시간 오프라인 테스트 |

### 3.6 Phase E: AI/ML 파이프라인 (4 Sprint)

| Sprint | 작업 | 상세 |
|--------|------|------|
| **E-1** | 엣지 AI 추론 | TFLite 5종 모델 통합, 건강 점수(0~100), Z-score 이상 탐지 |
| **E-2** | 예측/추천 엔진 | 7일/30일 시계열 예측, 개인화 추천, 위험도 분류 |
| **E-3** | 448→896차원 확장 | 핑거프린트 다중모드 융합, 비표적 분석(Milvus), 전자코/전자혀 통합 |
| **E-4** | 연합학습 + 데이터 플랫폼 | Flower 프레임워크, K-익명성, 연구용 API |

### 3.7 Phase F': 테스트/품질 (3 Sprint)

| Sprint | 목표 | 상세 |
|--------|------|------|
| **F'-1** | Go 커버리지 70%+ | REAL 15서비스 엣지케이스, CRUD 비즈니스 규칙, SKELETON 도메인 테스트 |
| **F'-2** | Flutter 커버리지 60%+ | 위젯 테스트 20+, 화면 통합 10+, Golden 스냅샷 5+ |
| **F'-3** | E2E + 통합 | Gateway E2E 전체 플로우, Rust-Flutter-Go 전계층 파이프라인 |

### 3.8 Phase G: 운영/배포 (2 Sprint)

| Sprint | 작업 | 상세 |
|--------|------|------|
| **G-1** | CI/CD 파이프라인 | GitHub Actions 강화, Docker 이미지 빌드, ECR 푸시 |
| **G-2** | K8s 배포 + 모니터링 | Dev/Staging/Prod 환경, Prometheus+Grafana 실운영 |

### 3.9 Phase H: 규제/인증 (2 Sprint)

| Sprint | 작업 | 상세 |
|--------|------|------|
| **H-1** | MFDS/FDA 510(k) | IEC 62304 소프트웨어 수명주기, ISO 14971 위험관리, DHF |
| **H-2** | CE-IVDR + 글로벌 | EU IVDR 인증, WCAG 2.1 AA 접근성, GDPR/개인정보 보호 |

---

## 제4부: Tech Spec v2.4.3 대비 미구현 모듈 분석

| 사양서 섹션 | 모듈 | 현재 상태 | 우선순위 |
|------------|------|----------|:---:|
| §4.1-4.3 | `harness/` (SensorTrait, CartridgeManifest, AfeRegistry) | 미구현 | Phase D |
| §6.1-6.3 | `signal/` (differential, dsp, peak_detector, baseline) | 기본 모듈만 | Phase E |
| §6.4 | `digital_twin/` (twin_engine, drift_detector) | 미구현 | Phase E |
| §6.5-6.7 | `ai/` 고도화 (xgboost, non_target, model_registry) | 기본 모듈만 | Phase E |
| §6.8-6.10 | `checkup/` (session, disease_risk, auto_classifier) | 미구현 | Phase B |
| §6.11-6.14 | `context/` (engine, nutrition, shopping, habit_tracker) | 미구현 | Phase B |
| §6.15 | `diagnostics/` (SelfDiagnostics 6종) | 미구현 | Phase D |
| §6.16 | `healing/` (SelfHealingOrchestrator 4종) | 미구현 | Phase D |

---

## 제5부: 의존관계 맵 및 병렬화 전략

```
시간축 →

Sprint 1-3:   [Phase A': 영속화 완성]  ←── 최우선, 모든 것의 기반
              ‖ 병렬
Sprint 1-3:   [Phase D'-1: Rust FFI]   ←── A'과 독립

Sprint 4-7:   [Phase B': 로직 고도화]  ←── A' 완료 후
              ‖ 병렬
Sprint 4-9:   [Phase C: 외부 연동]     ←── A' 완료 후
              ‖ 병렬
Sprint 4-6:   [Phase D'-2,3: BLE/NFC]  ←── D'-1 완료 후

Sprint 8-11:  [Phase E: AI/ML]         ←── B', C 부분 완료 후
              ‖ 병렬 (F는 전 과정에서 병렬)
Sprint 1-∞:   [Phase F': 테스트]       ←── 상시 병렬 진행

Sprint 12-13: [Phase G: 운영/배포]     ←── F' 기본 완료 후
Sprint 14-15: [Phase H: 규제/인증]     ←── G 완료 후
```

**최적 병렬화 시 총 일정: 15 Sprint** (순차 시 27 Sprint)

---

## 제6부: 핵심 리스크 및 완화 전략

| 리스크 | 영향 | 완화 전략 |
|--------|------|----------|
| PostgreSQL 마이그레이션 중 데이터 정합성 | 높음 | `DB_DSN` 환경변수 전환 패턴으로 메모리/DB 공존 |
| rustc 1.93.0 ICE 버그 | 중간 | `RUSTFLAGS='-A warnings'` 유지, rustc 업데이트 모니터링 |
| 외부 API 의존성 (OpenAI, Toss, Agora) | 높음 | 인터페이스 추상화 + 폴백 스텁 유지, 멀티벤더 전략 |
| 규제 인증 장기 소요 | 높음 | Phase H 조기 착수(DHF 문서화), 사전인증 플랫폼 설계 활용 |
| Flutter 화면 수 축소 (97→72) | 낮음 | UX 사이트맵 기준 누락 화면 검증 필요 |

---

## 제7부: 종합 결론

### 현재 위치: **53%** (2개월 전 44%에서 +9%p)

### 핵심 성과
1. PostgreSQL 전 서비스 구조 생성 — Phase A의 가장 큰 걸림돌 해소
2. Go 서비스 REAL 15개 달성 — 비즈니스 로직 심도 대폭 향상
3. Flutter 테스트 +70% — 품질 기반 강화
4. Rust 테스트 +45% — 코어 엔진 안정성 확보

### 100% 달성까지의 핵심 과제 (우선순위 순)
1. **DB 영속화 쿼리 구현** (Phase A') — 75% 잔여, 모든 것의 기반
2. **외부 연동 13종** (Phase C) — 70% 잔여, 사용자 경험의 핵심
3. **Rust FFI 실활성화** (Phase D') — 55% 잔여, 네이티브 측정의 핵심
4. **Tech Spec 누락 모듈** — harness, digital_twin, checkup, context, diagnostics, healing
5. **AI/ML 파이프라인** (Phase E) — 엣지 추론 + 예측 + 896차원
6. **테스트 80% 커버리지** (Phase F') — 38% 잔여
7. **K8s 실배포** (Phase G) — 구조 완료, 실행 필요
8. **규제 인증** (Phase H) — MFDS/FDA/CE-IVDR
