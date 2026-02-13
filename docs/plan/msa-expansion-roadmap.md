# 단계별 MSA 확장 로드맵

**문서번호**: MPK-PLAN-MSA-v1.0  
**기준**: 원본 기획안 IV. 4.1 마이크로서비스 도메인 맵  
**목적**: Phase별 서비스 도입 시점·담당 서비스명의 공식 기준

---

## 1. Phase 1 (MVP, Month 1–4) — 현행

| 서비스 | 도메인 | 상태 | 비고 |
|--------|--------|------|------|
| auth-service | 사용자 | 구현 중 | Keycloak OIDC, JWT, 회원가입/로그인 |
| user-service | 사용자 | 구현 중 | 프로필, 구독, 가족 그룹 |
| device-service | 측정 | 구현 중 | 리더기 등록/목록/OTA |
| measurement-service | 측정 | 구현 중 | 세션, 스트림, TimescaleDB·Milvus |

**Proto**: `manpasik.proto` (MeasurementService, DeviceService, UserService). auth는 Keycloak·gRPC 인터셉터.

---

## 2. Phase 2 (Core, Month 5–8)

| 서비스 | 도메인 | 도입 시점 | 비고 |
|--------|--------|-----------|------|
| subscription-service | 커머스 | Month 5 | SaaS 구독 관리, 티어 업/다운 |
| shop-service | 커머스 | Month 5 | 상품, 장바구니, 주문 |
| payment-service | 커머스 | Month 6 | PG 연동, 구독 결제 |
| ai-inference-service | AI | Month 6 | 실시간 추론, 모델 서빙 |
| coaching-service | AI | Month 7 | AI 건강 코칭, 개인화 추천 |
| cartridge-service | 측정 | Month 6 | 카트리지 인증, 보정 데이터, 사용 추적 |
| calibration-service | 측정 | Month 7 | 보정 모델 관리, 팩토리 보정 |

---

## 3. Phase 3 (Advanced, Month 9–12)

| 서비스 | 도메인 | 도입 시점 | 비고 |
|--------|--------|-----------|------|
| family-service | 사용자 | Month 9 | 가족 그룹, 보호자, 공유 설정 (또는 user-service 확장) |
| health-record-service | 의료 | Month 9 | 건강 기록, FHIR 호환 |
| telemedicine-service | 의료 | Month 10 | 화상진료, 의료진 매칭 |
| reservation-service | 의료 | Month 10 | 병원/약국 예약 |
| prescription-service | 의료 | Month 11 | 처방전, 약물 상호작용 |
| community-service | 커뮤니티 | Month 10 | 포럼, Q&A, 챌린지 |
| translation-service | 커뮤니티 | Month 11 | 실시간 번역, 자막 |
| video-service | 커뮤니티 | Month 11 | WebRTC 시그널링, 미디어 |
| notification-service | 커뮤니티 | Month 9 | 푸시, 이메일, SMS, 인앱 |
| admin-service | 관리 | Month 10 | 계층형 관리자, 권한 |

---

## 4. Phase 4 (Ecosystem, Month 13–18)

| 서비스 | 도메인 | 도입 시점 | 비고 |
|--------|--------|-----------|------|
| marketplace-service | 커머스 | Month 13 | SDK 마켓, 서드파티 카트리지 |
| ai-training-service | AI | Month 14 | 모델 학습, 연합학습 코디네이션 |
| vision-service | AI | Month 14 | 음식 사진, 칼로리 분석 |
| nlp-service | AI | Month 15 | 자연어, 번역, 음성 |
| inventory-service | 관리 | Month 14 | 재고, 공급망 추적 |
| logistics-service | 관리 | Month 15 | 배송, 추적 |
| analytics-service | 관리 | Month 15 | 비즈니스 인텔리전스 |
| iot-gateway-service | IoT | Month 16 | 디바이스 연결, MQTT 브로커 |
| location-service | IoT | Month 16 | 리더기 위치, 지오펜싱 |
| emergency-service | IoT | Month 17 | 긴급 대응, 119 연동, AI 통화 |

---

## 5. 원본 4.1 도메인 맵과의 대응

| 원본 도메인 | Phase 1 | Phase 2 | Phase 3 | Phase 4 |
|-------------|---------|---------|---------|---------|
| 사용자 | auth, user | — | family(또는 user 확장) | — |
| 측정 | device, measurement | cartridge, calibration | — | — |
| AI | — | ai-inference, coaching | — | ai-training, vision, nlp |
| 의료 | — | — | health-record, telemedicine, reservation, prescription | — |
| 커머스 | — | subscription, shop, payment | — | marketplace |
| 관리 | — | — | admin | inventory, logistics, analytics |
| 커뮤니티 | — | — | community, translation, video, notification | — |
| IoT | — | — | — | iot-gateway, location, emergency |

---

**참조**: 원본 기획안 IV. 4.1, `backend/services/`, `docs/plan/MPK-ECO-PLAN-v1.1-COMPLETE.md`
