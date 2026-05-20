# 만파식(ManPaSik) 100% 달성 구현 계획안 v1.0

> **작성일**: 2026-02-23
> **기준**: 현재 구현 상태 정밀 측정 후 수립
> **목표**: 전 차원 100% 달성

---

## 1. 현재 상태 정밀 측정 (2026-02-23 기준)

### 1.1 자산 인벤토리

| 계층 | 측정치 | 상세 |
|------|--------|------|
| **Proto** | 33 서비스, 233 RPC, 4,274줄 | 구조 100% |
| **Go 백엔드** | 33 서비스 + 1 게이트웨이 | 빌드/테스트 ALL PASS |
| **Gateway** | 16 라우트 파일, 240+ 엔드포인트 | 구조 100% |
| **DB 스키마** | 34 SQL 스크립트, 112 테이블 | 구조 100% |
| **Rust 엔진** | 6,508줄, 62 테스트 | BLE/NFC/Crypto/DSP |
| **Flutter 화면** | 97 프리젠테이션 파일 | UI 100% |
| **Flutter REST** | 230 메서드, 2,791줄 | 구조 100% |
| **Flutter 위젯** | 49 커스텀 위젯 | Sanggam Orbit 디자인 |
| **Flutter 저장소** | 16 REST 저장소 | DioException 패턴 |
| **Flutter 테스트** | 30 파일, 259 테스트 | 저장소 중심 |
| **Flutter 라우트** | 90 라우트 정의 | RBAC 가드 포함 |
| **기획 문서** | 135 마크다운 | 11-Layer 설계 |

### 1.2 Go 서비스 비즈니스 로직 깊이 분류

| 분류 | 서비스 (개수) | 특징 |
|------|--------------|------|
| **REAL (6)** | auth, health-record, family, notification, admin, payment | JWT/FHIR/RBAC 실로직 |
| **CRUD (15)** | user, subscription, shop, device, measurement, cartridge, calibration, coaching, reservation, prescription, community, video, telemedicine, translation, cartridge-store | 기본 CRUD+검증 |
| **SKELETON (12)** | ai-inference, analytics, emergency, marketplace, nlp, iot-gateway, vision, assistant, concept, voice-profile, data-platform, gateway | 인터페이스만 |

### 1.3 외부 연동 현황

| 연동 | 현재 상태 | 필요 작업 |
|------|----------|-----------|
| 카카오 OAuth | **REAL** | 완료 |
| Toss 결제 PG | **REAL** (콜백 구조) | 실제 키 연동 |
| Kafka/Redpanda | **REAL** | 완료 |
| Redis | **REAL** | 완료 |
| PostgreSQL/TimescaleDB | **설계만** | 메모리→DB 마이그레이션 |
| FHIR R4 | **REAL** | LOINC 매핑 보강 |
| WebRTC (Agora/Twilio) | **STUB** | 서버 연동 필요 |
| Vision API (OpenAI) | **STUB** | 모델 연동 필요 |
| 번역 API (Google/Azure) | **STUB** | 엔진 연동 필요 |
| FCM/APNS 푸시 | **STUB** | Firebase 연동 필요 |
| HealthKit/Google Health | **STUB** | SDK 연동 필요 |
| Rust FFI (BLE/NFC) | **STUB** | flutter_rust_bridge 활성화 |
| Google/Apple/Naver OAuth | **STUB** | API 연동 필요 |

### 1.4 완성도 매트릭스 (현재 → 목표)

| 차원 | 현재 | 목표 | 갭 |
|------|------|------|-----|
| 구조적 완성도 | 95% | 100% | 5% |
| Go 비즈니스 로직 | 40% | 100% | 60% |
| Flutter 데이터 바인딩 | 87% | 100% | 13% |
| 외부 연동 | 25% | 100% | 75% |
| 테스트 커버리지 | 35% | 80% | 45% |
| DB 영속화 | 0% | 100% | 100% |
| 네이티브 기능 | 10% | 80% | 70% |
| 운영/배포 | 15% | 90% | 75% |
| 규제 문서 | 80% | 100% | 20% |
| **종합** | **~44%** | **100%** | **~56%** |

