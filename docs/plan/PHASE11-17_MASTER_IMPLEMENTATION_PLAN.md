# ManPaSik 전체 세부 구현 계획서 (Phase 11~17)

> **작성일**: 2026-02-11
> **기반**: 전수 미구현 조사 (47개 항목) + 유사 시스템 기술 조사
> **참고**: 프로덕션 체크리스트, 헬스케어 플랫폼 사례, 오픈소스 베스트 프랙티스

---

## 현재 완료 현황 (Phase 1~10)

| 항목 | 수치 |
|---|---|
| 마이크로서비스 | 21개 (20 + Gateway) |
| gRPC RPC | 152개 (핸들러 구현: 147/152 = 96.7%) |
| REST API 엔드포인트 | 66개 |
| PostgreSQL Repository | 16/20 (4개 누락: community, video, translation, telemedicine) |
| Docker Compose 서비스 | 21/21 |
| Kubernetes 매니페스트 | 39개 (3환경) |
| CI/CD | 22서비스 빌드/테스트/배포 |
| 관측성 | Prometheus 메트릭 + gRPC Interceptor (21/21) |
| E2E 테스트 | 8개 파일 |
| Flutter 화면 | 8/14 (57%) |
| Flutter 백엔드 연동 | 4/21 (19%) |

---

## 미구현 전수 목록 (47개 → 7개 Phase로 분배)

### 분류 기준
- **의존성 순서**: 하위 인프라 → 미들웨어 → 서비스 → 프론트엔드 → 문서
- **리스크 우선**: 보안(P0) → 데이터무결성(P1) → 기능(P2) → 문서(P3)
- **병렬성**: 독립적인 작업은 동시 실행, 의존 작업은 순차 배치

---

## Phase 11: 핵심 인프라 통합 (Redis + Kafka + Auth + Validation)

> **목표**: 프로덕션 필수 인프라 4대 요소 통합
> **예상 규모**: 파일 30~40개 생성/수정
> **선행조건**: 없음 (현재 즉시 착수 가능)

### 11-A: Redis 통합 (Agent A)

**배경 조사 결과:**
- Redis는 JWT 토큰 저장, 세션 캐시, API 응답 캐시에 최적 (TTL 지원, 분산 검증)
- 현재 `shared/config/config.go`에 Redis 설정 존재하나 코드에서 미사용
- auth-service의 `token.go`가 인메모리, "Redis 연동은 추후" 주석 있음

**구현 항목:**

| # | 파일 | 작업 |
|---|---|---|
| 1 | `backend/shared/cache/redis.go` | Redis 클라이언트 공통 패키지 (`go-redis/v9`) |
| 2 | `backend/shared/cache/redis_test.go` | 연결, Set/Get/Del, TTL, Pipeline 테스트 |
| 3 | `backend/services/auth-service/internal/repository/redis/token.go` | JWT Refresh Token 저장소 (Redis 기반, TTL 자동 만료) |
| 4 | `backend/services/auth-service/cmd/main.go` | Redis 연결 → Redis TokenRepo fallback → Memory |
| 5 | `backend/shared/cache/session.go` | 세션 캐시 유틸 (사용자 세션, 디바이스 세션) |
| 6 | `backend/shared/cache/api_cache.go` | API 응답 캐시 미들웨어 (GET 요청 캐싱) |

**기술 스택:**
```
go get github.com/redis/go-redis/v9
```
- `redis.NewClient()` with `config.Redis.Addr()`
- TTL: RefreshToken 7d, Session 24h, API Cache 5m
- Health check: `client.Ping(ctx)`

### 11-B: Kafka/Redpanda 통합 (Agent B)

**배경 조사 결과:**
- Redpanda는 Kafka API 호환, 이미 docker-compose에 구성 (port 19092)
- `franz-go` 라이브러리가 Go 생태계에서 가장 활발 (Redpanda 공식 추천)
- 현재 EventBus는 인메모리 전용 → 서비스 재시작 시 이벤트 유실

**구현 항목:**

