# ManPaSik 서비스 간 통신 패턴 명세서 (Service Communication Patterns)

**문서번호**: MPK-COMM-PATTERNS-v1.0  
**작성일**: 2026-02-14  
**목적**: 29개 MSA 서비스 간 동기/비동기 통신 패턴, 장애 대응, 분산 트랜잭션 전략을 상세 정의  
**적용**: 전체 Go 백엔드 서비스, Rust 코어, Flutter 클라이언트

---

## 1. 통신 패턴 개요

### 1.1 통신 방식 분류

| 방식 | 프로토콜 | 용도 | 특성 |
|------|---------|------|------|
| **동기 (Request-Reply)** | gRPC (HTTP/2, Protobuf) | 즉시 응답 필요 작업 | 낮은 지연, 강한 결합 |
| **비동기 (Event-Driven)** | Kafka (Redpanda) | 이벤트 전파, 최종 일관성 | 느슨한 결합, 높은 확장성 |
| **스트리밍** | gRPC Bidirectional Stream | 실시간 측정 데이터 | 양방향, 저지연 |
| **외부 연동** | REST (Kong Gateway) | 클라이언트↔백엔드, 외부 API | 범용, 인증 통합 |

### 1.2 통신 원칙

1. **Command는 동기, Event는 비동기**: 사용자 요청(CRUD, 조회)은 gRPC, 부수 효과(알림, 분석, 로그)는 Kafka
2. **서비스 자율성**: 각 서비스는 자체 DB 소유, 다른 서비스 DB 직접 접근 금지
3. **Idempotency**: 모든 gRPC 핸들러와 Kafka 컨슈머는 멱등성 보장 (event_id 기반 중복 제거)
4. **Graceful Degradation**: 비핵심 서비스 장애 시 핵심 기능(측정, 인증) 계속 동작

---

## 2. 동기 통신 (gRPC) 상세

### 2.1 서비스 호출 관계 매트릭스

```
호출자(→) ↓ 피호출자
```

| 호출자 \ 피호출자 | auth | user | device | meas | sub | shop | pay | ai | cart | calib | coach | notif | family | hr | admin |
|------------------|------|------|--------|------|-----|------|-----|-----|------|-------|-------|-------|--------|-----|-------|
| **gateway** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **auth** | — | ✅ | — | — | — | — | — | — | — | — | — | — | — | — | — |
| **measurement** | — | — | ✅ | — | — | — | — | — | ✅ | ✅ | — | — | — | — | — |
| **ai-inference** | — | ✅ | — | ✅ | — | — | — | — | ✅ | — | ✅ | ✅ | — | — | — |
| **coaching** | — | ✅ | — | ✅ | ✅ | — | — | ✅ | — | — | — | ✅ | — | ✅ | — |
| **payment** | — | — | — | — | ✅ | ✅ | — | — | — | — | — | ✅ | — | — | — |
| **subscription** | — | ✅ | — | — | — | — | ✅ | — | ✅ | — | — | ✅ | — | — | — |
| **shop** | — | ✅ | — | — | ✅ | — | ✅ | — | ✅ | — | — | — | — | — | — |
| **family** | — | ✅ | — | ✅ | — | — | — | — | — | — | — | ✅ | — | ✅ | — |
| **admin** | — | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | — |
| **notification** | — | ✅ | ✅ | — | — | — | — | — | — | — | — | — | ✅ | — | — |

### 2.2 gRPC 서비스 디스커버리

```text
서비스 주소 규칙 (Kubernetes DNS):
  <service-name>.<namespace>.svc.cluster.local:<port>

예시:
  auth-service.manpasik.svc.cluster.local:50051
  user-service.manpasik.svc.cluster.local:50052

Docker Compose (개발 환경):
  auth-service:50051
  user-service:50052
```

