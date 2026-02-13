# ManPaSik 구현 패턴·유형·방식

> **용도**: 구현 내역·유형·방식을 실시간 공유·학습할 수 있도록 정리.  
> **업데이트**: 새 기능/레이어 추가 시 해당 패턴을 이 문서에 반영.

---

## 1. Flutter 앱

### 1.1 디렉터리 구조 (Feature-First)

```
lib/
├── core/                    # 앱 전역: 상수, 라우터, 서비스, Provider
│   ├── constants/
│   ├── router/
│   ├── services/            # gRPC 채널, 인터셉터
│   ├── providers/           # gRPC·Repository Provider
│   ├── theme/
│   └── utils/
├── features/
│   └── {feature}/           # auth, home, devices, measurement, settings, user
│       ├── data/            # Repository 구현체, DTO
│       ├── domain/          # Repository 인터페이스, 엔티티
│       └── presentation/    # Screen 위젯
├── generated/               # proto 수동 생성 (메시지·gRPC 스텁)
├── l10n/
└── shared/                  # 공통 Provider, 위젯, 모델
```

- **원칙**: 기능별로 `data` / `domain` / `presentation` 분리. 도메인은 인터페이스만, 데이터 소스 연동은 `data`에만 둠.

### 1.2 상태 관리 (Riverpod)

| 유형 | 사용처 | 방식 |
|------|--------|------|
| **StateNotifierProvider** | auth, theme, locale | `StateNotifier<T>` + `ref.read( dependencyProvider )` 주입 |
| **Provider** | Repository, GrpcClientManager | 싱글톤·무상태 서비스 |
| **FutureProvider** | 화면 데이터(목록·프로필) | `ref.watch(authProvider).userId` 등으로 의존, 비동기 1회 로드 |

- **의존성**: 인증이 필요한 Repository는 `accessTokenProvider: () => ref.read(authProvider).accessToken` 형태로 토큰 콜백 전달.

### 1.3 gRPC 클라이언트

| 구성요소 | 역할 | 위치 |
|----------|------|------|
| **GrpcClientManager** | 호스트/포트 설정, 서비스별 `ClientChannel` (auth/user/device/measurement) 생성·쇼다운 | `lib/core/services/grpc_client.dart` |
| **AuthInterceptor** | JWT Bearer 메타데이터 자동 첨부 (`ClientInterceptor`, `interceptUnary` / `interceptStreaming`) | `lib/core/services/auth_interceptor.dart` |
| **수동 생성 메시지·스텁** | protoc 미사용 환경 대응. `GeneratedMessage` 호환 `.pb.dart`, `Client` 상속 `.pbgrpc.dart` | `lib/generated/manpasik.pb.dart`, `manpasik.pbgrpc.dart` |

- **채널**: Auth 서비스는 인터셉터 없음. User/Device/Measurement는 `AuthInterceptor(tokenProvider)` 사용.
- **포트**: `AppConstants.grpcAuthPort`(50051) 등 상수 사용, 환경별 분리 가능하도록 설계.

### 1.4 Repository 패턴

- **인터페이스**: `lib/features/{feature}/domain/*_repository.dart` (추상 클래스).
- **구현체**: `lib/features/{feature}/data/*_repository_impl.dart` — gRPC 클라이언트 호출, 예외 → 도메인/UI 친화적 에러 변환.
- **Provider**: `lib/core/providers/grpc_provider.dart`에서 `grpcClientManagerProvider` → `*RepositoryProvider` (필요 시 `accessTokenProvider` 주입).
- **테스트**: `test/helpers/fake_repositories.dart`에 `Fake*Repository` 구현, 단위·위젯 테스트에서 `overrideWithValue(Fake*Repository())` 로 주입.

### 1.5 화면·데이터 연동