---

## 2. 100% 달성 로드맵 — 8 Phase, 32 Sprint

### 전체 구조

```
Phase A: 기반 영속화      ─── S-A1~A4  (4 스프린트)
Phase B: 비즈니스 로직    ─── S-B1~B6  (6 스프린트)
Phase C: 외부 연동        ─── S-C1~C6  (6 스프린트)
Phase D: 네이티브 통합    ─── S-D1~D4  (4 스프린트)
Phase E: AI/ML 파이프라인 ─── S-E1~E4  (4 스프린트)
Phase F: 테스트/품질      ─── S-F1~F4  (4 스프린트)
Phase G: 운영/배포        ─── S-G1~G2  (2 스프린트)
Phase H: 규제/인증        ─── S-H1~H2  (2 스프린트)
```

### 의존관계 맵

```
Phase A (영속화) ────────────────┐
                                 ├──→ Phase B (비즈니스 로직)
Phase C (외부 연동) ─── 병렬 ───┤
Phase D (네이티브) ──── 병렬 ───┘
                                 │
Phase E (AI/ML) ←── B, C 완료 후 ┘
                                 │
Phase F (테스트) ←── A~E 병렬 진행 가능
Phase G (운영) ←── F 완료 후
Phase H (규제) ←── G 완료 후
```

---

## Phase A: 기반 영속화 (메모리 → PostgreSQL)

> **현재**: 33개 서비스가 모두 `internal/repository/memory/` 사용 → 재시작 시 데이터 소실
> **목표**: 모든 서비스가 PostgreSQL/TimescaleDB 영속 저장소 사용

### Sprint A-1: 저장소 추상화 레이어 + 핵심 6 서비스 DB 연동

| 작업 | 파일 | 상세 |
|------|------|------|
| Repository 인터페이스 추출 | `backend/shared/repository/interfaces.go` | 모든 서비스 공통 패턴 추출 |
| PostgreSQL 드라이버 래퍼 | `backend/shared/database/postgres.go` | pgx 풀링, 마이그레이션, 헬스체크 |
| auth-service DB 저장소 | `internal/repository/postgres/auth.go` | users, sessions, tokens |
| user-service DB 저장소 | `internal/repository/postgres/user.go` | profiles, health_profiles |
| payment-service DB 저장소 | `internal/repository/postgres/payment.go` | payments, refunds |
| admin-service DB 저장소 | `internal/repository/postgres/admin.go` | audit_logs, settings |
| health-record-service DB | `internal/repository/postgres/healthrecord.go` | records, consents |
| family-service DB 저장소 | `internal/repository/postgres/family.go` | groups, members |
| 환경변수 전환 로직 | 각 서비스 `cmd/main.go` | `DB_DSN` 존재 시 postgres, 미존재 시 memory |

### Sprint A-2: CRUD 15 서비스 DB 연동 (배치 1)

| 서비스 | DB 테이블 | 비고 |
|--------|-----------|------|
| subscription-service | subscriptions, plan_features | 티어 기능 매핑 |
| shop-service | products, carts, orders, order_items | 트랜잭션 처리 |
| device-service | devices, device_status, ota_updates | TimescaleDB hypertable |
| measurement-service | measurements, measurement_results | TimescaleDB 시계열 |
| cartridge-service | cartridges, cartridge_usage | NFC 태그 매핑 |
| calibration-service | calibration_data, calibration_models | 모델 버전 관리 |
| coaching-service | health_goals, coaching_sessions | 일일/주간 리포트 |
| notification-service | notifications, preferences | 다중 채널 |

### Sprint A-3: CRUD 15 서비스 DB 연동 (배치 2)