| # | 파일 | 작업 |
|---|---|---|
| 1 | `backend/shared/events/kafka_adapter.go` | Kafka Producer/Consumer 어댑터 (franz-go) |
| 2 | `backend/shared/events/kafka_adapter_test.go` | Producer/Consumer 유닛 테스트 |
| 3 | `backend/shared/events/eventbus.go` 수정 | `KafkaEventBus` 구현 (EventBus 인터페이스 확장) |
| 4 | `backend/shared/events/schema.go` | 이벤트 스키마 정의 (JSON Schema, 버전 관리) |
| 5 | `backend/shared/events/dead_letter.go` | Dead Letter Queue 처리 (실패 이벤트 재처리) |

**아키텍처:**
```
서비스 → KafkaProducer → Redpanda Topic → KafkaConsumer → 서비스
                                           ↘ DLQ Topic (실패 시)
```

**토픽 설계:**
```
manpasik.reservation.created     (partition: 3, retention: 7d)
manpasik.prescription.created    (partition: 3, retention: 7d)
manpasik.measurement.completed   (partition: 6, retention: 30d)
manpasik.payment.completed       (partition: 3, retention: 90d)
manpasik.notification.send       (partition: 6, retention: 3d)
manpasik.dlq                     (partition: 1, retention: 30d)
```

**기술 스택:**
```
go get github.com/twmb/franz-go
go get github.com/twmb/franz-go/pkg/kadm
```

### 11-C: Auth Interceptor 전체 서비스 적용 (Agent C)

**배경 조사 결과:**
- 현재 auth-service에만 AuthInterceptor 적용, 나머지 19개 서비스 인증 없음
- Best Practice: API Gateway에서 1차 인증, 각 서비스에서 2차 검증
- RBAC은 Role 필드는 존재하나 핸들러에서 미체크

**구현 항목:**

| # | 파일 | 작업 |
|---|---|---|
| 1 | `backend/shared/middleware/auth.go` 수정 | 서비스별 Public Method 화이트리스트 설정 가능하게 확장 |
| 2 | `backend/shared/middleware/rbac.go` | 역할 기반 접근 제어 (admin/medical_staff/user/family_member 권한 매트릭스) |
| 3 | `backend/shared/middleware/rate_limit.go` | Token Bucket 기반 Rate Limiting 인터셉터 (Redis 백엔드) |
| 4 | `backend/shared/middleware/request_id.go` | Request ID 생성/전파 (X-Request-ID 헤더, gRPC metadata) |
| 5 | 19개 서비스 `cmd/main.go` | Auth + RBAC + RateLimit + RequestID 인터셉터 체이닝 |

**RBAC 권한 매트릭스:**
```
admin:          모든 RPC 접근
medical_staff:  예약/처방/건강기록/원격진료 RPC
user:           본인 데이터 CRUD, 측정, 커뮤니티
family_member:  공유된 가족 데이터 읽기
researcher:     익명화된 데이터 읽기 (향후)
```

### 11-D: 입력 검증 유틸리티 (Agent D)

**배경 조사 결과:**
- protoc-gen-validate (PGV)가 Proto 레벨 검증 표준
- 현재 핸들러에서 `req == nil` 체크만 존재, 필드 레벨 검증 부재

**구현 항목:**

| # | 파일 | 작업 |
|---|---|---|
| 1 | `backend/shared/validation/validator.go` | 공통 검증 유틸 (이메일, 전화번호, UUID, 날짜범위, 문자열길이) |
| 2 | `backend/shared/validation/validator_test.go` | 검증 함수 테스트 (정상/비정상 케이스) |
| 3 | `backend/shared/validation/sanitizer.go` | 입력 살균 (XSS, SQL Injection 방지) |
| 4 | 핵심 핸들러 5개 적용 | auth, measurement, prescription, payment, health-record 핸들러 |

---

## Phase 12: 외부 시스템 연동 (Milvus + Elasticsearch + S3 + DB Migration)

> **목표**: 벡터 검색, 풀텍스트 검색, 파일 저장소, DB 마이그레이션
> **예상 규모**: 파일 25~30개 생성/수정
> **선행조건**: Phase 11 Redis 통합 완료

### 12-A: Milvus 벡터 DB 연동 (Agent A)

