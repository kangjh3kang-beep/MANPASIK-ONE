# Phase B — SKELETON→REAL 비즈니스 로직 고도화 계획서

**문서 ID**: MPS-PLAN-PHASE-B-v1.0
**작성일**: 2026-04-25
**상태**: Phase B-1~B-3 구현 완료

---

## 1. 실사 결과 요약

### 1.1 SKELETON 6개 서비스 재분류

초기에 SKELETON으로 분류된 6개 서비스를 실사한 결과:

| 서비스 | 초기 분류 | 실사 결과 | 비즈니스 로직 | 테스트 |
|--------|----------|----------|-------------|--------|
| analytics-service | SKELETON | **REAL** | 6메서드(이벤트추적/퍼널/리텐션) | 27개 PASS |
| emergency-service | SKELETON | **REAL** | 10메서드(응급신고/119dispatch/심각도분류) | 21개 PASS |
| assistant-service | SKELETON | **80% REAL** | 7메서드(세션관리/AI응답/인텐트분류) | 테스트 있음 |
| concept-service | SKELETON | **60% REAL** | 14메서드(CRUD+멤버십) | 테스트 있음 |
| data-platform-service | SKELETON | **REAL (Phase E-4)** | 20+메서드(K-Anonymity/연합학습) | 88개 PASS |
| gateway | SKELETON | **80% REAL** | 12/16 라우트 연결, 4개 STUB | 2개 E2E |

### 1.2 핵심 발견
- **data-platform-service**는 이미 Phase E-4 수준 (K-Anonymity, FedAvg 연합학습 구현)
- **analytics/emergency**는 비즈니스 로직은 REAL이나 REST API 노출이 부족
- **gateway**는 핸들러 코드는 완성, 서비스 연결만 누락

---

## 2. Phase B 구현 완료 내역

### B-1: Gateway 서비스 연결 (완료)

`backend/services/gateway/cmd/main.go` 수정:
- ServiceClients 구조체에 7개 필드 추가 (Assistant, Vision, Concept, Organization, DataPlatform, DataProvision, VoiceProfile)
- connectServices()에 5개 서비스 연결 추가 (assistant:50074, vision:50075, concept:50076, data-platform:50077, voice-profile:50078)
- main()에서 5개 setter 호출 (SetAssistantClient, SetVisionClient, SetConceptClients, SetDataPlatformClients, SetVoiceProfileClient)

**결과**: 16/16 라우트 모두 서비스 연결됨

### B-2: Analytics-Service REST API (완료)

`backend/services/analytics-service/internal/handler/rest.go` 신규:
- `POST /api/v1/analytics/events` — 이벤트 기록
- `GET /api/v1/analytics/events` — 최근 이벤트 목록
- `GET /api/v1/analytics/users/{userID}` — 사용자 분석 요약
- `GET /api/v1/analytics/daily/{date}` — 일별 통계
- `GET /api/v1/analytics/retention/{userID}` — 유지율
- `POST /api/v1/analytics/funnel` — 이벤트 퍼널

`backend/services/analytics-service/cmd/main.go` 수정:
- handler 패키지 import + RegisterRoutes 호출

### B-3: Emergency-Service REST API 고도화 (완료)

`backend/services/emergency-service/internal/handler/grpc.go` 전면 재작성:
- `POST /api/v1/emergency/report` — 응급 신고 (기존)
- `POST /api/v1/emergency/resolve` — 응급 해결 (신규)
- `GET /api/v1/emergency/history` — 이력 조회 (신규)
- `GET /api/v1/emergency/contacts` — 연락처 조회 (기존)
- `POST /api/v1/emergency/contacts` — 연락처 추가 (신규)
- `DELETE /api/v1/emergency/contacts` — 연락처 삭제 (신규)
- `GET /api/v1/emergency/settings` — 설정 조회 (분리)
- `PUT /api/v1/emergency/settings` — 설정 변경 (분리)

---

## 3. 검증 결과

| 항목 | 결과 |
|------|------|
| 전체 빌드 | **35/35 PASS** |
| 전체 테스트 | **35/35 PASS** |

---

## 4. SKELETON 재분류 결과 (Phase B 후)

| 서비스 | 분류 변경 | 근거 |
|--------|----------|------|
| analytics-service | SKELETON → **REAL** | REST API 6개 + 비즈니스 로직 6메서드 |
| emergency-service | SKELETON → **REAL** | REST API 8개 + 119 dispatch + notifier |
| data-platform-service | SKELETON → **REAL** | 20+ 메서드, K-Anonymity, 연합학습 |
| gateway | SKELETON → **REAL** | 16/16 라우트 연결 완료 |
| assistant-service | SKELETON → **CRUD+** | 7메서드, OpenAI 통합, 세션관리 |
| concept-service | SKELETON → **CRUD** | 14메서드, 통계/대시보드 stub |

### 최종 Go 서비스 분류 (Phase B 후)

- **REAL (21)**: auth, health-record, family, notification, admin, payment, coaching, community, prescription, reservation, telemedicine, translation, video, calibration, cartridge-store, **analytics**, **emergency**, **data-platform**, **gateway**
- **CRUD+ (10)**: user, subscription, shop, device, measurement, cartridge, ai-inference, marketplace, iot-gateway, **assistant**
- **CRUD (4)**: nlp, vision, digital-twin, **concept**
- **기타 (2)**: audit, voice-profile

---

*자동 생성: 2026-04-25 | 만파식(萬波息) Phase B 비즈니스 로직 고도화 계획서*