- **목록/프로필**: `FutureProvider`(예: `measurementHistoryProvider`, `deviceListProvider`, `userProfileProvider`)로 gRPC 1회 호출.
- **UI**: `ref.watch(*Provider).when(data: ..., loading: ..., error: ...)` 로 로딩/에러/빈 상태 처리.
- **인증 필요 RPC**: Repository가 이미 `AuthInterceptor`로 토큰 전달받음. 화면은 `ref.read(authProvider).userId` 등으로 인증 여부만 필요 시 사용.

### 1.6 테스트

- **단위**: Repository·GrpcClientManager·AuthInterceptor 등은 Fake 또는 mocktail Mock으로 gRPC 없이 검증.
- **위젯**: `ProviderScope(overrides: [ authRepositoryProvider.overrideWithValue(FakeAuthRepository()), authProvider.overrideWith(...) ])` 로 Notifier·Repository 교체 후 `pumpWidget`·`pumpAndSettle`.
- **테스트 수**: 50개+ 목표. auth/theme/locale/validators + repository + grpc_client + screen 위젯으로 구성.

### 1.7 Rust FFI (S5b)

- **실연동 전**: `lib/core/services/rust_ffi_stub.dart` — `RustFfiStub.bleScan()`, `nfcReadCartridge()` 등 스텁 제공. BLE/NFC UI는 이 스텁 호출.
- **실연동 시**: flutter_rust_bridge 코드생성 후 생성된 API로 스텁 교체.
- **차트**: `fl_chart` 사용. 측정 결과 화면(`/measurement/result`)에서 GetMeasurementHistory 기반 트렌드 라인 차트.

---

## 2. Go 백엔드

### 2.1 서비스 구조

```
backend/services/{auth|user|device|measurement}-service/
├── cmd/main.go          # 설정 로드, DB/저장소 선택, gRPC 서버 기동
└── internal/
    ├── handler/          # gRPC 핸들러 (요청 → 서비스 호출)
    ├── repository/
    │   ├── memory/       # 인메모리 구현 (개발/테스트)
    │   └── postgres/     # PostgreSQL 구현 (선택)
    └── service/          # 비즈니스 로직 (Repository 인터페이스 주입)
```

- **원칙**: `main` → handler → service → repository. Proto 메시지는 `backend/shared/gen/go/v1/` 공용 사용.

### 2.2 설정·DB

- **공통 설정**: `backend/shared/config/config.go` — `LoadFromEnv(serviceName)`. DB/Redis/Kafka/JWT 등 환경변수 기반.
- **DB 미연결 시**: 인메모리 저장소로 fallback (로그 출력 후 계속 기동).

### 2.3 gRPC·인증

- **인증 필요 RPC**: `middleware.AuthInterceptor(validator)` 로 토큰 검증. `TokenValidator`는 해당 서비스의 ValidateToken 로직 래핑.
- **Auth 서비스**: Register/Login/RefreshToken/Logout/ValidateToken. JWT Access/Refresh 발급·갱신.

### 2.4 이벤트 발행 (Kafka / Memory 폴백)

- **인터페이스**: 각 서비스의 `internal/service/*.go`에 `EventPublisher` 인터페이스 정의 (예: `PublishPaymentCompleted`, `PublishMeasurementCompleted` 등). 서비스는 인터페이스만 의존.
- **구현체**:
  - **Kafka**: `internal/repository/kafka/event_publisher.go` — `shared/events.KafkaEventBus` 사용, envelope(`event_id`, `event_type`, `version`, `timestamp`, `source`, `payload`) 형식으로 토픽 발행.
  - **Memory**: `internal/repository/memory/event_publisher.go` — no-op 구현 (로그 없이 즉시 `nil` 반환). 로컬/테스트용.
- **조건부 초기화** (`cmd/main.go`): `KAFKA_BROKERS` 설정 시 `events.NewKafkaEventBus` 연결 시도 → 실패 시 또는 미설정 시 `memory.NewEventPublisher()` 사용. 서비스는 `SetEventPublisher(ep)` 로 주입.
- **에러 처리**: 이벤트 발행 실패 시 비치명적 처리 — 로그 `Warn` 후 계속 진행. 주문/결제 등 트랜잭션 자체는 롤백하지 않음.
- **적용 서비스**: payment, subscription, device, measurement (및 reservation/prescription 등 옵션 이벤트 인터페이스).