**배경 조사 결과:**
- Milvus Go SDK v2.5.x (`github.com/milvus-io/milvus/client/v2`)
- measurement-service의 `vector.go`가 인메모리 코사인 유사도로 스텁됨
- docker-compose에 Milvus 2.4.10 이미 구성 (port 19530)

**구현 항목:**

| # | 파일 | 작업 |
|---|---|---|
| 1 | `backend/shared/vectordb/milvus.go` | Milvus 클라이언트 공통 패키지 (연결, 컬렉션 관리) |
| 2 | `backend/shared/vectordb/milvus_test.go` | 연결, 삽입, 검색 테스트 |
| 3 | `backend/services/measurement-service/internal/repository/milvus/vector.go` | 실제 Milvus 벡터 저장소 (VectorRepository 구현) |
| 4 | `backend/services/measurement-service/cmd/main.go` | Milvus 연결 → MilvusVectorRepo fallback → MemoryVectorRepo |
| 5 | `backend/shared/config/config.go` 수정 | Milvus 설정 추가 (Host, Port, CollectionName) |

**컬렉션 스키마:**
```
Collection: manpasik_fingerprints
Fields:
  - id: VARCHAR (primary key)
  - session_id: VARCHAR
  - user_id: VARCHAR
  - vector: FLOAT_VECTOR (dim=1792)  // 88→448→896→1792 fingerprint
  - created_at: INT64
Index: IVF_FLAT (nlist=1024, metric=COSINE)
```

### 12-B: Elasticsearch 검색 연동 (Agent B)

**배경 조사 결과:**
- Elasticsearch 8.14.0이 docker-compose에 이미 구성 (port 9200)
- Go 클라이언트 v8 (`github.com/elastic/go-elasticsearch/v8`) Typed API 사용
- 헬스케어 데이터 인덱싱에 적합 (건강기록, 커뮤니티 글, 의약품)

**구현 항목:**

| # | 파일 | 작업 |
|---|---|---|
| 1 | `backend/shared/search/elasticsearch.go` | ES 클라이언트 공통 패키지 (연결, 인덱스 관리, 검색) |
| 2 | `backend/shared/search/elasticsearch_test.go` | 인덱스 생성, 문서 삽입, 검색 테스트 |
| 3 | `backend/services/health-record-service/internal/search/indexer.go` | 건강기록 인덱서 (CRUD → ES 동기화) |
| 4 | `backend/services/community-service/internal/search/indexer.go` | 커뮤니티 글 인덱서 |
| 5 | `backend/services/prescription-service/internal/search/indexer.go` | 의약품 검색 인덱서 |
| 6 | `backend/shared/config/config.go` 수정 | Elasticsearch 설정 추가 |

**인덱스 설계:**
```
manpasik_health_records:  record_id, user_id, diagnosis, symptoms, notes (Korean analyzer)
manpasik_community_posts: post_id, title, content, tags, author (Korean + English)
manpasik_medications:     medication_name, active_ingredient, interactions, dosage
```

### 12-C: S3/MinIO 파일 저장소 (Agent C)

**배경 조사 결과:**
- MinIO가 docker-compose에 이미 Milvus 스토리지로 존재
- 추가 버킷 생성으로 앱 파일 저장 가능
- 프로필 이미지, 측정 보고서 PDF, 의료 영상 등 필요

**구현 항목:**

| # | 파일 | 작업 |
|---|---|---|
| 1 | `backend/shared/storage/s3.go` | S3 호환 스토리지 클라이언트 (MinIO SDK) |
| 2 | `backend/shared/storage/s3_test.go` | 업로드, 다운로드, 삭제, Pre-signed URL 테스트 |
| 3 | `backend/gateway/internal/router/upload_handlers.go` | 파일 업로드 REST 엔드포인트 (multipart/form-data) |
| 4 | `backend/shared/config/config.go` 수정 | S3 설정 추가 (Endpoint, AccessKey, SecretKey, Bucket) |

