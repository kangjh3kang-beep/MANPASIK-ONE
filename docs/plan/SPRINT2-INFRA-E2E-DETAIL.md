# Sprint 2 — 인프라 갱신 및 E2E 테스트 확장 세부 구현 기획서

> **문서 버전**: v1.0  
> **작성일**: 2026-02-12  
> **프로젝트**: ManPaSik (만파식) — 의료 IoT 디바이스 + 건강관리 플랫폼  
> **범위**: I-1 Docker Compose 갱신, I-3 E2E 테스트 확장, I-5 K8s Overlay 갱신

---

## 목차

1. [개요](#1-개요)
2. [I-1: Docker Compose 전체 서비스 목록](#2-i-1-docker-compose-전체-서비스-목록)
3. [I-3: E2E 테스트 파일별 명세](#3-i-3-e2e-테스트-파일별-명세)
4. [I-5: K8s 서비스 매니페스트 목록](#4-i-5-k8s-서비스-매니페스트-목록)
5. [환경변수 통합 목록](#5-환경변수-통합-목록)
6. [네트워크 토폴로지 및 통신 매트릭스](#6-네트워크-토폴로지-및-통신-매트릭스)
7. [CI/CD 파이프라인 갱신 사항](#7-cicd-파이프라인-갱신-사항)
8. [검증 기준](#8-검증-기준)

---

## 1. 개요

### 1.1 현재 상태 분석

| 항목 | 현재 상태 | 목표 |
|------|----------|------|
| Docker Compose 서비스 | 인프라 10개 + 백엔드 21개 + 게이트웨이 1개 = 32개 | 누락 서비스 추가 + 환경변수 보강 + 리소스 제한 |
| E2E 테스트 | 7개 파일 (gRPC 직접 호출 1개 + 이벤트버스 기반 4개 + REST 1개 + 헬스체크 1개) | gRPC 직접 호출 테스트 8개 파일 이상으로 확장 |
| K8s 서비스 YAML | 22개 파일 (gateway + 20 서비스 + deployments.yaml) | 누락 서비스 추가 + HPA/PDB 보완 |
| K8s Overlay | dev/staging/production 3개 환경 | Overlay 완성 + 환경별 차별화 |

### 1.2 구현된 Go 서비스 전체 목록 (22개)

| # | 서비스명 | gRPC 포트 | Phase | 역할 |
|---|---------|----------|-------|------|
| 1 | gateway | 8080 (HTTP) | 1 | REST→gRPC 변환 API 게이트웨이 |
| 2 | auth-service | 50051 | 1 | 인증/인가 (JWT, Register, Login) |
| 3 | user-service | 50052 | 1 | 사용자 프로필 관리 |
| 4 | device-service | 50053 | 1 | IoT 디바이스 등록/관리 |
| 5 | measurement-service | 50054 | 1 | 측정 세션/데이터 관리 |
| 6 | subscription-service | 50055 | 2 | 구독 플랜 관리 |
| 7 | shop-service | 50056 | 2 | 카트리지/리더기 쇼핑몰 |
| 8 | payment-service | 50057 | 2 | 결제 (Toss Payments 연동) |
| 9 | ai-inference-service | 50058 | 2 | AI 분석/건강 스코어 |
| 10 | cartridge-service | 50059 | 2 | NFC 카트리지 관리 |
| 11 | calibration-service | 50060 | 2 | 디바이스 보정 |
| 12 | coaching-service | 50061 | 2 | AI 건강 코칭 |
| 13 | family-service | 50063 | 3 | 가족 그룹/데이터 공유 |
| 14 | health-record-service | 50064 | 3 | 건강 기록 (FHIR 호환) |
| 15 | community-service | 50065 | 3 | 커뮤니티 게시판/챌린지 |
| 16 | reservation-service | 50066 | 3 | 병원/약국 예약 |
| 17 | admin-service | 50067 | 3 | 관리자 대시보드 |
| 18 | notification-service | 50068 | 3 | 알림 (Push/Email/SMS) |
| 19 | prescription-service | 50069 | 3 | 처방전 관리 |
| 20 | video-service | 50070 | 3 | 화상 상담 (WebRTC) |
| 21 | telemedicine-service | 50071 | 3 | 원격 진료 세션 |
| 22 | translation-service | 50073 | 3 | 의료 번역 (다국어) |

> **참고**: vision-service (cmd/main.go 존재)는 현재 Docker Compose에 미포함. 포트 50072 배정 예정.

---

## 2. I-1: Docker Compose 전체 서비스 목록

### 2.1 현재 상태 vs 변경 사항

현재 `infrastructure/docker/docker-compose.dev.yml`에 이미 대부분의 서비스가 정의되어 있으나, 다음 항목이 누락/보완 필요:

#### 2.1.1 누락 서비스 추가

| 서비스 | 포트 | 설명 |
|--------|------|------|
| vision-service | 50072 | 음식 칼로리 추정 비전 AI |

#### 2.1.2 게이트웨이 환경변수 보강

현재 게이트웨이에 누락된 서비스 주소:

```yaml
# 추가 필요한 게이트웨이 환경변수
SUBSCRIPTION_SERVICE_ADDR: "manpasik-subscription-service:50055"
SHOP_SERVICE_ADDR: "manpasik-shop-service:50056"
PAYMENT_SERVICE_ADDR: "manpasik-payment-service:50057"
AI_INFERENCE_SERVICE_ADDR: "manpasik-ai-inference-service:50058"
CARTRIDGE_SERVICE_ADDR: "manpasik-cartridge-service:50059"
CALIBRATION_SERVICE_ADDR: "manpasik-calibration-service:50060"
COACHING_SERVICE_ADDR: "manpasik-coaching-service:50061"
VISION_SERVICE_ADDR: "manpasik-vision-service:50072"
```

#### 2.1.3 기존 서비스 환경변수 보강

각 서비스에 공통적으로 추가해야 할 환경변수:

```yaml
# Redis (세션/캐시)
REDIS_HOST: redis
REDIS_PORT: "6379"

# Kafka (이벤트 발행)
KAFKA_BROKERS: "redpanda:9092"

# S3/MinIO (파일 스토리지)
S3_ENDPOINT: "minio:9000"
S3_ACCESS_KEY: manpasik
S3_SECRET_KEY: manpasik_dev_2026

# 모니터링
OTEL_EXPORTER_OTLP_ENDPOINT: ""
METRICS_PORT: "9100"
```

서비스별 추가 환경변수:

| 서비스 | 추가 환경변수 |
|--------|-------------|
| payment-service | `TOSS_SECRET_KEY`, `TOSS_API_URL` |
| ai-inference-service | `LLM_API_KEY`, `MILVUS_HOST`, `MILVUS_PORT` |
| coaching-service | `AI_INFERENCE_SERVICE_ADDR` |
| notification-service | `FCM_SERVER_KEY`, `SMTP_HOST`, `SMTP_PORT`, `SMTP_PASSWORD` |
| video-service | `VIDEO_STREAM_SECRET`, `TURN_SERVER_URL` |
| translation-service | `TRANSLATION_API_KEY` |
| measurement-service | `MILVUS_HOST`, `MILVUS_PORT`, `ELASTICSEARCH_URL`, `TIMESCALE_HOST` |
| health-record-service | `ELASTICSEARCH_URL` |
| admin-service | `CONFIG_ENCRYPTION_KEY` |

### 2.2 Docker Compose 전체 서비스 정의 (변경 대상)

#### 2.2.1 vision-service (신규 추가)

```yaml
vision-service:
  build:
    context: ../../backend
    dockerfile: services/vision-service/Dockerfile
  image: manpasik/vision-service:dev
  container_name: manpasik-vision-service
  environment:
    GRPC_PORT: ":50072"
    DB_HOST: postgres
    DB_PORT: "5432"
    DB_USER: manpasik
    DB_PASSWORD: manpasik_dev_2026
    DB_NAME: manpasik
    DB_SSLMODE: disable
    REDIS_HOST: redis
    REDIS_PORT: "6379"
    S3_ENDPOINT: "minio:9000"
    S3_ACCESS_KEY: manpasik
    S3_SECRET_KEY: manpasik_dev_2026
    S3_BUCKET: manpasik-vision
    LLM_API_KEY: "${LLM_API_KEY:-}"
  ports:
    - "50072:50072"
  depends_on:
    postgres:
      condition: service_healthy
  networks:
    - manpasik-network
  restart: unless-stopped
  deploy:
    resources:
      limits:
        cpus: '1.0'
        memory: 1G
      reservations:
        cpus: '0.25'
        memory: 256M
```

#### 2.2.2 기존 서비스 환경변수 보강 패턴

모든 Phase 2/3 서비스에 다음 환경변수 블록을 공통 추가:

```yaml
# 공통 추가 환경변수 (각 서비스에 반영)
REDIS_HOST: redis
REDIS_PORT: "6379"
KAFKA_BROKERS: "redpanda:9092"
LOG_LEVEL: "debug"
ENVIRONMENT: "development"
```

#### 2.2.3 리소스 제한 (deploy.resources)

| 서비스 카테고리 | CPU 제한 | 메모리 제한 | CPU 예약 | 메모리 예약 |
|---------------|---------|-----------|---------|-----------|
| gateway | 1.0 | 1G | 0.25 | 256M |
| Phase 1 핵심 (auth, measurement) | 0.5 | 512M | 0.1 | 128M |
| Phase 2 (payment, ai-inference) | 1.0 | 1G | 0.25 | 256M |
| Phase 2 기타 | 0.5 | 512M | 0.1 | 128M |
| Phase 3 일반 | 0.5 | 512M | 0.1 | 128M |
| video-service | 1.0 | 1G | 0.25 | 256M |
| 인프라 (PostgreSQL) | 2.0 | 2G | 0.5 | 512M |
| 인프라 (Redis) | 0.5 | 256M | 0.1 | 64M |
| 인프라 (Elasticsearch) | 2.0 | 2G | 0.5 | 512M |
| 인프라 (Milvus) | 2.0 | 2G | 0.5 | 512M |

#### 2.2.4 Healthcheck 추가

gRPC 서비스 공통 healthcheck 패턴:

```yaml
healthcheck:
  test: ["CMD", "/bin/grpc_health_probe", "-addr=:${GRPC_PORT}"]
  interval: 15s
  timeout: 5s
  retries: 3
  start_period: 10s
```

#### 2.2.5 게이트웨이 depends_on 완성

```yaml
gateway:
  depends_on:
    - auth-service
    - user-service
    - device-service
    - measurement-service
    - subscription-service    # 추가
    - shop-service            # 추가
    - payment-service         # 추가
    - ai-inference-service    # 추가
    - cartridge-service       # 추가
    - calibration-service     # 추가
    - coaching-service        # 추가
    - family-service
    - health-record-service
    - community-service
    - reservation-service
    - admin-service
    - notification-service
    - prescription-service
    - video-service
    - telemedicine-service
    - translation-service
    - vision-service          # 추가
```

### 2.3 네트워크 구성

모든 서비스가 단일 `manpasik-network` (bridge) 네트워크를 공유. 서비스 간 통신은 컨테이너명 DNS 사용.

---

## 3. I-3: E2E 테스트 파일별 명세

### 3.1 현재 E2E 테스트 분석

| 파일명 | 빌드 태그 | 방식 | 내용 |
|--------|----------|------|------|
| `flow_test.go` | 없음 | gRPC 직접 호출 | Register → Login → ValidateToken → StartSession → EndSession → GetHistory |
| `health_test.go` | 없음 | gRPC HealthCheck | 18개 서비스 헬스체크 + 차동측정 계산 |
| `commerce_flow_test.go` | `integration` | EventBus | subscription → shop → payment 이벤트 체인 |
| `ai_hardware_flow_test.go` | `integration` | EventBus | measurement → ai → coaching, cartridge → calibration |
| `community_admin_flow_test.go` | `integration` | EventBus | post → comment → notification, admin → notification, consent grant/revoke |
| `medical_flow_test.go` | `integration` | EventBus | reservation → prescription → pharmacy 이벤트 |
| `gateway_rest_test.go` | `integration` | HTTP REST | Gateway REST 엔드포인트 검증 |

### 3.2 확장 대상 테스트 목록

#### 3.2.1 `env.go` 확장

추가 주소 헬퍼 함수:

```go
func VideoAddr() string         { return getEnv("VIDEO_SERVICE_ADDR", "127.0.0.1:50070") }
func TelemedicineAddr() string  { return getEnv("TELEMEDICINE_SERVICE_ADDR", "127.0.0.1:50071") }
func VisionAddr() string        { return getEnv("VISION_SERVICE_ADDR", "127.0.0.1:50072") }
func TranslationAddr() string   { return getEnv("TRANSLATION_SERVICE_ADDR", "127.0.0.1:50073") }
```

#### 3.2.2 `health_test.go` 확장

```go
// TestServiceHealth에 추가할 서비스
{"video-service", VideoAddr()},
{"telemedicine-service", TelemedicineAddr()},
{"vision-service", VisionAddr()},
{"translation-service", TranslationAddr()},
```

---

### 3.3 신규 E2E 테스트 파일 명세

#### 3.3.1 `payment_flow_test.go` — 결제 플로우 E2E

| 항목 | 내용 |
|------|------|
| **파일 경로** | `backend/tests/e2e/payment_flow_test.go` |
| **빌드 태그** | 없음 (gRPC 직접 호출) |
| **선행 조건** | auth-service, payment-service 기동 |

**테스트 함수:**

```
TestPaymentFlow_CreateToRefund
```

**시나리오:**

1. **사전 준비**: auth-service에서 Register → Login → ValidateToken으로 `user_id` 획득
2. **CreatePayment**: PaymentService/CreatePayment 호출
   - 입력: `user_id`, `order_id="test-order-001"`, `payment_type=ONE_TIME`, `amount_krw=29900`, `payment_method="card"`
   - 검증: `payment_id` 비어있지 않음, `status == PAYMENT_STATUS_PENDING`
3. **ConfirmPayment**: PaymentService/ConfirmPayment 호출
   - 입력: `payment_id`, `payment_key="test-payment-key"`, `pg_provider="toss"`
   - 검증: `status == PAYMENT_STATUS_COMPLETED`
4. **GetPayment**: PaymentService/GetPayment 호출
   - 입력: `payment_id`
   - 검증: 반환된 결제 상세가 이전 단계 값과 일치
5. **RefundPayment**: PaymentService/RefundPayment 호출
   - 입력: `payment_id`, `refund_amount_krw=0` (전액 환불), `reason="테스트 환불"`
   - 검증: `refund_id` 비어있지 않음, `payment_status == PAYMENT_STATUS_REFUNDED`

**gRPC 엔드포인트:**

| RPC | 메서드 경로 |
|-----|-----------|
| CreatePayment | `/manpasik.v1.PaymentService/CreatePayment` |
| ConfirmPayment | `/manpasik.v1.PaymentService/ConfirmPayment` |
| GetPayment | `/manpasik.v1.PaymentService/GetPayment` |
| RefundPayment | `/manpasik.v1.PaymentService/RefundPayment` |

---

#### 3.3.2 `subscription_flow_test.go` — 구독 플로우 E2E

| 항목 | 내용 |
|------|------|
| **파일 경로** | `backend/tests/e2e/subscription_flow_test.go` |
| **빌드 태그** | 없음 (gRPC 직접 호출) |
| **선행 조건** | auth-service, subscription-service 기동 |

**테스트 함수:**

```
TestSubscriptionFlow_CreateToCancel
```

**시나리오:**

1. **사전 준비**: Register → Login → ValidateToken으로 `user_id` 획득
2. **CreateSubscription**: SubscriptionService/CreateSubscription 호출
   - 입력: `user_id`, `tier=SUBSCRIPTION_TIER_FREE`
   - 검증: `subscription_id` 비어있지 않음, `tier == FREE`, `status == ACTIVE`
3. **UpdateSubscription (업그레이드)**: SubscriptionService/UpdateSubscription 호출
   - 입력: `user_id`, `new_tier=SUBSCRIPTION_TIER_PRO`, `payment_id="test-pay-001"`
   - 검증: `tier == PRO`, `ai_coaching_enabled == true`
4. **GetSubscription**: SubscriptionService/GetSubscription 호출
   - 입력: `user_id`
   - 검증: 반환 정보가 업그레이드 후 상태와 일치
5. **CheckFeatureAccess**: SubscriptionService/CheckFeatureAccess 호출
   - 입력: `user_id`, `feature_name="ai_coaching"`
   - 검증: `allowed == true`
6. **CancelSubscription**: SubscriptionService/CancelSubscription 호출
   - 입력: `user_id`, `reason="테스트 해지"`
   - 검증: `success == true`, `effective_until` 미래 시점

**gRPC 엔드포인트:**

| RPC | 메서드 경로 |
|-----|-----------|
| CreateSubscription | `/manpasik.v1.SubscriptionService/CreateSubscription` |
| UpdateSubscription | `/manpasik.v1.SubscriptionService/UpdateSubscription` |
| GetSubscription | `/manpasik.v1.SubscriptionService/GetSubscription` |
| CheckFeatureAccess | `/manpasik.v1.SubscriptionService/CheckFeatureAccess` |
| CancelSubscription | `/manpasik.v1.SubscriptionService/CancelSubscription` |

---

#### 3.3.3 `admin_config_flow_test.go` — 관리자 설정 플로우 E2E

| 항목 | 내용 |
|------|------|
| **파일 경로** | `backend/tests/e2e/admin_config_flow_test.go` |
| **빌드 태그** | 없음 (gRPC 직접 호출) |
| **선행 조건** | admin-service 기동 |

**테스트 함수:**

```
TestAdminConfigFlow_SetGetConfig
TestAdminFlow_CreateAdminAndAudit
```

**시나리오 1 — TestAdminConfigFlow_SetGetConfig:**

1. **SetSystemConfig**: AdminService/SetSystemConfig 호출
   - 입력: `key="max_devices_per_user"`, `value="5"`, `description="사용자당 최대 디바이스 수"`
   - 검증: 반환된 `SystemConfig.key == "max_devices_per_user"`, `value == "5"`
2. **GetSystemConfig**: AdminService/GetSystemConfig 호출
   - 입력: `key="max_devices_per_user"`
   - 검증: `value == "5"`, `description` 일치
3. **SetSystemConfig (업데이트)**: 동일 키로 값 변경
   - 입력: `key="max_devices_per_user"`, `value="10"`
   - 검증: `value == "10"`
4. **GetSystemConfig (재확인)**: 변경 값 확인
   - 검증: `value == "10"`

**시나리오 2 — TestAdminFlow_CreateAdminAndAudit:**

1. **CreateAdmin**: AdminService/CreateAdmin 호출
   - 입력: `email="admin-e2e@manpasik.test"`, `password="Admin123!"`, `display_name="E2E Admin"`, `role=ADMIN_ROLE_REGIONAL`, `region="seoul"`
   - 검증: `admin_id` 비어있지 않음
2. **GetAdmin**: AdminService/GetAdmin 호출
   - 입력: `admin_id`
   - 검증: `email`, `role`, `region` 일치
3. **GetSystemStats**: AdminService/GetSystemStats 호출
   - 검증: 응답에 `total_users`, `total_devices` 필드 존재
4. **GetAuditLog**: AdminService/GetAuditLog 호출
   - 입력: `admin_id`, `limit=10`
   - 검증: 로그 항목 존재 확인

**gRPC 엔드포인트:**

| RPC | 메서드 경로 |
|-----|-----------|
| SetSystemConfig | `/manpasik.v1.AdminService/SetSystemConfig` |
| GetSystemConfig | `/manpasik.v1.AdminService/GetSystemConfig` |
| CreateAdmin | `/manpasik.v1.AdminService/CreateAdmin` |
| GetAdmin | `/manpasik.v1.AdminService/GetAdmin` |
| GetSystemStats | `/manpasik.v1.AdminService/GetSystemStats` |
| GetAuditLog | `/manpasik.v1.AdminService/GetAuditLog` |

---

#### 3.3.4 `device_measurement_flow_test.go` — 디바이스+측정 확장 플로우 E2E

| 항목 | 내용 |
|------|------|
| **파일 경로** | `backend/tests/e2e/device_measurement_flow_test.go` |
| **빌드 태그** | 없음 (gRPC 직접 호출) |
| **선행 조건** | auth-service, device-service, measurement-service 기동 |

**테스트 함수:**

```
TestDeviceMeasurementFlow_RegisterToHistory
```

**시나리오:**

1. **사전 준비**: Register → Login → ValidateToken으로 `user_id` 획득
2. **RegisterDevice**: DeviceService/RegisterDevice 호출
   - 입력: `device_id="e2e-dev-001"`, `serial_number="SN-E2E-2026"`, `firmware_version="1.0.0"`, `user_id`
   - 검증: `device_id` 비어있지 않음, `registration_token` 존재
3. **ListDevices**: DeviceService/ListDevices 호출
   - 입력: `user_id`
   - 검증: 방금 등록한 디바이스가 목록에 포함
4. **StartSession**: MeasurementService/StartSession 호출
   - 입력: `device_id="e2e-dev-001"`, `cartridge_id="e2e-cart-001"`, `user_id`
   - 검증: `session_id` 비어있지 않음
5. **EndSession**: MeasurementService/EndSession 호출
   - 입력: `session_id`
   - 검증: `session_id` 일치
6. **GetMeasurementHistory**: MeasurementService/GetMeasurementHistory 호출
   - 입력: `user_id`, `limit=10`
   - 검증: 응답 정상 수신 (total_count >= 0)

**gRPC 엔드포인트:**

| RPC | 메서드 경로 |
|-----|-----------|
| RegisterDevice | `/manpasik.v1.DeviceService/RegisterDevice` |
| ListDevices | `/manpasik.v1.DeviceService/ListDevices` |
| StartSession | `/manpasik.v1.MeasurementService/StartSession` |
| EndSession | `/manpasik.v1.MeasurementService/EndSession` |
| GetMeasurementHistory | `/manpasik.v1.MeasurementService/GetMeasurementHistory` |

---

#### 3.3.5 `notification_flow_test.go` — 알림 플로우 E2E

| 항목 | 내용 |
|------|------|
| **파일 경로** | `backend/tests/e2e/notification_flow_test.go` |
| **빌드 태그** | 없음 (gRPC 직접 호출) |
| **선행 조건** | auth-service, notification-service 기동 |

**테스트 함수:**

```
TestNotificationFlow_SendAndRead
```

**시나리오:**

1. **사전 준비**: Register → Login → ValidateToken으로 `user_id` 획득
2. **SendNotification**: NotificationService/SendNotification 호출
   - 입력: `user_id`, `type=NOTIFICATION_TYPE_HEALTH_ALERT`, `title="건강 경고"`, `body="혈당이 정상 범위를 초과했습니다"`, `priority=HIGH`, `channel=IN_APP`
   - 검증: `notification_id` 비어있지 않음
3. **GetUnreadCount**: NotificationService/GetUnreadCount 호출
   - 입력: `user_id`
   - 검증: `count >= 1`
4. **ListNotifications**: NotificationService/ListNotifications 호출
   - 입력: `user_id`, `limit=10`
   - 검증: 방금 보낸 알림이 목록에 포함
5. **MarkAsRead**: NotificationService/MarkAsRead 호출
   - 입력: `notification_id`
   - 검증: `success == true`
6. **GetUnreadCount (재확인)**: 읽지 않은 수 감소 확인

**gRPC 엔드포인트:**

| RPC | 메서드 경로 |
|-----|-----------|
| SendNotification | `/manpasik.v1.NotificationService/SendNotification` |
| GetUnreadCount | `/manpasik.v1.NotificationService/GetUnreadCount` |
| ListNotifications | `/manpasik.v1.NotificationService/ListNotifications` |
| MarkAsRead | `/manpasik.v1.NotificationService/MarkAsRead` |

---

#### 3.3.6 `community_flow_test.go` — 커뮤니티 플로우 E2E (gRPC 직접 호출)

| 항목 | 내용 |
|------|------|
| **파일 경로** | `backend/tests/e2e/community_grpc_flow_test.go` |
| **빌드 태그** | 없음 (gRPC 직접 호출) |
| **선행 조건** | auth-service, community-service 기동 |

**테스트 함수:**

```
TestCommunityFlow_PostCommentLike
```

**시나리오:**

1. **사전 준비**: Register → Login → ValidateToken으로 `user_id` 획득
2. **CreatePost**: CommunityService/CreatePost 호출
   - 입력: `author_id=user_id`, `title="E2E 테스트 게시글"`, `content="만파식 E2E 테스트입니다"`, `category=POST_CATEGORY_GENERAL`
   - 검증: `post_id` 비어있지 않음
3. **GetPost**: CommunityService/GetPost 호출
   - 입력: `post_id`
   - 검증: `title`, `content` 일치
4. **CreateComment**: CommunityService/CreateComment 호출
   - 입력: `post_id`, `author_id=user_id`, `content="댓글 테스트"`
   - 검증: `comment_id` 비어있지 않음
5. **LikePost**: CommunityService/LikePost 호출
   - 입력: `post_id`, `user_id`
   - 검증: `success == true`, `like_count >= 1`
6. **ListPosts**: CommunityService/ListPosts 호출
   - 검증: 방금 생성한 게시글이 목록에 포함

---

#### 3.3.7 `family_health_flow_test.go` — 가족+건강기록 플로우 E2E

| 항목 | 내용 |
|------|------|
| **파일 경로** | `backend/tests/e2e/family_health_flow_test.go` |
| **빌드 태그** | 없음 (gRPC 직접 호출) |
| **선행 조건** | auth-service, family-service, health-record-service 기동 |

**테스트 함수:**

```
TestFamilyHealthFlow_GroupToSharing
```

**시나리오:**

1. **사전 준비**: 2명의 사용자 Register/Login (owner, member)
2. **CreateFamilyGroup**: FamilyService/CreateFamilyGroup 호출
   - 입력: `owner_user_id`, `name="E2E 가족그룹"`, `description="테스트용"`
   - 검증: `group_id` 비어있지 않음
3. **InviteMember**: FamilyService/InviteMember 호출
   - 입력: `group_id`, `inviter_user_id=owner`, `invitee_email=member_email`, `role=FAMILY_ROLE_MEMBER`
   - 검증: `invitation_id` 비어있지 않음
4. **CreateRecord**: HealthRecordService/CreateRecord 호출
   - 입력: `user_id=owner`, `record_type=HEALTH_RECORD_TYPE_VITAL_SIGN`, `title="혈압 기록"`, `provider="자가 측정"`
   - 검증: `record_id` 비어있지 않음
5. **ListRecords**: HealthRecordService/ListRecords 호출
   - 입력: `user_id=owner`, `limit=10`
   - 검증: 생성한 기록이 목록에 포함
6. **GetHealthSummary**: HealthRecordService/GetHealthSummary 호출
   - 입력: `user_id=owner`
   - 검증: `total_records >= 1`

---

#### 3.3.8 `reservation_prescription_flow_test.go` — 예약+처방 플로우 E2E

| 항목 | 내용 |
|------|------|
| **파일 경로** | `backend/tests/e2e/reservation_prescription_flow_test.go` |
| **빌드 태그** | 없음 (gRPC 직접 호출) |
| **선행 조건** | auth-service, reservation-service, prescription-service 기동 |

**테스트 함수:**

```
TestReservationPrescriptionFlow
```

**시나리오:**

1. **사전 준비**: Register → Login → ValidateToken으로 `user_id` 획득
2. **CreateReservation**: ReservationService/CreateReservation 호출
   - 입력: `user_id`, `facility_id="fac-e2e"`, `slot_id="slot-e2e"`, `specialty=DOCTOR_SPECIALTY_INTERNAL`, `reason="정기 검진"`
   - 검증: `reservation_id` 비어있지 않음, `status == PENDING`
3. **GetReservation**: ReservationService/GetReservation 호출
   - 입력: `reservation_id`
   - 검증: 상세 정보 일치
4. **ListReservations**: ReservationService/ListReservations 호출
   - 입력: `user_id`, `limit=10`
   - 검증: 목록에 방금 생성한 예약 포함
5. **CreatePrescription**: PrescriptionService/CreatePrescription 호출
   - 입력: `user_id`, `doctor_id="doc-e2e"`, `doctor_name="김의사"`, `facility_id="fac-e2e"`, `diagnosis="일반 검진"`
   - 검증: `prescription_id` 비어있지 않음
6. **GetPrescription**: PrescriptionService/GetPrescription 호출
   - 입력: `prescription_id`
   - 검증: 상세 정보 일치
7. **CancelReservation**: ReservationService/CancelReservation 호출
   - 입력: `reservation_id`, `reason="테스트 취소"`
   - 검증: `success == true`

---

### 3.4 E2E 테스트 전체 매트릭스

| # | 파일명 | 태그 | 방식 | 대상 서비스 | 상태 |
|---|--------|------|------|-----------|------|
| 1 | `flow_test.go` | - | gRPC | auth, measurement | 기존 |
| 2 | `health_test.go` | - | gRPC Health | 전체 서비스 | 기존 (확장) |
| 3 | `commerce_flow_test.go` | integration | EventBus | subscription, shop, payment | 기존 |
| 4 | `ai_hardware_flow_test.go` | integration | EventBus | measurement, ai, coaching, cartridge, calibration | 기존 |
| 5 | `community_admin_flow_test.go` | integration | EventBus | community, admin, notification, family | 기존 |
| 6 | `medical_flow_test.go` | integration | EventBus | reservation, prescription | 기존 |
| 7 | `gateway_rest_test.go` | integration | HTTP REST | gateway (+ 전 서비스) | 기존 |
| 8 | **`payment_flow_test.go`** | - | gRPC | auth, payment | **신규** |
| 9 | **`subscription_flow_test.go`** | - | gRPC | auth, subscription | **신규** |
| 10 | **`admin_config_flow_test.go`** | - | gRPC | admin | **신규** |
| 11 | **`device_measurement_flow_test.go`** | - | gRPC | auth, device, measurement | **신규** |
| 12 | **`notification_flow_test.go`** | - | gRPC | auth, notification | **신규** |
| 13 | **`community_grpc_flow_test.go`** | - | gRPC | auth, community | **신규** |
| 14 | **`family_health_flow_test.go`** | - | gRPC | auth, family, health-record | **신규** |
| 15 | **`reservation_prescription_flow_test.go`** | - | gRPC | auth, reservation, prescription | **신규** |

---

## 4. I-5: K8s 서비스 매니페스트 목록

### 4.1 현재 상태

`infrastructure/kubernetes/base/services/`에 이미 22개 YAML 파일이 존재하며 기본 구조는 완성됨:

- `gateway.yaml` — Deployment + Service (LoadBalancer)
- `auth-service.yaml` ~ `translation-service.yaml` — 각각 Deployment + Service (ClusterIP)
- `deployments.yaml` — 추가 Deployment 정의

### 4.2 누락 서비스

| 서비스 | YAML 파일 | 상태 |
|--------|----------|------|
| vision-service | `vision-service.yaml` | **신규 추가 필요** |
| telemedicine-service | `telemedicine-service.yaml` | **신규 추가 필요** |

### 4.3 서비스별 K8s 스펙

#### 4.3.1 공통 Deployment 패턴

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {service-name}
  labels:
    app: {service-name}
    tier: backend
spec:
  replicas: 2          # base (overlay에서 환경별 조정)
  selector:
    matchLabels:
      app: {service-name}
  template:
    metadata:
      labels:
        app: {service-name}
        tier: backend
    spec:
      containers:
        - name: {service-name}
          image: manpasik/{service-name}:latest
          ports:
            - containerPort: {grpc-port}
              name: grpc
            - containerPort: 9100
              name: metrics
          envFrom:
            - configMapRef:
                name: manpasik-config
            - secretRef:
                name: manpasik-secrets
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 500m
              memory: 512Mi
          readinessProbe:
            exec:
              command: ["/bin/grpc_health_probe", "-addr=:{grpc-port}"]
            initialDelaySeconds: 5
            periodSeconds: 10
          livenessProbe:
            exec:
              command: ["/bin/grpc_health_probe", "-addr=:{grpc-port}"]
            initialDelaySeconds: 15
            periodSeconds: 20
---
apiVersion: v1
kind: Service
metadata:
  name: {service-name}
spec:
  selector:
    app: {service-name}
  ports:
    - port: {grpc-port}
      targetPort: {grpc-port}
      name: grpc
    - port: 9100
      targetPort: 9100
      name: metrics
```

#### 4.3.2 신규 추가 서비스 YAML

**vision-service.yaml**:

| 항목 | 값 |
|------|---|
| containerPort | 50072 |
| resources.requests | cpu: 200m, memory: 256Mi |
| resources.limits | cpu: 1000m, memory: 1Gi |
| 추가 env | `LLM_API_KEY` (secretRef) |

**telemedicine-service.yaml**:

| 항목 | 값 |
|------|---|
| containerPort | 50071 |
| resources.requests | cpu: 100m, memory: 128Mi |
| resources.limits | cpu: 500m, memory: 512Mi |

### 4.4 kustomization.yaml 갱신

`infrastructure/kubernetes/base/kustomization.yaml`에 추가:

```yaml
resources:
  # ... 기존 ...
  - services/vision-service.yaml       # 추가
  - services/telemedicine-service.yaml  # 추가
```

### 4.5 ConfigMap 갱신

`infrastructure/kubernetes/base/config/configmap.yaml`에 추가:

```yaml
data:
  # 기존 항목 유지 + 추가:
  VISION_SERVICE_ADDR: "vision-service:50072"
  VISION_SERVICE_URL: "vision-service.manpasik.svc.cluster.local:50072"
  TELEMEDICINE_SERVICE_ADDR: "telemedicine-service:50071"
  TELEMEDICINE_SERVICE_URL: "telemedicine-service.manpasik.svc.cluster.local:50071"

  # 신규 인프라 환경변수
  CONFIG_ENCRYPTION_KEY: ""     # Secrets로 이동 권장
  TOSS_API_URL: "https://api.tosspayments.com"
  TRANSLATION_API_PROVIDER: "google"
  VIDEO_TURN_SERVER_URL: ""
  MQTT_BROKER_URL: "mosquitto.manpasik-data.svc.cluster.local:1883"
```

### 4.6 Secrets 갱신

`infrastructure/kubernetes/base/config/secrets.yaml`에 추가:

```yaml
stringData:
  # 기존 항목 유지 + 추가:
  CONFIG_ENCRYPTION_KEY: "CHANGE_ME_AES256_CONFIG_KEY"
  TOSS_SECRET_KEY: "CHANGE_ME_TOSS_SECRET"
  LLM_API_KEY: "CHANGE_ME_LLM_API_KEY"
  TRANSLATION_API_KEY: "CHANGE_ME_TRANSLATION_KEY"
  TURN_SERVER_SECRET: "CHANGE_ME_TURN_SECRET"
  MQTT_USERNAME: "manpasik"
  MQTT_PASSWORD: "CHANGE_ME_MQTT_PASSWORD"
```

### 4.7 Overlay 환경별 차이

#### 4.7.1 dev (로컬 개발)

| 항목 | 값 |
|------|---|
| replicas | 1 (모든 서비스) |
| resources.limits | cpu: 250m, memory: 256Mi |
| DB_SSLMODE | disable |
| LOG_LEVEL | debug |

#### 4.7.2 staging (검증)

| 항목 | 값 |
|------|---|
| replicas | 2 (모든 서비스) |
| resources.limits | cpu: 500m, memory: 512Mi |
| DB_SSLMODE | require |
| LOG_LEVEL | info |

#### 4.7.3 production (운영)

| 항목 | 값 |
|------|---|
| replicas | 3 (핵심 서비스), 2 (기타) |
| resources.limits | 서비스별 차등 (아래 표) |
| DB_SSLMODE | require |
| LOG_LEVEL | warn |
| HPA | 핵심 7개 서비스 |
| PDB | 모든 서비스 minAvailable: 1 |

**Production 리소스 차등:**

| 서비스 | CPU 요청 | CPU 제한 | 메모리 요청 | 메모리 제한 | HPA min/max |
|--------|---------|---------|-----------|-----------|-------------|
| gateway | 200m | 1000m | 256Mi | 1Gi | 3/10 |
| auth-service | 200m | 500m | 256Mi | 512Mi | 3/8 |
| measurement-service | 200m | 1000m | 256Mi | 1Gi | 3/10 |
| payment-service | 200m | 500m | 256Mi | 512Mi | 3/6 |
| ai-inference-service | 500m | 2000m | 512Mi | 2Gi | 3/8 |
| video-service | 200m | 1000m | 256Mi | 1Gi | 3/8 |
| notification-service | 200m | 500m | 256Mi | 512Mi | 3/6 |
| 기타 서비스 | 100m | 500m | 128Mi | 512Mi | - |

---

## 5. 환경변수 통합 목록

### 5.1 카테고리별 환경변수

#### A. 데이터베이스 (PostgreSQL)

| 환경변수 | 설명 | 기본값 (dev) | 사용 서비스 |
|---------|------|------------|-----------|
| `DB_HOST` | PostgreSQL 호스트 | `postgres` | 전체 |
| `DB_PORT` | PostgreSQL 포트 | `5432` | 전체 |
| `DB_USER` | DB 사용자 | `manpasik` | 전체 |
| `DB_PASSWORD` | DB 비밀번호 | `manpasik_dev_2026` | 전체 |
| `DB_NAME` | DB 이름 | `manpasik` | 전체 |
| `DB_SSLMODE` | SSL 모드 | `disable` | 전체 |
| `DB_MAX_CONNS` | 최대 연결 수 | `20` | 전체 |
| `TIMESCALE_HOST` | TimescaleDB 호스트 | `timescaledb` | measurement |
| `TIMESCALE_PORT` | TimescaleDB 포트 | `5432` | measurement |

#### B. 캐시 (Redis)

| 환경변수 | 설명 | 기본값 (dev) | 사용 서비스 |
|---------|------|------------|-----------|
| `REDIS_HOST` | Redis 호스트 | `redis` | auth, 세션 관련 전체 |
| `REDIS_PORT` | Redis 포트 | `6379` | auth, 세션 관련 전체 |
| `REDIS_PASSWORD` | Redis 비밀번호 | (빈 값) | 전체 |
| `REDIS_DB` | Redis DB 번호 | `0` | 전체 |

#### C. 메시지 브로커 (Kafka/Redpanda)

| 환경변수 | 설명 | 기본값 (dev) | 사용 서비스 |
|---------|------|------------|-----------|
| `KAFKA_BROKERS` | Kafka 브로커 주소 | `redpanda:9092` | 이벤트 발행 전체 |
| `KAFKA_GROUP_ID` | Consumer Group ID | `{service-name}` | 전체 |

#### D. 인증/보안

| 환경변수 | 설명 | 기본값 (dev) | 사용 서비스 |
|---------|------|------------|-----------|
| `JWT_SECRET` | JWT 서명 키 | `dev-secret-change-in-production` | auth |
| `JWT_ACCESS_TTL_MINUTES` | 액세스 토큰 TTL | `15` | auth |
| `JWT_REFRESH_TTL_DAYS` | 리프레시 토큰 TTL | `7` | auth |
| `JWT_ISSUER` | JWT 발급자 | `manpasik-auth` | auth |
| `KEYCLOAK_URL` | Keycloak URL | `http://keycloak:8080` | auth |
| `KEYCLOAK_REALM` | Keycloak Realm | `manpasik` | auth |
| `KEYCLOAK_CLIENT_ID` | Keycloak Client ID | `manpasik-api` | auth |
| `KEYCLOAK_CLIENT_SECRET` | Keycloak Client Secret | (빈 값) | auth |
| `CONFIG_ENCRYPTION_KEY` | 설정 암호화 키 (AES-256) | (빈 값) | admin |
| `ENCRYPTION_MASTER_KEY` | 마스터 암호화 키 | (빈 값) | 전체 |

#### E. 벡터 DB (Milvus)

| 환경변수 | 설명 | 기본값 (dev) | 사용 서비스 |
|---------|------|------------|-----------|
| `MILVUS_HOST` | Milvus 호스트 | `milvus` | measurement, ai-inference |
| `MILVUS_PORT` | Milvus 포트 | `19530` | measurement, ai-inference |
| `MILVUS_COLLECTION` | 컬렉션 이름 | `manpasik_fingerprints` | measurement, ai-inference |

#### F. 검색 엔진 (Elasticsearch)

| 환경변수 | 설명 | 기본값 (dev) | 사용 서비스 |
|---------|------|------------|-----------|
| `ELASTICSEARCH_URL` | ES URL | `http://elasticsearch:9200` | measurement, health-record |
| `ELASTICSEARCH_USERNAME` | ES 사용자 | (빈 값) | measurement, health-record |
| `ELASTICSEARCH_PASSWORD` | ES 비밀번호 | (빈 값) | measurement, health-record |

#### G. 오브젝트 스토리지 (MinIO/S3)

| 환경변수 | 설명 | 기본값 (dev) | 사용 서비스 |
|---------|------|------------|-----------|
| `S3_ENDPOINT` | S3 엔드포인트 | `minio:9000` | video, vision, community |
| `S3_ACCESS_KEY` | S3 Access Key | `minioadmin` | video, vision, community |
| `S3_SECRET_KEY` | S3 Secret Key | `minioadmin` | video, vision, community |
| `S3_BUCKET` | S3 버킷 | `manpasik` | video, vision, community |
| `S3_REGION` | S3 리전 | `us-east-1` | video, vision, community |
| `S3_USE_SSL` | S3 SSL 사용 | `false` | video, vision, community |

#### H. 결제 (Toss Payments)

| 환경변수 | 설명 | 기본값 (dev) | 사용 서비스 |
|---------|------|------------|-----------|
| `TOSS_SECRET_KEY` | Toss 시크릿 키 | (빈 값) | payment |
| `TOSS_API_URL` | Toss API URL | `https://api.tosspayments.com` | payment |

#### I. AI/LLM

| 환경변수 | 설명 | 기본값 (dev) | 사용 서비스 |
|---------|------|------------|-----------|
| `LLM_API_KEY` | LLM API 키 | (빈 값) | ai-inference, coaching, vision |
| `AI_MODEL_API_KEY` | AI 모델 API 키 | (빈 값) | ai-inference |

#### J. 알림

| 환경변수 | 설명 | 기본값 (dev) | 사용 서비스 |
|---------|------|------------|-----------|
| `FCM_SERVER_KEY` | Firebase Cloud Messaging 키 | (빈 값) | notification |
| `SMTP_HOST` | SMTP 호스트 | (빈 값) | notification |
| `SMTP_PORT` | SMTP 포트 | `587` | notification |
| `SMTP_PASSWORD` | SMTP 비밀번호 | (빈 값) | notification |

#### K. 비디오/WebRTC

| 환경변수 | 설명 | 기본값 (dev) | 사용 서비스 |
|---------|------|------------|-----------|
| `VIDEO_STREAM_SECRET` | 비디오 스트림 시크릿 | (빈 값) | video |
| `TURN_SERVER_URL` | TURN 서버 URL | (빈 값) | video |
| `TURN_SERVER_SECRET` | TURN 서버 시크릿 | (빈 값) | video |

#### L. 번역

| 환경변수 | 설명 | 기본값 (dev) | 사용 서비스 |
|---------|------|------------|-----------|
| `TRANSLATION_API_KEY` | 번역 API 키 | (빈 값) | translation |
| `TRANSLATION_API_PROVIDER` | 번역 제공자 | `google` | translation |

#### M. IoT/MQTT

| 환경변수 | 설명 | 기본값 (dev) | 사용 서비스 |
|---------|------|------------|-----------|
| `MQTT_BROKER_URL` | MQTT 브로커 URL | `mosquitto:1883` | device |
| `MQTT_USERNAME` | MQTT 사용자 | `manpasik` | device |
| `MQTT_PASSWORD` | MQTT 비밀번호 | (빈 값) | device |

#### N. 서비스 메타

| 환경변수 | 설명 | 기본값 (dev) | 사용 서비스 |
|---------|------|------------|-----------|
| `GRPC_PORT` | gRPC 리슨 포트 | 서비스별 상이 | 전체 |
| `HTTP_PORT` | HTTP 리슨 포트 | `8080` | gateway |
| `LOG_LEVEL` | 로그 레벨 | `info` | 전체 |
| `ENVIRONMENT` | 환경 식별자 | `development` | 전체 |
| `METRICS_PORT` | Prometheus 메트릭 포트 | `9100` | 전체 |
| `SHUTDOWN_TIMEOUT_SECONDS` | Graceful shutdown 타임아웃 | `5` | 전체 |
| `MAX_FINGERPRINT_DIMENSION` | 핑거프린트 차원 | `896` | measurement, ai-inference |
| `DEFAULT_ALPHA` | 기본 알파 보정 계수 | `0.95` | measurement, calibration |
| `OFFLINE_MODE_ENABLED` | 오프라인 모드 | `true` | measurement |

#### O. 서비스 간 통신 주소

| 환경변수 | 값 (Docker Compose) | 값 (K8s) |
|---------|-------------------|----------|
| `AUTH_SERVICE_ADDR` | `manpasik-auth-service:50051` | `auth-service:50051` |
| `USER_SERVICE_ADDR` | `manpasik-user-service:50052` | `user-service:50052` |
| `DEVICE_SERVICE_ADDR` | `manpasik-device-service:50053` | `device-service:50053` |
| `MEASUREMENT_SERVICE_ADDR` | `manpasik-measurement-service:50054` | `measurement-service:50054` |
| `SUBSCRIPTION_SERVICE_ADDR` | `manpasik-subscription-service:50055` | `subscription-service:50055` |
| `SHOP_SERVICE_ADDR` | `manpasik-shop-service:50056` | `shop-service:50056` |
| `PAYMENT_SERVICE_ADDR` | `manpasik-payment-service:50057` | `payment-service:50057` |
| `AI_INFERENCE_SERVICE_ADDR` | `manpasik-ai-inference-service:50058` | `ai-inference-service:50058` |
| `CARTRIDGE_SERVICE_ADDR` | `manpasik-cartridge-service:50059` | `cartridge-service:50059` |
| `CALIBRATION_SERVICE_ADDR` | `manpasik-calibration-service:50060` | `calibration-service:50060` |
| `COACHING_SERVICE_ADDR` | `manpasik-coaching-service:50061` | `coaching-service:50061` |
| `FAMILY_SERVICE_ADDR` | `manpasik-family-service:50063` | `family-service:50063` |
| `HEALTH_RECORD_SERVICE_ADDR` | `manpasik-health-record-service:50064` | `health-record-service:50064` |
| `COMMUNITY_SERVICE_ADDR` | `manpasik-community-service:50065` | `community-service:50065` |
| `RESERVATION_SERVICE_ADDR` | `manpasik-reservation-service:50066` | `reservation-service:50066` |
| `ADMIN_SERVICE_ADDR` | `manpasik-admin-service:50067` | `admin-service:50067` |
| `NOTIFICATION_SERVICE_ADDR` | `manpasik-notification-service:50068` | `notification-service:50068` |
| `PRESCRIPTION_SERVICE_ADDR` | `manpasik-prescription-service:50069` | `prescription-service:50069` |
| `VIDEO_SERVICE_ADDR` | `manpasik-video-service:50070` | `video-service:50070` |
| `TELEMEDICINE_SERVICE_ADDR` | `manpasik-telemedicine-service:50071` | `telemedicine-service:50071` |
| `VISION_SERVICE_ADDR` | `manpasik-vision-service:50072` | `vision-service:50072` |
| `TRANSLATION_SERVICE_ADDR` | `manpasik-translation-service:50073` | `translation-service:50073` |

---

## 6. 네트워크 토폴로지 및 통신 매트릭스

### 6.1 gRPC 서비스 간 통신 매트릭스

| 호출자 ↓ / 대상 → | auth | user | device | meas | sub | shop | pay | ai | cart | calib | coach | family | hr | comm | rsrv | admin | notif | presc | video | tele | trans | vision |
|--|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|
| **gateway** | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● |
| **auth** | - | | | | | | | | | | | | | | | | | | | | | |
| **measurement** | ● | | ● | - | | | | ● | ● | ● | | | | | | | ● | | | | | |
| **device** | ● | | - | | | | | | ● | ● | | | | | | | ● | | | | | |
| **subscription** | ● | ● | | | - | | ● | | | | | | | | | | ● | | | | | |
| **payment** | ● | | | | ● | ● | - | | | | | | | | | | ● | | | | | |
| **coaching** | ● | ● | | ● | | | | ● | | | - | | | | | | ● | | | | | |
| **ai-inference** | | | | ● | | | | - | | | | | | | | | | | | | | ● |
| **notification** | | ● | | | | | | | | | | | | | | | - | | | | ● | |
| **reservation** | ● | ● | | | | | ● | | | | | | | | - | | ● | | | | ● | |
| **prescription** | ● | | | | | | ● | | | | | | ● | | | | ● | - | | | ● | |
| **telemedicine** | ● | ● | | | | | ● | | | | | | | | ● | | ● | ● | ● | - | ● | |
| **community** | ● | ● | | | | | | | | | | | | - | | | ● | | | | ● | |
| **admin** | ● | ● | | | ● | | | | | | | | | | | - | ● | | | | | |
| **family** | ● | ● | | ● | | | | | | | | - | ● | | | | ● | | | | | |
| **health-record** | ● | | | ● | | | | | | | | | - | | | | ● | | | | | |

> ● = 호출 관계 존재, - = 자기 자신, 빈칸 = 직접 호출 없음

### 6.2 Kafka 토픽 매트릭스

| 토픽 | 발행 서비스 | 구독 서비스 |
|------|-----------|-----------|
| `measurement.completed` | measurement | ai-inference, coaching, health-record |
| `calibration.completed` | calibration | device, measurement |
| `cartridge.replaced` | device | calibration |
| `ai.analysis_completed` | ai-inference | coaching, notification |
| `coaching.tip_delivered` | coaching | notification |
| `health_alert.triggered` | ai-inference, measurement | notification, family |
| `subscription.created` | subscription | shop, notification |
| `subscription.cancelled` | subscription | payment, notification |
| `payment.completed` | payment | subscription, shop, notification |
| `shop.order_created` | shop | payment, notification |
| `reservation.created` | reservation | notification |
| `reservation.cancelled` | reservation | notification |
| `prescription.created` | prescription | notification, health-record |
| `prescription.sent_to_pharmacy` | prescription | notification |
| `consent.granted` | health-record | family |
| `consent.revoked` | health-record | family |
| `family.data_shared` | family | notification |
| `community.post_created` | community | notification |
| `community.comment_created` | community | notification |
| `admin.action_performed` | admin | notification |
| `notification.sent` | notification | - (최종 소비) |

### 6.3 외부 포트 매핑 (Docker Compose)

| 호스트 포트 | 서비스 | 프로토콜 |
|-----------|--------|---------|
| 5432 | PostgreSQL | TCP |
| 5433 | TimescaleDB | TCP |
| 6379 | Redis | TCP |
| 8000 | Kong Proxy | HTTP |
| 8001 | Kong Admin API | HTTP |
| 8080 | Keycloak | HTTP |
| 8090 | Gateway | HTTP |
| 9000 | MinIO API | HTTP |
| 9010 | MinIO Console | HTTP |
| 9090 | Prometheus | HTTP |
| 9200 | Elasticsearch | HTTP |
| 19092 | Redpanda (Kafka) | TCP |
| 19530 | Milvus | TCP |
| 1883 | Mosquitto (MQTT) | TCP |
| 3000 | Grafana | HTTP |
| 50051-50073 | gRPC 서비스 | gRPC |

---

## 7. CI/CD 파이프라인 갱신 사항

### 7.1 CI 워크플로우 (`.github/workflows/ci.yml`) 변경

#### 7.1.1 `go-build` Job — Build all services 확장

현재 21개 서비스 빌드 중. 추가:

```yaml
- name: Build all services
  run: |
    # ... 기존 서비스 ...
    go build -v ./services/vision-service/cmd        # 추가
    go build -v ./services/telemedicine-service/cmd   # 기존에 있으나 확인
```

#### 7.1.2 `docker-build` Job — matrix.service 확장

```yaml
strategy:
  matrix:
    service:
      # ... 기존 22개 ...
      - vision-service       # 추가
```

#### 7.1.3 `e2e-test` Job — Docker Compose 기반 전환

현재 E2E 테스트는 서비스 미기동 시 Skip하는 방식. 개선안:

```yaml
e2e-test:
  name: E2E Test
  runs-on: ubuntu-latest
  needs: [go-build]

  services:
    postgres:
      image: postgres:16-alpine
      env:
        POSTGRES_USER: manpasik
        POSTGRES_PASSWORD: manpasik_dev_2026
        POSTGRES_DB: manpasik
      options: >-
        --health-cmd pg_isready
        --health-interval 10s
        --health-timeout 5s
        --health-retries 5
      ports:
        - 5432:5432

    redis:
      image: redis:7-alpine
      options: >-
        --health-cmd "redis-cli ping"
        --health-interval 10s
        --health-timeout 5s
        --health-retries 5
      ports:
        - 6379:6379

  steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: ${{ env.GO_VERSION }}

    - name: Start core services
      working-directory: backend
      run: |
        # 핵심 서비스만 백그라운드 실행
        go run ./services/auth-service/cmd &
        go run ./services/user-service/cmd &
        go run ./services/device-service/cmd &
        go run ./services/measurement-service/cmd &
        go run ./services/payment-service/cmd &
        go run ./services/subscription-service/cmd &
        go run ./services/admin-service/cmd &
        go run ./services/notification-service/cmd &
        go run ./services/community-service/cmd &
        go run ./services/family-service/cmd &
        go run ./services/health-record-service/cmd &
        go run ./services/reservation-service/cmd &
        go run ./services/prescription-service/cmd &
        sleep 15  # 서비스 기동 대기
      env:
        DB_HOST: localhost
        DB_PORT: "5432"
        DB_USER: manpasik
        DB_PASSWORD: manpasik_dev_2026
        DB_NAME: manpasik
        DB_SSLMODE: disable
        REDIS_HOST: localhost
        JWT_SECRET: test-secret

    - name: Run E2E tests
      working-directory: backend
      run: go test -v -count=1 -timeout 120s ./tests/e2e/...
      env:
        AUTH_SERVICE_ADDR: "127.0.0.1:50051"
        USER_SERVICE_ADDR: "127.0.0.1:50052"
        DEVICE_SERVICE_ADDR: "127.0.0.1:50053"
        MEASUREMENT_SERVICE_ADDR: "127.0.0.1:50054"
        SUBSCRIPTION_SERVICE_ADDR: "127.0.0.1:50055"
        PAYMENT_SERVICE_ADDR: "127.0.0.1:50057"
        ADMIN_SERVICE_ADDR: "127.0.0.1:50067"
        NOTIFICATION_SERVICE_ADDR: "127.0.0.1:50068"
        COMMUNITY_SERVICE_ADDR: "127.0.0.1:50065"
        FAMILY_SERVICE_ADDR: "127.0.0.1:50063"
        HEALTH_RECORD_SERVICE_ADDR: "127.0.0.1:50064"
        RESERVATION_SERVICE_ADDR: "127.0.0.1:50066"
        PRESCRIPTION_SERVICE_ADDR: "127.0.0.1:50069"

    - name: Run Integration tests (EventBus)
      working-directory: backend
      run: go test -tags=integration -v -count=1 -timeout 60s ./tests/e2e/...
```

#### 7.1.4 E2E 테스트 결과 리포팅

```yaml
    - name: Upload E2E test results
      if: always()
      uses: actions/upload-artifact@v4
      with:
        name: e2e-test-results
        path: backend/tests/e2e/
```

### 7.2 CD 워크플로우 (`.github/workflows/cd.yml`) 변경

#### 7.2.1 Staging 배포 — 서비스 검증 목록 갱신

```yaml
- name: Verify all services
  run: |
    SERVICES=(
      gateway auth-service user-service device-service measurement-service
      subscription-service shop-service payment-service
      ai-inference-service cartridge-service calibration-service coaching-service
      family-service notification-service health-record-service
      prescription-service reservation-service telemedicine-service
      admin-service community-service video-service translation-service
      vision-service  # 추가
    )
```

#### 7.2.2 Rollback — 서비스 목록 동기화

```yaml
- name: Rollback all services
  run: |
    SERVICES=(
      # ... 위와 동일한 23개 서비스 목록 ...
      vision-service  # 추가
    )
```

### 7.3 Makefile 갱신

#### 7.3.1 `GO_SERVICES` 변수 확장

```makefile
GO_SERVICES := gateway auth-service user-service device-service measurement-service \
  subscription-service shop-service payment-service ai-inference-service \
  cartridge-service calibration-service coaching-service \
  family-service health-record-service community-service \
  reservation-service admin-service notification-service \
  prescription-service video-service telemedicine-service \
  translation-service vision-service
```

#### 7.3.2 `build-go` 타겟 확장

```makefile
build-go:
	@echo "Go 서비스 빌드..."
	@for svc in $(GO_SERVICES); do \
		if [ "$$svc" = "gateway" ]; then \
			cd backend && $(GO) build $(GOFLAGS) $(LDFLAGS) -o ../bin/$$svc ./gateway/cmd; \
		else \
			cd backend && $(GO) build $(GOFLAGS) $(LDFLAGS) -o ../bin/$$svc ./services/$$svc/cmd; \
		fi; \
	done
	@echo "Go 빌드 완료"
```

#### 7.3.3 `docker-build` / `docker-push` 확장

```makefile
docker-build:
	@echo "Docker 이미지 빌드..."
	@for svc in $(GO_SERVICES); do \
		docker build -t $(DOCKER_REGISTRY)/$$svc:$(DOCKER_TAG) -f backend/services/$$svc/Dockerfile .; \
	done
	@echo "Docker 빌드 완료"
```

#### 7.3.4 신규 타겟 추가

```makefile
# E2E 테스트 (gRPC 직접 호출, 서비스 기동 필요)
test-e2e:
	@echo "E2E 테스트..."
	cd backend && $(GO) test -v -count=1 -timeout 120s ./tests/e2e/...
	@echo "E2E 테스트 완료"

# E2E 테스트 (integration 태그, EventBus 기반)
test-e2e-integration:
	@echo "E2E 통합 테스트 (EventBus)..."
	cd backend && $(GO) test -tags=integration -v -count=1 -timeout 60s ./tests/e2e/...
	@echo "E2E 통합 테스트 완료"

# Docker Compose로 전체 스택 기동 후 E2E 실행
test-e2e-full:
	@echo "전체 E2E (Docker Compose + 테스트)..."
	cd infrastructure/docker && $(DOCKER_COMPOSE) -f docker-compose.dev.yml up -d
	sleep 30
	cd backend && $(GO) test -v -count=1 -timeout 180s ./tests/e2e/...
	cd infrastructure/docker && $(DOCKER_COMPOSE) -f docker-compose.dev.yml down
	@echo "전체 E2E 완료"
```

---

## 8. 검증 기준

### 8.1 Docker Compose 검증

| # | 검증 항목 | 명령어 | 기대 결과 |
|---|---------|--------|---------|
| 1 | 전체 서비스 기동 | `docker compose -f infrastructure/docker/docker-compose.dev.yml up -d` | exit code 0 |
| 2 | 서비스 상태 확인 | `docker compose -f infrastructure/docker/docker-compose.dev.yml ps` | 모든 서비스 `Up` 또는 `Healthy` |
| 3 | 인프라 헬스체크 | PostgreSQL: `pg_isready`, Redis: `redis-cli ping` | 모두 정상 |
| 4 | gRPC 서비스 응답 | `grpc_health_probe -addr=localhost:{port}` (각 서비스) | `SERVING` |
| 5 | Gateway REST 응답 | `curl http://localhost:8090/health` | 200 OK |
| 6 | 로그 에러 없음 | `docker compose logs --tail=100 2>&1 \| grep -c "ERROR\|FATAL"` | 0 |

### 8.2 E2E 테스트 검증

| # | 검증 항목 | 명령어 | 기대 결과 |
|---|---------|--------|---------|
| 1 | gRPC 직접 호출 테스트 | `go test -v -count=1 ./tests/e2e/... -timeout 120s` | PASS (또는 서비스 미기동 시 SKIP) |
| 2 | EventBus 통합 테스트 | `go test -tags=integration -v ./tests/e2e/...` | 모든 테스트 PASS |
| 3 | 결제 플로우 | `TestPaymentFlow_CreateToRefund` | PASS |
| 4 | 구독 플로우 | `TestSubscriptionFlow_CreateToCancel` | PASS |
| 5 | 관리자 설정 플로우 | `TestAdminConfigFlow_SetGetConfig` | PASS |
| 6 | 디바이스+측정 플로우 | `TestDeviceMeasurementFlow_RegisterToHistory` | PASS |
| 7 | 알림 플로우 | `TestNotificationFlow_SendAndRead` | PASS |
| 8 | 커뮤니티 플로우 | `TestCommunityFlow_PostCommentLike` | PASS |
| 9 | 가족+건강기록 플로우 | `TestFamilyHealthFlow_GroupToSharing` | PASS |
| 10 | 예약+처방 플로우 | `TestReservationPrescriptionFlow` | PASS |

### 8.3 K8s 검증

| # | 검증 항목 | 명령어 | 기대 결과 |
|---|---------|--------|---------|
| 1 | Kustomize 빌드 (dev) | `kubectl kustomize infrastructure/kubernetes/overlays/dev/` | YAML 정상 출력 |
| 2 | Kustomize 빌드 (staging) | `kubectl kustomize infrastructure/kubernetes/overlays/staging/` | YAML 정상 출력 |
| 3 | Kustomize 빌드 (production) | `kubectl kustomize infrastructure/kubernetes/overlays/production/` | YAML 정상 출력 |
| 4 | dry-run 적용 | `kubectl apply -k infrastructure/kubernetes/overlays/dev/ --dry-run=client` | 에러 없음 |
| 5 | 리소스 수 확인 | Deployment 23개 + Service 23개 + HPA 7개 + PDB + ConfigMap + Secret | 수량 일치 |
| 6 | 포트 충돌 없음 | 모든 서비스 포트 유니크 | 중복 없음 |

### 8.4 CI/CD 검증

| # | 검증 항목 | 기대 결과 |
|---|---------|---------|
| 1 | CI 빌드 (모든 서비스) | 23개 서비스 빌드 성공 |
| 2 | CI E2E 테스트 | gRPC + EventBus 테스트 모두 통과 |
| 3 | Docker 이미지 빌드 | 23개 이미지 빌드 성공 |
| 4 | Staging 배포 | 모든 서비스 rollout 성공 |
| 5 | Production 배포 | Canary → Full rollout 성공 |

---

## 부록 A: 구현 우선순위

| 단계 | 작업 | 예상 소요 | 의존성 |
|------|------|---------|--------|
| 1 | Docker Compose 환경변수 보강 + vision-service 추가 | 2시간 | 없음 |
| 2 | Docker Compose 리소스 제한 + healthcheck 추가 | 1시간 | 단계 1 |
| 3 | `env.go` 확장 + `health_test.go` 서비스 추가 | 30분 | 없음 |
| 4 | `payment_flow_test.go` 작성 | 1시간 | Proto 정의 |
| 5 | `subscription_flow_test.go` 작성 | 1시간 | Proto 정의 |
| 6 | `admin_config_flow_test.go` 작성 | 1시간 | Proto 정의 |
| 7 | `device_measurement_flow_test.go` 작성 | 1시간 | Proto 정의 |
| 8 | `notification_flow_test.go` 작성 | 45분 | Proto 정의 |
| 9 | `community_grpc_flow_test.go` 작성 | 45분 | Proto 정의 |
| 10 | `family_health_flow_test.go` 작성 | 1시간 | Proto 정의 |
| 11 | `reservation_prescription_flow_test.go` 작성 | 1시간 | Proto 정의 |
| 12 | K8s vision/telemedicine YAML 추가 | 30분 | 없음 |
| 13 | K8s ConfigMap/Secrets 갱신 | 30분 | 없음 |
| 14 | K8s kustomization.yaml 갱신 | 15분 | 단계 12 |
| 15 | CI/CD 워크플로우 갱신 | 1시간 | 단계 1-14 |
| 16 | Makefile 갱신 | 30분 | 없음 |
| **합계** | | **약 13시간** | |

## 부록 B: 포트 할당 현황

| 포트 범위 | 용도 |
|----------|------|
| 50051-50061 | Phase 1-2 gRPC 서비스 |
| 50062 | (미사용, 예약) |
| 50063-50073 | Phase 3 gRPC 서비스 |
| 8080 | Keycloak / Gateway (내부) |
| 8090 | Gateway (외부 매핑) |
| 8000-8002 | Kong |
| 9000-9010 | MinIO |
| 9090-9100 | Prometheus / 메트릭 |
| 9200 | Elasticsearch |
| 19092 | Redpanda (Kafka) |
| 19530 | Milvus |
| 1883 | MQTT |
| 3000 | Grafana |
| 5432-5433 | PostgreSQL / TimescaleDB |
| 6379 | Redis |
