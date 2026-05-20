# Phase C — 외부 연동 통합 계획서

**문서 ID**: MPS-PLAN-PHASE-C-v1.0
**작성일**: 2026-04-25
**상태**: Phase C-1~C-3 구현 완료

---

## 1. 실사 결과 요약

### 1.1 외부 연동 현황 (Phase C 이전)

초기 MEMORY.md에서 외부 연동 완성도를 25%로 추정했으나, 정밀 실사 결과 **~65%** 이미 구현 완료 확인:

| 연동 영역 | 구현체 | 파일 | 상태 |
|-----------|--------|------|------|
| OAuth (Google/Apple/Kakao/Naver) | 4개 Provider + 테스트 | `auth-service/internal/oauth/` | **100%** |
| Redis 캐시/토큰 | go-redis/v9 + TokenRepository | `shared/cache/redis.go` | **100%** |
| Kafka 이벤트버스 | franz-go Producer/Consumer | `shared/events/kafka_adapter.go` | **90%** |
| S3/MinIO 스토리지 | minio-go Upload/Download | `shared/storage/s3.go` | **100%** |
| Milvus 벡터DB | milvus-sdk-go v2.4.2 | `shared/vectordb/milvus.go` | **100%** |
| Elasticsearch 검색 | HTTP Client Search/Index | `shared/search/elasticsearch.go` | **100%** |
| Toss Payments PG | Confirm/Cancel | `payment-service/internal/pg/toss.go` | **90%** |
| Agora WebRTC | Token Provider | `telemedicine-service/internal/webrtc/agora.go` | **100%** |
| OpenAI Vision | Image Analyzer | `vision-service/internal/vision/openai.go` | **100%** |
| OpenAI NLP | Intent Classifier | `nlp-service/internal/classifier/openai.go` | **100%** |
| FCM 푸시 | ServerKey HTTP | `notification-service/internal/push/fcm.go` | **100%** |
| K-Anonymity/연합학습 | FedAvg + Privacy | `data-platform-service/` | **100%** |

### 1.2 실제 GAP 3가지

1. **Payment Webhook**: PG사(Toss)로부터의 비동기 결제 확인 수신 핸들러 부재
2. **검색 API 미노출**: Milvus 벡터 유사도 검색, ES 전문 검색이 REST로 노출되지 않음
3. **Kafka 이벤트 소비자**: 4개 서비스가 이벤트 발행하지만 소비하는 서비스 0개

---

## 2. Phase C 구현 완료 내역

### C-1: Payment REST API + Toss Webhook (완료)

**신규 파일**: `backend/services/payment-service/internal/handler/rest.go`
- `POST /api/v1/payments` — 결제 생성
- `POST /api/v1/payments/confirm` — 결제 확인 (Toss Confirm API)
- `GET /api/v1/payments/{paymentID}` — 결제 조회
- `GET /api/v1/payments` — 결제 목록 (사용자별)
- `POST /api/v1/payments/{paymentID}/refund` — 환불 처리
- `POST /webhooks/payments/toss` — **Toss Webhook** (HMAC-SHA256 서명 검증)

**수정 파일**:
- `payment-service/internal/service/payment.go` — `GetByOrderID` 인터페이스 추가 + `ConfirmPaymentByWebhook` 메서드 (멱등성 보장)
- `payment-service/internal/repository/memory/payment.go` — `GetByOrderID` 구현
- `payment-service/internal/repository/postgres/payment.go` — `GetByOrderID` 구현
- `payment-service/internal/service/payment_test.go` — fakePayRepo에 `GetByOrderID` 추가
- `payment-service/cmd/main.go` — REST 핸들러 등록 + TOSS_WEBHOOK_SECRET 환경변수

### C-2: Measurement/Community 검색 REST API (완료)

**신규 파일**: `backend/services/measurement-service/internal/handler/rest.go`
- `GET /api/v1/measurements/history` — 측정 이력 조회 (user_id, limit, offset)
- `POST /api/v1/measurements/search/similar` — **Milvus 벡터 유사도 검색** (핑거프린트)

**수정 파일**:
- `measurement-service/internal/service/measurement.go` — `SearchSimilarFingerprints` 메서드 추가
- `measurement-service/cmd/main.go` — REST 핸들러 등록
- `community-service/internal/service/community.go` — `SearchPosts` 인터페이스 + 메서드 추가
- `community-service/internal/repository/elasticsearch/search.go` — `SearchPosts` ES multi_match 쿼리
- `community-service/internal/repository/memory/search.go` — `SearchPosts` fallback (no-op)

### C-3: Notification Kafka 이벤트 소비자 (완료)

**수정 파일**: `backend/services/notification-service/cmd/main.go`
- Kafka EventBus 연결 (KAFKA_BROKERS 환경변수 기반)
- 3개 이벤트 구독:
  - `payment.completed` → 결제 완료 알림 (금액 포함)
  - `measurement.completed` → 측정 결과 알림
  - `health_alert.triggered` → 긴급 건강 알림 (PriorityHigh)
- ConfigWatcher로 FCM 설정 핫리로드 지원

---

## 3. 검증 결과

| 항목 | 결과 |
|------|------|
| 전체 빌드 | **35/35 PASS** |
| 전체 테스트 | **35/35 PASS** |

---

## 4. 외부 연동 완성도 변화

| 시점 | 완성도 | 근거 |
|------|--------|------|
| Phase C 이전 (추정) | 25% | MEMORY.md 초기 추정 |
| Phase C 이전 (실사) | ~65% | OAuth/Redis/Kafka/S3/Milvus/ES/PG/WebRTC 모두 구현 확인 |
| **Phase C 이후** | **~75%** | Webhook + 검색 API + Kafka Consumer 추가 |

### 남은 외부 연동 (25%)
- SMTP 이메일 전송기 (현재 Noop)
- SMS 전송기 (Twilio/NHN Cloud)
- Apple HealthKit / Google Fit 연동
- FHIR R4 HL7 연동
- 실제 PG사 라이브 키 설정 + 검증
- Firebase Cloud Messaging v1 API 마이그레이션 (현재 legacy ServerKey)

---

*자동 생성: 2026-04-25 | 만파식(萬波息) Phase C 외부 연동 통합 계획서*