**버킷 설계:**
```
manpasik-profiles:   프로필 이미지 (max 5MB, image/*)
manpasik-reports:    측정 보고서 PDF (max 20MB, application/pdf)
manpasik-medical:    의료 데이터 (HIPAA 암호화 필수)
manpasik-community:  커뮤니티 첨부파일 (max 10MB)
```

### 12-D: golang-migrate DB 마이그레이션 (Agent D)

**배경 조사 결과:**
- `golang-migrate/migrate` v4.19.1 (Go 1.24 지원, 18K+ Stars)
- 현재 init SQL만 존재, 스키마 변경 관리 불가
- 프로덕션에서 expand-contract 마이그레이션 필수

**구현 항목:**

| # | 파일 | 작업 |
|---|---|---|
| 1 | `backend/migrations/000001_initial_schema.up.sql` | 전체 초기 스키마 (init SQL 통합) |
| 2 | `backend/migrations/000001_initial_schema.down.sql` | 롤백 스크립트 |
| 3 | `backend/migrations/000002_add_indexes.up.sql` | 성능 인덱스 추가 |
| 4 | `backend/cmd/migrate/main.go` | 마이그레이션 CLI 도구 |
| 5 | `Makefile` 수정 | `make migrate-up`, `make migrate-down`, `make migrate-create` |

**기술 스택:**
```
go get -u github.com/golang-migrate/migrate/v4
go get -u github.com/golang-migrate/migrate/v4/database/pgx5
go get -u github.com/golang-migrate/migrate/v4/source/file
```

---

## Phase 13: 스트리밍 + 실시간 통신 (gRPC Stream + WebRTC + WebSocket)

> **목표**: 실시간 측정 데이터, 디바이스 스트리밍, 영상 통화
> **예상 규모**: 파일 15~20개 생성/수정
> **선행조건**: Phase 11 완료 (Redis 세션 관리)

### 13-A: gRPC 양방향 스트리밍 (Agent A)

**배경 조사 결과:**
- `StreamMeasurement` (measurement-service): 디바이스 → 서버 실시간 데이터
- `StreamDeviceStatus` (device-service): 디바이스 상태 실시간 모니터링
- Go 패턴: 별도 goroutine으로 Send/Recv 분리, `io.EOF`로 종료 감지

**구현 항목:**

| # | 파일 | 작업 |
|---|---|---|
| 1 | `measurement-service/internal/handler/stream.go` | StreamMeasurement 양방향 스트림 핸들러 |
| 2 | `measurement-service/internal/service/stream.go` | 실시간 측정 데이터 처리 로직 (버퍼링, 이상치 감지) |
| 3 | `device-service/internal/handler/stream.go` | StreamDeviceStatus 양방향 스트림 핸들러 |
| 4 | `device-service/internal/handler/grpc.go` 수정 | RequestOtaUpdate 핸들러 추가 |
| 5 | `device-service/internal/service/stream.go` | 디바이스 상태 스트림 로직 |

**스트리밍 아키텍처:**
```
클라이언트 → StreamMeasurement(stream MeasurementData)
서버: goroutine1 → Recv() 측정값
      goroutine2 → Send() 실시간 피드백 (정상/이상/캘리브레이션 필요)
      channel로 goroutine 간 통신
```

### 13-B: WebRTC 영상 통화 (Agent B)

**배경 조사 결과:**
- Pion (`github.com/pion/webrtc/v4`)이 Go WebRTC 표준 라이브러리
- NeoCareSphere (HIPAA 원격진료 플랫폼)이 유사 아키텍처 참고
- 현재 video-service는 Signal 저장만 구현, 실제 P2P 연결 없음

**구현 항목:**

| # | 파일 | 작업 |
|---|---|---|
| 1 | `backend/services/video-service/internal/signaling/server.go` | WebRTC 시그널링 서버 (Offer/Answer/ICE 교환) |
| 2 | `backend/services/video-service/internal/signaling/room_manager.go` | 방 별 피어 관리, SDP 교환 중재 |
| 3 | `backend/services/video-service/internal/signaling/turn_config.go` | TURN/STUN 서버 설정 (coturn 연동) |
| 4 | `infrastructure/docker/docker-compose.dev.yml` 수정 | coturn 컨테이너 추가 (TURN 서버) |
| 5 | `backend/services/video-service/cmd/main.go` 수정 | WebSocket 시그널링 엔드포인트 추가 |