| 서비스 | DNS 이름 | 포트 | 프로토콜 |
|--------|---------|------|---------|
| auth-service | auth-service | 50051 | gRPC (TLS) |
| user-service | user-service | 50052 | gRPC (TLS) |
| device-service | device-service | 50053 | gRPC (TLS) |
| measurement-service | measurement-service | 50054 | gRPC (TLS) |
| subscription-service | subscription-service | 50055 | gRPC (TLS) |
| shop-service | shop-service | 50056 | gRPC (TLS) |
| payment-service | payment-service | 50057 | gRPC (TLS) |
| ai-inference-service | ai-inference-service | 50058 | gRPC (TLS) |
| cartridge-service | cartridge-service | 50059 | gRPC (TLS) |
| calibration-service | calibration-service | 50060 | gRPC (TLS) |
| coaching-service | coaching-service | 50061 | gRPC (TLS) |
| notification-service | notification-service | 50062 | gRPC (TLS) |
| family-service | family-service | 50063 | gRPC (TLS) |
| health-record-service | health-record-service | 50064 | gRPC (TLS) |
| prescription-service | prescription-service | 50065 | gRPC (TLS) |
| reservation-service | reservation-service | 50066 | gRPC (TLS) |
| community-service | community-service | 50067 | gRPC (TLS) |
| translation-service | translation-service | 50068 | gRPC (TLS) |
| video-service | video-service | 50069 | gRPC (TLS) |
| admin-service | admin-service | 50070 | gRPC (TLS) |
| telemedicine-service | telemedicine-service | 50071 | gRPC (TLS) |

### 2.3 gRPC 연결 관리

```go
// 공통 gRPC 클라이언트 팩토리 패턴
type ServiceClients struct {
    AuthClient    pb.AuthServiceClient
    UserClient    pb.UserServiceClient
    DeviceClient  pb.DeviceServiceClient
    // ...
}

func NewServiceClients(cfg *config.Config) (*ServiceClients, error) {
    opts := []grpc.DialOption{
        grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
        grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
            grpc_retry.UnaryClientInterceptor(retryOpts...),
            otelgrpc.UnaryClientInterceptor(),
            grpc_auth.UnaryClientInterceptor(authFunc),
        )),
        grpc.WithKeepaliveParams(keepalive.ClientParameters{
            Time:                30 * time.Second,
            Timeout:             10 * time.Second,
            PermitWithoutStream: true,
        }),
    }
    // ... 각 서비스 연결 초기화
}
```

**연결 풀링 설정:**

| 파라미터 | 값 | 근거 |
|---------|-----|------|
| `MaxConcurrentStreams` | 100 | 서비스별 동시 스트림 수 |
| `InitialWindowSize` | 1MB | 흐름 제어 초기 윈도우 |
| `Keepalive.Time` | 30s | 유휴 연결 확인 주기 |
| `Keepalive.Timeout` | 10s | Keepalive 응답 대기 |
| `MaxRecvMsgSize` | 10MB | 측정 데이터 스트림 대응 |

---

## 3. 비동기 통신 (Kafka) 상세

### 3.1 이벤트 발행-구독 매핑