| 서비스 | DB 테이블 | 비고 |
|--------|-----------|------|
| reservation-service | reservations, facilities, doctors | 지역 검색 PostGIS |
| prescription-service | prescriptions, fulfillments | 토큰 시스템 |
| community-service | posts, comments, likes, challenges | 풀텍스트 검색 |
| video-service | videos, video_metadata | MinIO 연동 |
| telemedicine-service | sessions, chat_messages, recordings | WebRTC 세션 |
| translation-service | translations, language_cache | 캐시 레이어 |
| cartridge-store-service | store_listings, purchases, reviews | 마켓플레이스 |

### Sprint A-4: SKELETON 12 서비스 DB 연동 + 마이그레이션 검증

| 서비스 | DB 테이블 | 비고 |
|--------|-----------|------|
| ai-inference-service | ai_models, inference_results | 모델 메타데이터 |
| analytics-service | events, user_stats, aggregations | TimescaleDB |
| emergency-service | alerts, escalations | 긴급 상태 추적 |
| marketplace-service | listings, reviews, purchases | 개발자 마켓 |
| nlp-service | queries, intents, entities | 의도 파싱 캐시 |
| iot-gateway-service | iot_devices, commands, telemetry | TimescaleDB |
| vision-service | food_analyses, meal_logs | 이미지 분석 결과 |
| assistant-service | sessions, turns, preferences | 대화 이력 |
| concept-service | concepts, assignments, organizations | 목적별 컨셉 |
| voice-profile-service | voice_profiles, synthesis_jobs | TTS 모델 관리 |
| data-platform-service | datasets, access_requests, insights | 익명 집계 |
| **검증** | 전체 마이그레이션 | docker-compose 통합 테스트 |

---

## Phase B: 비즈니스 로직 고도화 (SKELETON → REAL)

> **현재**: 12개 SKELETON + 15개 CRUD = 얕은 로직
> **목표**: 전 서비스 도메인-특화 비즈니스 규칙 구현

### Sprint B-1: AI 추론 파이프라인 실구현

| 작업 | 상세 |
|------|------|
| ai-inference-service 고도화 | 88차원→448차원 핑거프린트 분석 엔진 |
| 모델 관리 | TFLite/ONNX 모델 로드, 버전 관리, A/B 테스트 |
| 스코어링 엔진 | 건강 점수(0~100), 카테고리별 세부 점수 |
| 이상 탐지 | Z-score 기반 이상치 자동 감지 |
| 추세 예측 | ARIMA/Prophet 시계열 예측 (7일/30일) |
| 스트림 처리 | 실시간 측정 데이터 → Kafka → 추론 → 결과 발행 |

### Sprint B-2: 의료 도메인 실로직

| 작업 | 상세 |
|------|------|
| telemedicine-service | WebRTC 세션 관리, 녹음/기록, 의사-환자 매칭 |
| reservation-service | PostGIS 지역 검색, 시간 슬롯 관리, 충돌 방지 |
| prescription-service | 처방 토큰 시스템, 약국 전송, QR 코드 생성 |
| health-record-service | LOINC 코드 매핑 (88차원→LOINC), FHIR Bundle 생성 |
| emergency-service | 이상 감지 → 에스컬레이션 체인 → 119 연동 로직 |

### Sprint B-3: 커머스 도메인 실로직

| 작업 | 상세 |
|------|------|
| shop-service | 재고 관리, 프로모션 엔진, 쿠폰 시스템 |
| payment-service | Toss PG 실거래 흐름, 환불 프로세스, 정산 |
| subscription-service | 자동 갱신, 업/다운그레이드, 프로레이트 계산 |
| marketplace-service | 3rd-party 카트리지 등록/심사/배포 |
| cartridge-store-service | 수익 분배(70/30), 개발자 정산, 리뷰 시스템 |

### Sprint B-4: 데이터/분석 도메인 실로직

| 작업 | 상세 |
|------|------|
| analytics-service | 이벤트 집계, 코호트 분석, 유지율 계산 |
| data-platform-service | K-익명성 적용, 지역 집계, 연구 데이터 제공 |
| vision-service | 이미지 분류기, 칼로리 추정, 영양 DB 매핑 |
| coaching-service | 개인화 알고리즘, 목표 달성 예측, 행동 추천 |