**시그널링 흐름:**
```
환자 → WebSocket → SignalingServer → WebSocket → 의사
  1. CreateRoom → room_id
  2. JoinRoom(room_id) → WebSocket 연결
  3. Offer SDP → Server relay → Answer SDP
  4. ICE Candidates 교환
  5. P2P 미디어 스트림 수립
  6. EndRoom → 정리
```

### 13-C: 누락 핸들러 + Repository 구현 (Agent C)

**구현 항목:**

| # | 서비스 | 작업 |
|---|---|---|
| 1 | subscription-service | `CheckCartridgeAccess`, `ListAccessibleCartridges` 핸들러 + 서비스 메서드 |
| 2 | community-service | PostgreSQL Repository 구현 |
| 3 | video-service | PostgreSQL Repository 구현 |
| 4 | translation-service | PostgreSQL Repository 구현 |
| 5 | telemedicine-service | PostgreSQL Repository 구현 |

**PostgreSQL Repository 패턴:** 기존 16개 서비스와 동일한 `pgxpool` 패턴 적용

### 13-D: 실시간 알림 WebSocket (Agent D)

**구현 항목:**

| # | 파일 | 작업 |
|---|---|---|
| 1 | `backend/gateway/internal/ws/hub.go` | WebSocket Hub (연결 관리, 방 구독) |
| 2 | `backend/gateway/internal/ws/client.go` | WebSocket 클라이언트 핸들러 |
| 3 | `backend/gateway/internal/router/ws_handlers.go` | `GET /ws` WebSocket 업그레이드 엔드포인트 |
| 4 | `backend/gateway/cmd/main.go` 수정 | WebSocket 라우트 등록 |

**기술 스택:**
```
go get github.com/gorilla/websocket
```

**메시지 형식:**
```json
{
  "type": "notification",
  "event": "measurement.completed",
  "data": { "session_id": "...", "status": "normal" },
  "timestamp": "2026-02-11T15:00:00Z"
}
```

---

## Phase 14: 프론트엔드 확장 (Flutter 6개 화면 + 백엔드 연동)

> **목표**: 누락 6개 주요 화면 구현, REST 클라이언트 통합
> **예상 규모**: 파일 40~50개 생성
> **선행조건**: Phase 11~13 백엔드 완료

### 14-A: Data Hub + AI Coach 화면 (Agent A)

**구현 항목:**

| # | 파일 | 작업 |
|---|---|---|
| 1 | `lib/features/data/data_hub_screen.dart` | 건강 타임라인, 트렌드 차트 (fl_chart) |
| 2 | `lib/features/data/widgets/health_timeline.dart` | 시간순 건강 데이터 위젯 |
| 3 | `lib/features/data/widgets/trend_chart.dart` | 측정값 트렌드 차트 |
| 4 | `lib/features/data/data_provider.dart` | 데이터 허브 Riverpod Provider |
| 5 | `lib/features/coach/coach_screen.dart` | AI 코칭 화면 (대화형 UI) |
| 6 | `lib/features/coach/widgets/coaching_card.dart` | 코칭 메시지 카드 |
| 7 | `lib/features/coach/widgets/goal_tracker.dart` | 목표 진행률 위젯 |
| 8 | `lib/features/coach/coach_provider.dart` | AI Coach Riverpod Provider |

### 14-B: Market + Community 화면 (Agent B)

| # | 파일 | 작업 |
|---|---|---|
| 1 | `lib/features/market/market_screen.dart` | 카트리지 쇼핑 메인 |
| 2 | `lib/features/market/product_detail_screen.dart` | 상품 상세 |
| 3 | `lib/features/market/cart_screen.dart` | 장바구니 |
| 4 | `lib/features/market/subscription_screen.dart` | 구독 관리 |
| 5 | `lib/features/market/market_provider.dart` | 마켓 Provider |
| 6 | `lib/features/community/community_screen.dart` | 커뮤니티 메인 |
| 7 | `lib/features/community/post_detail_screen.dart` | 글 상세 |
| 8 | `lib/features/community/challenge_screen.dart` | 챌린지 |
| 9 | `lib/features/community/community_provider.dart` | 커뮤니티 Provider |