| 토픽 | 생산자 | 소비자 | 패턴 |
|------|--------|--------|------|
| `manpasik.measurement.completed` | measurement-service | ai-inference, coaching, notification, health-record | Fan-out |
| `manpasik.user.registered` | auth-service | notification, user, coaching | Fan-out |
| `manpasik.payment.completed` | payment-service | subscription, shop, notification | Fan-out |
| `manpasik.payment.failed` | payment-service | notification, admin | Fan-out |
| `manpasik.subscription.changed` | subscription-service | cartridge, notification, user | Fan-out |
| `manpasik.cartridge.verified` | cartridge-service | measurement, notification | Fan-out |
| `manpasik.cartridge.depleted` | cartridge-service | notification, shop | Fan-out |
| `manpasik.notification.send` | 여러 서비스 | notification-service | Many-to-One |
| `manpasik.ai.analysis_completed` | ai-inference-service | coaching, notification, health-record | Fan-out |
| `manpasik.ai.risk.detected` | ai-inference-service | notification, family, admin | Fan-out (긴급) |
| `manpasik.device.status_changed` | device-service | notification, admin | Fan-out |
| `manpasik.coaching.generated` | coaching-service | notification | Point-to-Point |
| `manpasik.community.post.created` | community-service | notification, translation | Fan-out |
| `manpasik.reservation.created` | reservation-service | notification, telemedicine | Fan-out |
| `manpasik.prescription.created` | prescription-service | notification, health-record | Fan-out |
| `manpasik.config.changed` | admin-service | 전체 서비스 (ConfigWatcher) | Broadcast |
| `manpasik.dlq` | 전체 (실패 시) | admin (수동 재처리) | DLQ |

### 3.2 Consumer Group 전략

| Consumer Group | 서비스 | 구독 토픽 | 인스턴스 수 |
|---------------|--------|----------|-----------|
| `cg-ai-inference` | ai-inference-service | measurement.completed | 2~10 |
| `cg-coaching` | coaching-service | measurement.completed, ai.analysis_completed | 2~5 |
| `cg-notification` | notification-service | notification.send, *.risk.detected, *.completed | 3~10 |
| `cg-health-record` | health-record-service | measurement.completed, prescription.created | 2~5 |
| `cg-admin-monitor` | admin-service | config.changed, dlq, *.failed | 1~3 |

### 3.3 메시지 전달 보장

| 항목 | 설정 | 근거 |
|------|------|------|
| **Producer `acks`** | `all` (-1) | 모든 ISR 복제 확인 후 응답 |
| **`min.insync.replicas`** | 2 | 최소 2개 복제본 확인 |
| **Replication Factor** | 3 | 3-way 복제 |
| **Consumer Offset** | `enable.auto.commit=false` | 수동 커밋 (처리 완료 후) |
| **Idempotence** | `enable.idempotence=true` | 중복 메시지 방지 |
| **Retry** | `retries=3`, 지수 백오프 | 일시적 장애 대응 |
| **DLQ** | 3회 실패 → `manpasik.dlq` | 영구 실패 격리 |

---

## 4. 스트리밍 통신

### 4.1 측정 데이터 스트림 (gRPC Bidirectional)

```text
Flutter App ←→ measurement-service (StreamMeasurement RPC)

클라이언트 → 서버:
  MeasurementStreamRequest {
    session_id, device_id, raw_packet (bytes), sequence_no, timestamp
  }
  전송 주기: 10~50ms (카트리지 유형별)

서버 → 클라이언트:
  MeasurementStreamResponse {
    processed_value, confidence, progress_percent, status, alerts[]
  }
  응답 주기: 100~500ms (배치 처리 후)

세션 생명주기:
  StartSession → Stream (N packets) → EndSession
  타임아웃: 300초 무신호 시 자동 종료
  재연결: session_id + last_sequence_no로 이어받기
```

### 4.2 실시간 알림 스트림 (Server-Sent Events via Gateway)

```text
Flutter App → Kong Gateway → notification-service

SSE 연결:
  GET /api/v1/notifications/stream
  Headers: Authorization: Bearer <JWT>

이벤트 유형:
  - notification.new: 새 알림 도착
  - measurement.progress: 측정 진행률 (다른 기기에서 측정 중)
  - device.status: 기기 상태 변경
  - risk.alert: 위험 경보 (긴급)
```

---

## 5. 장애 대응 패턴

### 5.1 Circuit Breaker 설정 (서비스별)