### Sprint B-5: 소통/커뮤니티 실로직

| 작업 | 상세 |
|------|------|
| community-service | 랭킹 알고리즘, 스팸 필터, 전문가 인증 |
| translation-service | 번역 품질 평가, 용어 사전, 의학 용어 특화 |
| nlp-service | 의도 분류(건강 질의), 증상 추출, 엔티티 인식 |
| assistant-service | 14개 도메인 액션 매핑, 다중 턴 컨텍스트, 확인 플로우 |
| voice-profile-service | TTS 모델 학습 파이프라인, 음성 합성, 다국어 |

### Sprint B-6: IoT/디바이스 실로직

| 작업 | 상세 |
|------|------|
| iot-gateway-service | MQTT 브로커 연동, 디바이스 명령 큐, 텔레메트리 수집 |
| device-service | OTA 업데이트 전체 흐름, 펌웨어 버전 관리 |
| concept-service | 목적별 디바이스 그룹핑, 대시보드 집계 |
| calibration-service | 팩토리/현장 보정 자동화, 드리프트 감지 |

---

## Phase C: 외부 서비스 연동

> **현재**: 3개 REAL / 10개 STUB
> **목표**: 13개 전체 REAL

### Sprint C-1: 소셜 로그인 + 푸시 알림

| 연동 | 서비스 | 작업 |
|------|--------|------|
| Google OAuth 2.0 | auth-service | ID 토큰 검증, 프로필 싱크 |
| Apple Sign-In | auth-service | JWT 검증, 이메일 릴레이 |
| Naver 로그인 | auth-service | OAuth 2.0 콜백 |
| Firebase Cloud Messaging | notification-service | FCM 토큰 등록, 토픽 발행 |
| APNS (Apple Push) | notification-service | 디바이스 토큰, 사일런트 푸시 |

### Sprint C-2: 결제 게이트웨이

| 연동 | 서비스 | 작업 |
|------|--------|------|
| Toss Payments | payment-service | 실제 API 키, 결제 승인/취소/환불 |
| Toss 빌링키 | subscription-service | 정기결제, 빌링키 발급/갱신 |
| 세금계산서 | payment-service | 전자세금계산서 발행 API |
| PG 웹훅 | gateway | 결제 결과 콜백 수신 |

### Sprint C-3: 화상진료 + 실시간 통신

| 연동 | 서비스 | 작업 |
|------|--------|------|
| Agora/Twilio WebRTC | telemedicine-service | 채널 생성, 토큰 발급, 녹화 |
| WebSocket 서버 | gateway | 실시간 채팅, 디바이스 상태 스트림 |
| SignalR/Socket.IO | community-service | 실시간 포럼 업데이트 |

### Sprint C-4: AI/Vision/번역 API

| 연동 | 서비스 | 작업 |
|------|--------|------|
| OpenAI Vision API | vision-service | 음식 이미지 분석, 칼로리 추정 |
| OpenAI GPT-4 | assistant-service | 자연어 건강 상담, 코칭 |
| Google Translate | translation-service | 8개 언어 실시간 번역 |
| Azure Cognitive | nlp-service | 의도 분류, 엔티티 추출 |

### Sprint C-5: 건강 플랫폼 연동

| 연동 | 서비스 | 작업 |
|------|--------|------|
| Apple HealthKit | Flutter (iOS) | 심박수, 걸음수, 수면 동기화 |
| Google Health Connect | Flutter (Android) | 건강 데이터 읽기/쓰기 |
| FHIR R4 서버 | health-record-service | 외부 EMR 데이터 임포트/익스포트 |
| LOINC 코드 DB | health-record-service | 88차원 → 표준 LOINC 매핑 |

### Sprint C-6: 인프라 연동

| 연동 | 서비스 | 작업 |
|------|--------|------|
| MinIO (S3 호환) | video-service, vision-service | 파일 업로드/다운로드 |
| Elasticsearch | analytics-service | 풀텍스트 검색, 로그 분석 |
| Milvus | ai-inference-service | 벡터 유사도 검색 |
| 119 긴급 API | emergency-service | 긴급 상황 자동 신고 |
| PASS 본인인증 | auth-service | 휴대폰 본인인증 |