---

## 3. 인프라·Docker

### 3.1 Go 서비스 이미지

- **Dockerfile 위치**: `backend/services/{service}/Dockerfile`.
- **빌드 컨텍스트**: 반드시 `backend/` (상위). 예: `docker build -f services/auth-service/Dockerfile -t manpasik/auth-service:dev .`
- **구성**: multi-stage — builder에서 `go build ./services/{service}/cmd`, 런타임 이미지는 Alpine + ca-certificates.
- **포트**: 서비스별 50051(auth), 50052(user), 50053(device), 50054(measurement).

### 3.2 Docker Compose

- **파일**: `infrastructure/docker/docker-compose.dev.yml`.
- **Go 4서비스**: `build.context: ../../backend`, `dockerfile: services/{service}/Dockerfile`. `depends_on: postgres (healthy)`, 환경변수로 `DB_*`, `GRPC_PORT`, `JWT_SECRET` 등 전달.

### 3.3 E2E 테스트

- **위치**: `backend/tests/e2e/` (backend 모듈 내 — shared/gen/go/v1 사용).
- **실행**: `cd backend && go test -v ./tests/e2e/...` 또는 `make test-integration`.
- **구성**: `health_test.go`(4서비스 헬스체크, 차동측정 단위), `flow_test.go`(Register→Login→ValidateToken→StartSession→EndSession→GetHistory). gRPC 클라이언트 미생성 환경에서는 `conn.Invoke(ctx, fullMethod, req, resp)` + v1 메시지 사용.
- **서비스 미기동 시**: 연결/헬스 실패 시 `t.Skipf`로 스킵하여 CI 통과. 서비스 기동 후 동일 명령으로 전체 플로우 검증 가능.

### 3.4 Kubernetes (Kustomize)

- **위치**: `infrastructure/kubernetes/base/` (ConfigMap, Secrets, Deployment/Service 템플릿), `overlays/{dev|staging|production}/` (환경별 패치).
- **패턴**: `kustomization.yaml`로 base 리소스 + overlay 패치. 서비스별 `services/{service}.yaml` (Deployment, Service). Ingress·ConfigMap·Secrets는 base에 두고 환경별 값만 overlay에서 덮어씀.
- **배포**: `kubectl apply -k overlays/dev` 등으로 환경별 적용. CI/CD에서 카나리/단계 배포 가능.

---

## 4. Proto·gRPC 계약

- **정의**: `backend/shared/proto/manpasik.proto` (Auth, User, Device, Measurement 서비스).
- **Go**: `shared/gen/go/v1/` 에 수동 생성 또는 protoc 생성 `.pb.go`, `_grpc.pb.go`.
- **Flutter**: protoc 미사용 시 `lib/generated/` 에 수동 작성 `.pb.dart`(메시지), `.pbgrpc.dart`(클라이언트). 서비스 경로는 `manpasik.v1.*` / 메서드 경로 `/{ServiceName}/{MethodName}`.

---

## 5. 문서·품질 갱신 규칙

- **구현 반영**: 새 레이어·패턴 추가 시 이 문서(`docs/implementation-patterns.md`)에 유형·방식·위치 추가.
- **CHANGELOG**: 작업 단위별로 변경 파일·결정 사항·다음 단계 기록.
- **CONTEXT.md**: 기술 스택·Stage 상태·진행 단계 요약 유지.
- **QUALITY_GATES.md**: Stage 완료 시 해당 Gate 체크리스트·통과일 갱신.

---

**마지막 업데이트**: 2026-02-12 (Phase D-1: Kafka EventPublisher + Memory 폴백 패턴 반영)