| 대상 서비스 | 실패 임계값 | 오픈 타임아웃 | Half-Open 시도 | Fallback |
|-----------|-----------|-------------|---------------|----------|
| auth-service | 5회/10초 | 30초 | 3회 | 캐시된 토큰 검증 (Redis) |
| ai-inference-service | 3회/10초 | 60초 | 2회 | 기본 규칙 기반 분석 |
| notification-service | 10회/30초 | 60초 | 5회 | Kafka 큐에 적재 (지연 발송) |
| payment-service | 2회/10초 | 120초 | 1회 | 에러 반환 + 재시도 안내 |
| translation-service | 5회/10초 | 30초 | 3회 | 원문 표시 |
| video-service | 3회/10초 | 60초 | 2회 | "서비스 준비 중" 표시 |

### 5.2 Timeout 정책

| 호출 유형 | Deadline | 근거 |
|----------|---------|------|
| 일반 gRPC Unary | 5초 | P99 < 1초, 5배 마진 |
| AI 추론 | 10초 | GPU 연산 포함 |
| 파일 업로드 (MinIO) | 30초 | 대용량 이미지 |
| DB 쿼리 | 3초 | 인덱스 최적화 기준 |
| 외부 API (PG, FCM) | 10초 | 네트워크 변동 대응 |
| 측정 스트림 (패킷당) | 1초 | 실시간 요구 |

### 5.3 Retry 정책

```go
retryOpts := []grpc_retry.CallOption{
    grpc_retry.WithBackoff(grpc_retry.BackoffExponentialWithJitter(
        100*time.Millisecond,  // 초기 대기
        0.2,                    // 지터 비율
    )),
    grpc_retry.WithMax(3),                        // 최대 3회 재시도
    grpc_retry.WithCodes(                         // 재시도 대상 코드
        codes.Unavailable,     // 서비스 불가
        codes.ResourceExhausted, // 자원 소진
        codes.DeadlineExceeded,  // 타임아웃 (멱등 요청만)
    ),
}
```

**재시도 불가 RPC** (비멱등):
- `CreatePayment`, `PlaceOrder`, `CreateReservation` — 중복 생성 위험
- 이러한 RPC는 클라이언트 측 재시도 키(Idempotency-Key 헤더)로 보호

---

## 6. Saga 패턴 — 분산 트랜잭션

### 6.1 구독 결제 Saga

```text
[결제 Saga — Choreography 기반]

1. 사용자 → shop-service: CreateOrder(subscription_plan)
2. shop-service → payment-service: CreatePayment(order_id, amount)
3. payment-service → Toss PG: 결제 요청
4. Toss PG → payment-service: 결제 성공
5. payment-service ─→ Kafka: payment.completed
6. subscription-service ← Kafka: payment.completed
   └── SubscriptionService.UpgradeTier(user_id, new_tier)
7. cartridge-service ← Kafka: subscription.changed
   └── CartridgeService.UpdateAccessPolicy(user_id, new_tier)
8. notification-service ← Kafka: subscription.changed
   └── NotificationService.SendUpgradeConfirmation(user_id)

보상 트랜잭션 (결제 실패 시):
4'. Toss PG → payment-service: 결제 실패
5'. payment-service ─→ Kafka: payment.failed
6'. shop-service ← Kafka: payment.failed
    └── ShopService.CancelOrder(order_id)
7'. notification-service ← Kafka: payment.failed
    └── NotificationService.SendPaymentFailure(user_id)
```

### 6.2 측정→AI→코칭 Saga

```text
[측정 완료 Saga — Choreography 기반]

1. measurement-service: EndSession(session_id)
   ├── 측정 데이터 PostgreSQL/TimescaleDB 저장
   ├── 핑거프린트 Milvus 저장
   └── ─→ Kafka: measurement.completed

2. ai-inference-service ← Kafka: measurement.completed
   ├── AnalyzeMeasurement(fingerprint, model_type)
   ├── 결과 PostgreSQL 저장
   └── ─→ Kafka: ai.analysis_completed
       └── (risk_level >= high) ─→ Kafka: ai.risk.detected

3. coaching-service ← Kafka: ai.analysis_completed
   ├── GenerateCoaching(user_id, analysis_result)
   └── ─→ Kafka: coaching.generated

4. health-record-service ← Kafka: measurement.completed + ai.analysis_completed
   └── CreateHealthRecord(user_id, session_data, analysis_data)

5. notification-service ← Kafka: coaching.generated + ai.risk.detected
   └── SendPush/Email(user_id, coaching_summary | risk_alert)
```