---

## Phase D: 네이티브 통합

### Sprint D-1: Rust FFI 활성화

| 작업 | 상세 |
|------|------|
| flutter_rust_bridge 설정 | Cargo.toml, build.rs, codegen 실행 |
| Dart 바인딩 생성 | `rust_ffi_stub.dart` → 실제 바인딩 전환 |
| Android NDK 크로스 컴파일 | aarch64-linux-android, armv7 |
| iOS 빌드 | macOS 환경에서 lipo universal |
| 웹 폴백 | WASM 미지원 → 스텁 모드 유지 |

### Sprint D-2: BLE/NFC 실디바이스 통합

| 작업 | 상세 |
|------|------|
| BLE 스캔 실연동 | flutter_blue_plus → Rust BLE 파서 |
| NFC 태그 읽기 | nfc_manager → Rust NFC 파서 |
| 카트리지 자동 인식 | NFC 태그 → 카트리지 ID → 29종 매핑 |
| 측정 데이터 스트림 | BLE GATT → Rust 차동측정 → Flutter UI |
| 디바이스 OTA | 펌웨어 파일 → BLE 전송 → 상태 추적 |

### Sprint D-3: 센서/카메라 통합

| 작업 | 상세 |
|------|------|
| 카메라 (음식 분석) | camera 패키지 → 이미지 캡처 → Vision API |
| 마이크 (음성 명령) | speech_to_text → NLP → 액션 실행 |
| TTS (음성 응답) | flutter_tts → AI 비서 음성 출력 |
| 생체 센서 | 가속도계/자이로 → 운동 추적 |

### Sprint D-4: 오프라인 + 동기화

| 작업 | 상세 |
|------|------|
| Hive 로컬 DB | 측정 결과, 설정, 캐시 데이터 영속화 |
| CRDT 동기화 엔진 | Rust sync 모듈 → 충돌 해소 → 서버 동기화 |
| 오프라인 72시간 테스트 | 네트워크 차단 → 측정 → 복구 → 동기화 검증 |
| 큐 기반 동기화 | 오프라인 요청 큐 → 네트워크 복구 시 순차 전송 |

---

## Phase E: AI/ML 파이프라인

### Sprint E-1: 엣지 AI 추론

| 작업 | 상세 |
|------|------|
| TFLite 모델 통합 | Rust ai 모듈 → 5종 모델 로드 |
| 건강 점수 알고리즘 | 88차원 → 가중 평균 → 0~100점 |
| 이상 탐지 | Z-score + IQR → 실시간 경고 |
| 엣지→클라우드 파이프라인 | 디바이스 추론 → 결과 전송 → 클라우드 검증 |

### Sprint E-2: 예측/추천 엔진

| 작업 | 상세 |
|------|------|
| 시계열 예측 | 7일/30일 건강 추세 예측 |
| 개인화 추천 | 사용자 이력 기반 운동/식단/수면 추천 |
| 위험도 분류 | 안전/주의/위험 자동 분류 + 에스컬레이션 |
| CoachScreen 실데이터 | predictTrend() → 실제 API → Spirit Orb 반응 |

### Sprint E-3: 448→896차원 확장

| 작업 | 상세 |
|------|------|
| 448차원 핑거프린트 | Rust fingerprint 모듈 확장 |
| 896차원 정밀 분석 | E12-IF 다중 채널 융합 |
| 비표적 분석 (Untargeted) | 미지 물질 탐지 + 유사도 검색 (Milvus) |
| 전자코/전자혀 센서 | 추가 센서 데이터 통합 프로토콜 |

### Sprint E-4: 연합학습 + 데이터 플랫폼

| 작업 | 상세 |
|------|------|
| 연합학습 프레임워크 | 로컬 학습 → 모델 가중치만 전송 → 글로벌 집계 |
| K-익명성 데이터 | 지역별 익명 건강 통계 |
| 연구용 API | 대학/연구소 데이터 접근 API |
| 음성 복제 | TTS 모델 학습 → 다국어 음성 합성 |