### 14-C: Medical + Family 화면 (Agent C)

| # | 파일 | 작업 |
|---|---|---|
| 1 | `lib/features/medical/medical_screen.dart` | 의료 서비스 메인 |
| 2 | `lib/features/medical/reservation_screen.dart` | 병원 예약 |
| 3 | `lib/features/medical/prescription_screen.dart` | 처방전 관리 |
| 4 | `lib/features/medical/telemedicine_screen.dart` | 원격진료 (WebRTC 영상통화) |
| 5 | `lib/features/medical/medical_provider.dart` | 의료 서비스 Provider |
| 6 | `lib/features/family/family_screen.dart` | 가족 관리 메인 |
| 7 | `lib/features/family/member_detail_screen.dart` | 가족 구성원 상세 |
| 8 | `lib/features/family/family_provider.dart` | 가족 Provider |

### 14-D: 공통 인프라 + 라우팅 + 모델 (Agent D)

| # | 파일 | 작업 |
|---|---|---|
| 1 | `lib/shared/models/user.dart` | 사용자 도메인 모델 |
| 2 | `lib/shared/models/measurement.dart` | 측정 도메인 모델 |
| 3 | `lib/shared/models/health_record.dart` | 건강기록 도메인 모델 |
| 4 | `lib/shared/models/product.dart` | 상품/카트리지 도메인 모델 |
| 5 | `lib/shared/models/prescription.dart` | 처방전 도메인 모델 |
| 6 | `lib/shared/models/notification.dart` | 알림 도메인 모델 |
| 7 | `lib/core/router/app_router.dart` 수정 | 6개 신규 화면 라우트 추가 |
| 8 | `lib/features/home/home_screen.dart` 수정 | 네비게이션 바에 6개 화면 연결 |
| 9 | `lib/shared/widgets/notification_badge.dart` | 알림 배지 위젯 |
| 10 | `lib/core/services/rest_client.dart` → 각 Provider에서 사용 연결 | REST 클라이언트 실제 통합 |

---

## Phase 15: 모니터링 + 보안 강화

> **목표**: Grafana 대시보드, 알림 규칙, 분산 추적, TLS
> **예상 규모**: 파일 15~20개 생성
> **선행조건**: Phase 11 (Redis, Rate Limiting)

### 15-A: Grafana 대시보드 + Prometheus 알림 (Agent A)

| # | 파일 | 작업 |
|---|---|---|
| 1 | `infrastructure/docker/config/prometheus/prometheus.yml` 수정 | 전체 21개 서비스 scrape config |
| 2 | `infrastructure/docker/config/prometheus/alerts.yml` | 알림 규칙 (에러율 >1%, P95 >500ms, 서비스 다운) |
| 3 | `infrastructure/docker/config/grafana/dashboards/service-overview.json` | 서비스 개요 대시보드 |
| 4 | `infrastructure/docker/config/grafana/dashboards/business-metrics.json` | 비즈니스 메트릭 대시보드 (측정 수, 활성 사용자) |
| 5 | `infrastructure/docker/config/grafana/provisioning/dashboards.yml` | 대시보드 자동 프로비저닝 |
| 6 | `infrastructure/docker/config/grafana/provisioning/datasources.yml` | Prometheus 데이터소스 설정 |

### 15-B: OpenTelemetry 분산 추적 (Agent B)

**배경 조사 결과:**
- go.sum에 이미 opentelemetry 의존성 존재
- Jaeger/Tempo로 트레이스 수집

| # | 파일 | 작업 |
|---|---|---|
| 1 | `backend/shared/observability/tracing.go` | OpenTelemetry TracerProvider 초기화 |
| 2 | `backend/shared/observability/grpc_interceptor.go` 수정 | 트레이싱 인터셉터 추가 |
| 3 | `infrastructure/docker/docker-compose.dev.yml` 수정 | Jaeger 컨테이너 추가 |
| 4 | 21개 cmd/main.go | TracerProvider 초기화 코드 추가 |