### 6.3 예약→화상진료→처방 Saga

```text
[의료 Saga — Orchestration 기반 (Phase 3)]

Orchestrator: telemedicine-service

1. reservation-service.CreateReservation(doctor_id, patient_id, time)
   └── ─→ Kafka: reservation.created

2. (예약 시간 도래) telemedicine-service.InitiateSession(reservation_id)
   ├── video-service.CreateRoom(participants: [doctor_id, patient_id])
   └── notification-service: "화상 진료 시작" 알림

3. (진료 완료) telemedicine-service.CompleteSession(session_id)
   ├── video-service.EndRoom(room_id)
   ├── health-record.CreateRecordFromClinician(doctor_input)
   └── prescription-service.CreatePrescription(medications, diagnosis_code)
       └── ─→ Kafka: prescription.created

4. notification-service ← Kafka: prescription.created
   └── SendPush(patient_id, "처방전이 발행되었습니다")
```

---

## 7. API Versioning 전략

### 7.1 Proto 버전 관리

```protobuf
// 패키지 네이밍 규칙
package manpasik.v1;  // 현재 활성 버전

// 신규 버전 도입 시 (Breaking Change)
package manpasik.v2;  // v1과 병렬 운영

// 비-파괴적 변경은 기존 버전에서 확장
// - 새 필드 추가 (optional)
// - 새 RPC 추가
// - 새 Enum 값 추가
```

### 7.2 버전 호환성 정책

| 변경 유형 | 처리 방식 | 예시 |
|----------|----------|------|
| 필드 추가 (optional) | v1에서 그대로 확장 | `MeasurementResult`에 `risk_score` 추가 |
| 필드 제거 | Deprecated 주석 → 2 릴리스 후 제거 | `old_field` → `reserved 5;` |
| RPC 추가 | v1에서 그대로 확장 | `GetHealthSummary` 신규 |
| 필드 타입 변경 | v2 패키지 생성, v1 어댑터 유지 | `int32` → `google.protobuf.Timestamp` |
| RPC 시그니처 변경 | v2 패키지 생성, v1 어댑터 유지 | Request/Response 구조 변경 |

### 7.3 Gateway 버전 라우팅

```text
# Kong Gateway 라우팅 규칙
/api/v1/*  →  backend services (manpasik.v1 proto)
/api/v2/*  →  backend services (manpasik.v2 proto)  # Phase 4+

# 버전 전환 기간: 최소 6개월 병렬 운영
# Deprecation 알림: X-API-Deprecation 헤더
```

---

## 8. 오프라인→온라인 동기화 패턴

### 8.1 CRDT 기반 동기화

```text
오프라인 상태:
  Flutter App → Rust CRDT (sync/mod.rs) → 로컬 SQLite

온라인 복귀 시:
  1. Rust CRDT → gRPC StreamSync(changes[]) → measurement-service
  2. measurement-service: 서버 상태와 CRDT 병합 (LWW + Version Vector)
  3. 충돌 해결: 타임스탬프 기반 Last-Writer-Wins, 동시 변경 시 서버 우선
  4. 병합 완료 → 클라이언트에 최종 상태 반환
  5. 이벤트 발행: measurement.completed (지연 발행, delayed=true 플래그)
```

---

**참조**: `backend/shared/proto/manpasik.proto`, `docs/specs/event-schema-specification.md`, `infrastructure/docker-compose.yml`, `docs/specs/non-functional-requirements.md`