---

## Phase F: 테스트/품질 (현재 35% → 80%)

### Sprint F-1: Go 단위 테스트 강화

| 작업 | 현재 | 목표 |
|------|------|------|
| REAL 6 서비스 | 기본 테스트 | 엣지케이스 + 실패 시나리오 |
| CRUD 15 서비스 | 기본 CRUD | 비즈니스 규칙 검증 |
| SKELETON 12 서비스 | 최소 테스트 | 도메인 로직 테스트 |
| 목표 | 49 패키지 | 커버리지 70%+ |

### Sprint F-2: Flutter 위젯/화면 테스트

| 작업 | 현재 | 목표 |
|------|------|------|
| 저장소 테스트 | 14/16 | 16/16 |
| 위젯 테스트 | 0 | 핵심 위젯 20개 |
| 화면 통합 테스트 | 0 | 주요 플로우 10개 |
| Golden 테스트 | 0 | 핵심 화면 스냅샷 5개 |

### Sprint F-3: E2E + 통합 테스트

| 작업 | 상세 |
|------|------|
| 인증 E2E | 회원가입 → 로그인 → 토큰 갱신 → 로그아웃 |
| 측정 E2E | 디바이스 페어링 → 측정 → 결과 조회 → 히스토리 |
| 결제 E2E | 상품 탐색 → 장바구니 → 결제 → 주문 확인 |
| 의료 E2E | 시설 검색 → 예약 → 화상진료 → 처방 |
| 오프라인 E2E | 오프라인 측정 → 네트워크 복구 → 동기화 |
| Gateway E2E | 전 엔드포인트 정상 응답 확인 (240+) |

### Sprint F-4: 성능/보안 테스트

| 작업 | 상세 |
|------|------|
| 부하 테스트 | k6/Locust — 1,000 동시 사용자 |
| 레이턴시 벤치마크 | gRPC P95 < 100ms, REST P95 < 200ms |
| 보안 침투 테스트 | OWASP Top 10, SQL injection, XSS |
| SAST 정적 분석 | GoSec, SonarQube, Snyk |
| 의존성 감사 | npm audit, go mod tidy, cargo audit |

---

## Phase G: 운영/배포

### Sprint G-1: CI/CD + 컨테이너 오케스트레이션

| 작업 | 상세 |
|------|------|
| GitHub Actions | PR → 빌드 → 테스트 → 린트 자동화 |
| Docker Compose (prod) | 33 서비스 + PostgreSQL + Kafka + Redis |
| Kubernetes 매니페스트 | Helm 차트, 네임스페이스, HPA |
| 시크릿 관리 | Vault/SealedSecrets |
| 블루-그린 배포 | 무중단 배포 전략 |

### Sprint G-2: 모니터링 + 옵저버빌리티

| 작업 | 상세 |
|------|------|
| Prometheus 메트릭 | 서비스별 커스텀 메트릭, SLI/SLO |
| Grafana 대시보드 | 서비스 건강도, 비즈니스 KPI |
| Jaeger 분산 추적 | gRPC → REST 전체 호출 추적 |
| ELK 스택 | 구조화 로깅, 에러 알림 |
| PagerDuty 알림 | 에러율/레이턴시 임계치 초과 알림 |

---

## Phase H: 규제/인증

### Sprint H-1: 의료기기 인허가 문서

| 작업 | 상세 |
|------|------|
| IEC 62304 완성 | 소프트웨어 설계 이력 파일 (DHF) 마무리 |
| ISO 14971 보강 | 위험 관리 보고서 최종 검토 |
| MFDS 기술문서 | 한국 식약처 3등급 의료기기 기술문서 |
| FDA 510(k) 준비 | Predicate Device 조사, SE Report |
| CE-IVDR 대응 | EU IVD 규정 기술문서 |

### Sprint H-2: 데이터 보호/감사