### 15-C: K8s 보안 강화 (Agent C)

| # | 파일 | 작업 |
|---|---|---|
| 1 | `infrastructure/kubernetes/base/networkpolicies/default-deny.yaml` | 기본 거부 정책 |
| 2 | `infrastructure/kubernetes/base/networkpolicies/allow-gateway.yaml` | Gateway → 서비스 허용 |
| 3 | `infrastructure/kubernetes/base/networkpolicies/allow-data.yaml` | 서비스 → DB/Redis/Kafka 허용 |
| 4 | `infrastructure/kubernetes/base/data/postgres-statefulset.yaml` | PostgreSQL StatefulSet |
| 5 | `infrastructure/kubernetes/base/data/redis-deployment.yaml` | Redis Deployment |
| 6 | `infrastructure/kubernetes/base/data/redpanda-statefulset.yaml` | Redpanda StatefulSet |

### 15-D: 프로덕션 Docker Compose (Agent D)

| # | 파일 | 작업 |
|---|---|---|
| 1 | `infrastructure/docker/docker-compose.yml` | 프로덕션용 (TLS, 리소스 제한, 로그 드라이버) |
| 2 | `infrastructure/docker/docker-compose.monitoring.yml` | 모니터링 스택 (Prometheus + Grafana + Jaeger) |
| 3 | `infrastructure/docker/.env.example` | 환경변수 템플릿 |

---

## Phase 16: 문서화 + API 문서

> **목표**: OpenAPI 문서, 아키텍처 다이어그램, 개발자 가이드
> **예상 규모**: 파일 10~15개 생성
> **선행조건**: Phase 11~15 주요 기능 완료

### 16-A: OpenAPI/Swagger API 문서 (Agent A)

| # | 파일 | 작업 |
|---|---|---|
| 1 | `docs/api/openapi.yaml` | 전체 66+ REST API 엔드포인트 OpenAPI 3.0 스펙 |
| 2 | `docs/api/schemas/` | 요청/응답 스키마 정의 |
| 3 | `backend/gateway/cmd/main.go` 수정 | `/docs` Swagger UI 서빙 |

### 16-B: 아키텍처 문서 (Agent B)

| # | 파일 | 작업 |
|---|---|---|
| 1 | `docs/architecture/system-overview.md` | 전체 시스템 아키텍처 (Mermaid 다이어그램) |
| 2 | `docs/architecture/data-flow.md` | 주요 데이터 흐름 시퀀스 다이어그램 |
| 3 | `docs/architecture/deployment.md` | 배포 아키텍처 (K8s 토폴로지) |
| 4 | `docs/architecture/security.md` | 보안 아키텍처 (인증/인가 흐름) |

### 16-C: 개발자 가이드 (Agent C)

| # | 파일 | 작업 |
|---|---|---|
| 1 | `docs/guides/getting-started.md` | 개발 환경 설정 가이드 |
| 2 | `docs/guides/adding-new-service.md` | 새 마이크로서비스 추가 가이드 |
| 3 | `docs/guides/testing.md` | 테스트 전략 가이드 |
| 4 | `docs/guides/deployment.md` | 배포 가이드 (dev/staging/prod) |

### 16-D: Rust FFI + Flutter 통합 가이드 (Agent D)

| # | 파일 | 작업 |
|---|---|---|
| 1 | `frontend/flutter-app/pubspec.yaml` 수정 | flutter_rust_bridge 활성화 |
| 2 | `frontend/flutter-app/lib/core/services/rust_bridge.dart` | Rust FFI 브릿지 실제 연결 |
| 3 | `docs/guides/rust-ffi-integration.md` | Rust Core ↔ Flutter 통합 가이드 |

---

## Phase 17: Phase 4 신규 서비스 (생태계 확장)

> **목표**: 로드맵 Phase 4 서비스 중 우선순위 3개 구현
> **예상 규모**: 파일 30~40개 생성
> **선행조건**: Phase 11~16 완료

### 우선순위 서비스 (기술 조사 기반)

| 순위 | 서비스 | 사유 |
|---|---|---|
| 1 | **analytics-service** | 비즈니스 인텔리전스, 기존 데이터 활용 가능 |
| 2 | **vision-service** | 음식 사진 분석, TFLite 엣지 AI 활용 |
| 3 | **iot-gateway-service** | MQTT 브로커, 디바이스 스케일 지원 |

### 17-A: analytics-service (Agent A)

```
Proto: AnalyticsService
RPCs: GetUserAnalytics, GetDeviceAnalytics, GetBusinessMetrics, 
      GenerateReport, GetTrend, ExportAnalytics
Stack: TimescaleDB continuous aggregates + Grafana 임베딩
```

### 17-B: vision-service (Agent B)

```
Proto: VisionService
RPCs: AnalyzeFood, GetNutritionInfo, ClassifyImage, GetCalories
Stack: TFLite Go binding + Rust Core vision module
```

### 17-C: iot-gateway-service (Agent C)

```
Proto: IoTGatewayService
RPCs: RegisterDevice, StreamTelemetry, SendCommand, GetDeviceStatus
Stack: Eclipse Paho MQTT + gRPC bridge
```

---

## 전체 일정 요약

| Phase | 내용 | 에이전트 | 파일 수 | 의존성 |
|---|---|---|---|---|
| **11** | Redis + Kafka + Auth + Validation | A/B/C/D 병렬 | 30~40 | 없음 |
| **12** | Milvus + ES + S3 + DB Migration | A/B/C/D 병렬 | 25~30 | Phase 11 |
| **13** | gRPC Stream + WebRTC + 누락 Repo + WebSocket | A/B/C/D 병렬 | 15~20 | Phase 11 |
| **14** | Flutter 6화면 + 모델 + 라우팅 | A/B/C/D 병렬 | 40~50 | Phase 13 |
| **15** | Grafana + OTel + K8s 보안 + Prod Compose | A/B/C/D 병렬 | 15~20 | Phase 11 |
| **16** | OpenAPI + 아키텍처 + 가이드 + Rust FFI | A/B/C/D 병렬 | 10~15 | Phase 15 |
| **17** | analytics + vision + iot-gateway 신규 | A/B/C 병렬 | 30~40 | Phase 16 |

**총 예상 파일**: 165~215개 생성/수정
**총 미구현 해소**: 47/47 (100%)

---

## 기술 스택 추가 목록

| 라이브러리 | 버전 | 용도 |
|---|---|---|
| `github.com/redis/go-redis/v9` | latest | Redis 클라이언트 |
| `github.com/twmb/franz-go` | latest | Kafka/Redpanda 클라이언트 |
| `github.com/milvus-io/milvus/client/v2` | v2.5.x | Milvus 벡터 DB |
| `github.com/elastic/go-elasticsearch/v8` | v8.19+ | Elasticsearch 검색 |
| `github.com/minio/minio-go/v7` | latest | S3/MinIO 파일 저장 |
| `github.com/golang-migrate/migrate/v4` | v4.19+ | DB 마이그레이션 |
| `github.com/pion/webrtc/v4` | latest | WebRTC P2P 영상 |
| `github.com/gorilla/websocket` | latest | WebSocket 실시간 알림 |
| `go.opentelemetry.io/otel` | latest | 분산 추적 |
| `flutter_rust_bridge` | ^2.0.0 | Dart ↔ Rust FFI |

---

## 리스크 관리

| 리스크 | 확률 | 영향 | 완화 전략 |
|---|---|---|---|
| Milvus SDK 호환성 | 중 | 높 | 인메모리 fallback 유지 |
| WebRTC NAT 트래버설 실패 | 높 | 중 | TURN 서버 (coturn) 필수 배포 |
| Kafka 메시지 유실 | 저 | 높 | ack=all, DLQ, 멱등 컨슈머 |
| Flutter Rust FFI 빌드 오류 | 중 | 중 | 스텁 유지, 점진적 활성화 |
| DB 마이그레이션 충돌 | 저 | 높 | 트랜잭션 래핑, 테스트 환경 선행 |

---

> **다음 단계**: 사용자 승인 후 Phase 11부터 4개 에이전트 병렬 실행