| 작업 | 상세 |
|------|------|
| GDPR DPIA | 데이터 보호 영향 평가 |
| HIPAA BAA | 사업 관련 계약 체결 |
| SOC 2 Type II | 감사 보고서 준비 |
| 개인정보 암호화 검증 | AES-256-GCM 키 관리 검증 |
| 데이터 보존/삭제 정책 | GDPR Art.17 (잊힐 권리) 구현 |

---

## 3. 정량적 산출물 요약

| 항목 | 현재 | 100% 달성 시 | 증분 |
|------|------|-------------|------|
| Go 서비스 (REAL) | 6 | 33 | +27 고도화 |
| DB 저장소 (postgres/) | 0 | 33 | +33 신규 |
| 외부 연동 | 3 REAL | 13 REAL | +10 연동 |
| Flutter 테스트 | 30 파일 | 80+ 파일 | +50 파일 |
| Go 테스트 커버리지 | ~30% | 70%+ | +40%p |
| Rust FFI | STUB | REAL | 활성화 |
| AI 모델 | 0 배포 | 5+ 배포 | +5 모델 |
| K8s 매니페스트 | base | production | Helm 차트 |
| 규제 문서 | 80% | 100% | +7 문서 |

---

## 4. 우선순위별 실행 권고

### Tier 1: 즉시 실행 (Claude Code에서 구현 가능)

| 순위 | 작업 | Phase | 영향도 |
|------|------|-------|--------|
| 1 | SKELETON 12 서비스 → CRUD+ 비즈니스 로직 | B | HIGH |
| 2 | Flutter 위젯/화면 테스트 추가 | F-2 | HIGH |
| 3 | CoachScreen REST 실연동 | B-1 | MEDIUM |
| 4 | 미사용 REST 메서드 → 화면 연동 | - | MEDIUM |

### Tier 2: 환경/키 필요 (외부 의존)

| 순위 | 작업 | Phase | 필요 사항 |
|------|------|-------|-----------|
| 1 | PostgreSQL 마이그레이션 | A | Docker PostgreSQL 인스턴스 |
| 2 | Toss 결제 실연동 | C-2 | Toss 테스트 API 키 |
| 3 | Firebase FCM | C-1 | Firebase 프로젝트 설정 |
| 4 | WebRTC (Agora) | C-3 | Agora App ID |
| 5 | OpenAI Vision | C-4 | OpenAI API 키 |

### Tier 3: 하드웨어/환경 필요

| 순위 | 작업 | Phase | 필요 사항 |
|------|------|-------|-----------|
| 1 | Rust FFI (Android) | D-1 | Android NDK, 실기기 |
| 2 | Rust FFI (iOS) | D-1 | macOS + Xcode |
| 3 | BLE/NFC 실디바이스 | D-2 | 만파식 리더기 프로토타입 |
| 4 | HealthKit | C-5 | iPhone + 개발자 계정 |

---

## 5. 100% 달성 기준 정의

| 차원 | 100% 기준 | 검증 방법 |
|------|-----------|-----------|
| **구조** | 33 서비스 Proto/DB/Dockerfile/Gateway 완비 | `go build`, `flutter analyze` |
| **비즈니스 로직** | 33 서비스 모두 REAL 수준 | 코드 리뷰, 커버리지 |
| **외부 연동** | 13 외부 서비스 모두 REAL 연동 | 통합 테스트 PASS |
| **네이티브** | Rust FFI + BLE/NFC + 오프라인 72h | 실디바이스 테스트 |
| **테스트** | Go 70%+, Flutter 60%+, E2E 10 시나리오 | CI 리포트 |
| **운영** | K8s 프로덕션, CI/CD, 모니터링 | 배포 성공 |
| **규제** | 5국 인허가 서류 100% | 규제 체크리스트 |
| **문서** | 135 → 150+ 기술 문서 | 트레이서빌리티 매트릭스 |

---

## 변경 이력

| 날짜 | 버전 | 변경 내용 |
|------|------|-----------|
| 2026-02-23 | v1.0 | 초판 작성 — 현재 44% → 100% 로드맵 |
