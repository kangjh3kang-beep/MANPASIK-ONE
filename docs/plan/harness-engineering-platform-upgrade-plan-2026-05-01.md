# ManPaSik Harness Engineering Platform Upgrade Plan

**문서 ID**: MPS-HARNESS-UPGRADE-2026-05-01  
**작성일**: 2026-05-01  
**작성자**: Codex  
**목적**: 만파식 플랫폼을 화면 중심 구현에서 실제 운용 가능한 의료기기급 플랫폼으로 끌어올리기 위한 상세 구현, 구축, 검증 계획  
**범위**: Flutter 앱, Go Gateway/MSA, Rust Core/FFI, 데이터 저장소, 이벤트, 오프라인 동기화, 관측성, 규제 증적, E2E 하네스

---

## 1. 핵심 결론

현재 만파식 플랫폼은 메뉴, 화면, 서비스 디렉터리, 문서 체계가 넓게 구축되어 있다. 그러나 하네스 엔지니어링 관점에서는 아직 다음 간극이 남아 있다.

1. UI 메뉴는 많지만 일부 Provider가 demo/mock/stub/fallback에 의존한다.
2. Flutter-Rust-Go-DB-Milvus-Kafka 종단간 측정 흐름이 운영 기준으로 고정되어 있지 않다.
3. readiness, 장애 감지, 재시도, 감사 로그, traceability가 실제 의존성 상태를 충분히 반영하지 못한다.
4. 의료기기급 완성 기준인 요구사항-위험-테스트-릴리스 증적 연결이 자동화되어 있지 않다.

따라서 업그레이드 전략은 "화면 추가"가 아니라 **혈관형 연결성 강화**다. 모든 메뉴는 다음 완료 정의를 만족해야 한다.

```text
Route -> State Provider -> Repository -> REST/gRPC -> Service -> DB/Event/Search/Vector
      -> Observability -> Error/Offline State -> Test Evidence -> Regulatory Trace
```

---

## 2. 목표 상태

### 2.1 플랫폼 완성 기준

| 영역 | 목표 |
|---|---|
| 기능 연결성 | 모든 1차 메뉴와 핵심 2차 메뉴가 실제 Repository/API/Service/Storage와 연결 |
| 측정 실행성 | BLE/NFC/측정 세션/차동측정/핑거프린트/AI/저장/결과 표시가 하나의 골든 패스로 통과 |
| 오프라인 운용 | 네트워크 단절 시 측정, 저장, 결과 확인 가능. 복구 시 충돌 해결 및 서버 동기화 |
| 운영 가동성 | readiness가 DB, Redis, Kafka, Milvus, Elasticsearch, 하위 gRPC 상태를 반영 |
| 장애 대응 | 모든 핵심 흐름에 timeout, retry, circuit breaker, fallback 등급, 사용자 가시화 적용 |
| 규제 증적 | 요구사항 ID, 위험 ID, 테스트 ID, 빌드 ID, 감사 로그가 연결 |
| 품질 게이트 | mock 은퇴율, E2E 통과율, API 계약 일치율, observability 커버리지로 완료 판정 |

### 2.2 단계별 성숙도 목표

| 단계 | 기간 | 목표 완성도 |
|---|---:|---|
| H0 진단/지도화 | 1주 | 연결 맵 100%, mock register 100% |
| H1 측정 골든 패스 | 2주 | Measure E2E 90% 이상 |
| H2 메뉴별 실데이터 결선 | 4주 | 핵심 메뉴 API 결선 85% 이상 |
| H3 오프라인/동기화 | 3주 | 오프라인 측정/동기화 E2E 통과 |
| H4 운영/보안/규제 하네스 | 3주 | 운영 readiness, audit, traceability 적용 |
| H5 릴리스 안정화 | 2주 | smoke/regression/performance gate 통과 |

---

## 3. 하네스 레이어 모델

### 3.1 8계층 하네스

| 계층 | 이름 | 책임 | 완료 증거 |
|---|---|---|---|
| L1 | UX Route Harness | route, deep link, navigation fallback | route-map test, screenshot smoke |
| L2 | State Harness | Riverpod state, loading/error/offline/stale 상태 | provider unit/widget test |
| L3 | Client IO Harness | REST/gRPC, Rust FFI, local DB, secure storage | contract test, fault injection |
| L4 | Core Engine Harness | BLE/NFC packet, differential, fingerprint, AI | fixture test, FFI integration |
| L5 | Backend Harness | Gateway, gRPC service, repository, transaction | service test, API contract |
| L6 | Data/Event Harness | PostgreSQL/Timescale, Milvus, Redis, Kafka, ES, MinIO | migration, seed, event replay |
| L7 | Ops Harness | health, readiness, metrics, logs, tracing, alert | docker smoke, dashboard check |
| L8 | Compliance Harness | audit trail, traceability, risk control evidence | requirements-risk-test matrix |

### 3.2 완료 판정 공식

```text
Feature Completion Score =
  0.20 * route_coverage
+ 0.20 * real_data_binding
+ 0.15 * backend_persistence
+ 0.15 * failure_handling
+ 0.10 * offline_behavior
+ 0.10 * observability
+ 0.10 * test_and_trace_evidence
```

운영 릴리스 후보는 각 핵심 기능이 85점 이상이어야 한다. 측정, 인증, 결제, 의료 데이터 공유, 응급 알림은 92점 이상이어야 한다.

---

## 4. 현재 관찰된 핵심 보강 대상

| 대상 | 관찰 | 조치 |
|---|---|---|
| Flutter Provider | demo/mock 데이터와 예외 시 빈 값 반환 패턴 존재 | 명시적 DemoMode, RealMode, OfflineMode 분리 |
| Rust FFI | native 실패 시 stub/demo BLE로 자동 폴백 | 운영 빌드에서는 silent fallback 금지, 진단 이벤트 발행 |
| Measurement Service | DB/Milvus/Kafka 실패 시 memory fallback 가능 | 개발/테스트 허용, 운영에서는 readiness fail 및 degraded 표기 |
| Gateway | REST 라우트는 넓으나 CORS `*`, readiness 단순화 | env allowlist, dependency-aware readiness 적용 |
| 메뉴 연결 | route는 넓으나 일부 하위 기능은 실제 API 결선 불명확 | route-api-service-data-test 매트릭스 작성 |
| 규제 증적 | 문서 풍부, 자동 trace 부족 | REQ/RISK/TEST/BUILD ID 연결 |

---

## 5. 메뉴별 보강 계획

### 5.1 Home

**목표**: 사용자가 앱을 열었을 때 현재 건강, 장비, 동기화, 알림, 위험 상태를 한눈에 파악한다.

구현 항목:
1. `HomeDashboardProvider`를 신설해 측정 요약, 장비 상태, 알림, 추천, 동기화 상태를 통합한다.
2. 각 카드에 `dataSource`: live/cache/demo/offline 값을 포함한다.
3. 측정 기록 API 실패 시 빈 카드가 아니라 `stale` 상태와 마지막 성공 시각을 표시한다.
4. Home smoke test에서 핵심 CTA가 실제 route로 이동하는지 검증한다.

완료 기준:
- `/home -> /measure`, `/home -> /data`, `/home -> /devices`, `/home -> /notifications` 연결 테스트 통과
- 최소 5개 카드가 실제 Repository 또는 offline cache에서 데이터 수신

### 5.2 Measure

**목표**: 만파식의 핵심 골든 패스를 완성한다.

골든 패스:

```text
BLE scan
-> device connect
-> NFC cartridge read/validate
-> StartSession
-> raw packet stream
-> Rust parser
-> differential correction
-> fingerprint build
-> AI local inference
-> ProcessMeasurement
-> Timescale/PostgreSQL store
-> Milvus vector store
-> Kafka measurement.completed
-> result screen
-> history/data hub 반영
```

구현 항목:
1. `MeasurementOrchestrator`를 Flutter에 도입해 화면 상태와 IO 순서를 분리한다.
2. BLE packet fixture를 `rust-core`와 Flutter integration test에서 공유한다.
3. NFC cartridge fixture를 도입해 만료/잔여횟수/alpha 보정 케이스를 테스트한다.
4. `StartSession`, `ProcessMeasurement`, `EndSession`에 idempotency key를 적용한다.
5. fingerprint vector 차원 88/448/896을 API 계약에서 명시하고 validation한다.
6. 측정 실패 유형을 `deviceDisconnected`, `cartridgeInvalid`, `signalLow`, `engineUnavailable`, `serverSyncPending`으로 표준화한다.

완료 기준:
- 실측 장비가 없어도 hardware simulator E2E 통과
- 운영 빌드에서 Rust native 미초기화 시 사용자에게 진단 상태 노출
- 측정 완료 후 Data Hub history에 1분 이내 반영

### 5.3 Data Hub

구현 항목:
1. `DataHubRepository`를 측정 history, biomarker summary, trend, export 상태로 분리한다.
2. 차트는 mock 배열이 아니라 `MeasurementSummary`와 `BiomarkerSummary`에서 파생한다.
3. PDF/CSV/FHIR export 요청은 비동기 job으로 처리하고 다운로드 상태를 추적한다.
4. 가족 데이터는 권한과 동의 상태를 먼저 검증한다.

완료 기준:
- 기간 필터, biomarker 필터, export job E2E 통과
- offline cache에서 stale chart 렌더링 가능

### 5.4 Devices

구현 항목:
1. BLE scan 결과와 서버 등록 device를 분리해서 표시한다.
2. OTA 상태를 queued/downloading/installing/verified/failed로 표준화한다.
3. device ownership, family 공유, subscription entitlement를 DeviceService에서 검증한다.
4. device heartbeat 이벤트를 Kafka와 Redis status cache에 반영한다.

완료 기준:
- 신규 등록, 연결 해제, 상태 갱신, 배터리 읽기, OTA 요청 테스트 통과

### 5.5 AI Coach

구현 항목:
1. 추천 생성 입력을 `latestMeasurements`, `goals`, `riskFlags`, `userProfile`, `consentScope`로 고정한다.
2. 모든 추천에 `reasonCodes`와 `sourceMeasurementIds`를 포함한다.
3. LLM 응답은 의료 면책, 위험도 분류, escalation rule을 통과한 뒤 표시한다.

완료 기준:
- 측정 결과 기반 추천이 mock 없이 생성
- 위험 수치가 있으면 Medical/Emergency 경로로 escalation CTA 표시

### 5.6 Market/Subscription/Payment

구현 항목:
1. 상품, 카트리지, 구독, 결제, entitlement를 하나의 order lifecycle로 묶는다.
2. 카트리지 구매 완료 시 cartridge entitlement가 Devices/Measure에서 확인된다.
3. 결제 webhook idempotency와 환불/취소 이벤트를 audit에 저장한다.

완료 기준:
- 상품 선택 -> 카트 -> 결제 -> 주문 -> entitlement -> 측정 권한 확인 E2E 통과

### 5.7 Medical

구현 항목:
1. telemedicine reservation, video session, prescription, pharmacy transfer를 `MedicalCase`로 묶는다.
2. 건강 데이터 공유 동의 scope를 case 단위로 저장한다.
3. video call 실패 시 재입장/전화 대체/예약 변경 플로우를 제공한다.

완료 기준:
- 예약 -> 데이터 공유 동의 -> 영상 세션 -> 상담 결과 -> 처방 흐름 통과

### 5.8 Family/Emergency

구현 항목:
1. 가족 권한 모델을 owner/guardian/member/viewer로 고정한다.
2. 이상 감지 이벤트가 보호자 알림, emergency rule, audit trail로 전파된다.
3. 미성년/고령 사용자 보호 모드에 별도 consent rule을 적용한다.

완료 기준:
- 가족 초대, 권한 변경, 이상 알림, 보호자 확인 E2E 통과

### 5.9 Settings/Admin

구현 항목:
1. 동의, 보안, 접근성, 데이터 내보내기, 계정 삭제를 규제 증적과 연결한다.
2. Admin은 system stats뿐 아니라 dependency health, audit, queue lag, sync conflict를 보여준다.
3. 운영 설정 변경은 승인 워크플로우와 감사 로그를 필수화한다.

완료 기준:
- Admin RBAC, audit search, config change approval 테스트 통과

---

## 6. 백엔드 구축 계획

### 6.1 Gateway

1. CORS allowlist 적용.
2. JWT middleware와 RBAC middleware를 route group별 적용.
3. `/health/ready`에서 하위 gRPC health, DB, Redis, Kafka, Milvus, ES 상태 확인.
4. request id, user id, device id, session id를 구조화 로그에 포함.
5. REST-gRPC contract test 자동화.

### 6.2 Measurement Service

1. session state machine 도입: created, active, processing, completed, sync_pending, failed, cancelled.
2. `ProcessMeasurement`에서 session active 검증.
3. raw payload hash, correction params, model version, cartridge lot 저장.
4. Timescale hypertable migration과 retention policy 추가.
5. vector store 실패 시 retry queue로 이동하고 UI에는 `vectorSyncPending` 표시.

### 6.3 Device/Cartridge/Calibration

1. device registration token과 cartridge usage entitlement 연결.
2. calibration status가 만료되면 Measure 시작 차단 또는 warning.
3. NFC tag read 결과를 cartridge-service에서 재검증.
4. OTA artifact checksum과 signature 검증.

### 6.4 Event Mesh

필수 topic:

| Topic | Producer | Consumer |
|---|---|---|
| `manpasik.measurement.completed` | measurement-service | data-platform, ai-inference, notification |
| `manpasik.measurement.failed` | measurement-service | notification, admin |
| `manpasik.device.status_changed` | device-service/iot-gateway | devices, admin |
| `manpasik.cartridge.usage_recorded` | cartridge-service | subscription, market |
| `manpasik.payment.completed` | payment-service | subscription, audit |
| `manpasik.health.risk_detected` | ai-inference | emergency, family, medical |
| `manpasik.audit.recorded` | all services | audit-service |

---

## 7. Rust/FFI 구축 계획

1. `RustBridge.init()` 결과를 `nativeReady`, `stubFallback`, `unsupportedPlatform`, `loadFailed`로 세분화한다.
2. 운영 빌드에서 `stubFallback`은 측정 시작을 차단하고 진단 화면으로 연결한다.
3. BLE packet parser를 Rust 단일 구현으로 고정하고 Flutter는 DTO만 사용한다.
4. NFC read/write TODO를 플랫폼별 adapter trait로 분리한다.
5. differential/fingerprint/AI 결과에 engine version과 model version을 포함한다.
6. hardware simulator를 Rust CLI와 Flutter integration test에서 같이 사용한다.

---

## 8. Flutter 구축 계획

### 8.1 상태 모델 표준화

모든 핵심 Provider는 다음 상태를 사용한다.

```dart
sealed class FeatureState<T> {
  const FeatureState();
}

class Loading<T> extends FeatureState<T> {}
class Ready<T> extends FeatureState<T> {
  final T data;
  final DataSource source; // live, cache, demo, offline
  final DateTime loadedAt;
}
class Empty<T> extends FeatureState<T> {}
class Stale<T> extends FeatureState<T> {
  final T data;
  final DateTime lastSuccessAt;
}
class Failed<T> extends FeatureState<T> {
  final Object error;
  final bool retryable;
}
```

### 8.2 Repository 정책

1. Demo repository는 명시적 demo flag에서만 사용한다.
2. 운영 repository가 실패하면 빈 값 반환 금지. `Failed` 또는 `Stale`로 변환한다.
3. 화면은 `source`와 `lastSyncAt`을 작게 표시한다.
4. offline cache는 Repository 내부에서 우선순위를 관리한다.

---

## 9. 데이터/스키마 구축 계획

필수 보강 테이블:

| Table | 목적 |
|---|---|
| `measurement_raw_packets` | 원시 패킷 해시, payload metadata, sequence |
| `measurement_processing_runs` | engine version, model version, correction params |
| `measurement_sync_outbox` | 오프라인/재시도 서버 동기화 |
| `device_heartbeats` | 장비 heartbeat와 상태 이력 |
| `cartridge_usage_events` | 카트리지 사용, 로트, 잔여 횟수 |
| `audit_events` | PHI 접근, 설정 변경, 측정 결과 접근 |
| `traceability_links` | REQ/RISK/TEST/BUILD 연결 |

---

## 10. 테스트 하네스 계획

### 10.1 테스트 피라미드

| 계층 | 테스트 |
|---|---|
| Rust | packet parser, differential, fingerprint, NFC fixture, AI fixture |
| Flutter unit | Provider state, Repository error mapping, route guard |
| Flutter widget | 메뉴별 loading/error/stale/offline 화면 |
| Go unit | service state machine, repository, validation |
| Go integration | PostgreSQL, Redis, Kafka, Milvus, ES testcontainers |
| Contract | OpenAPI/gRPC proto compatibility |
| E2E | Login, device registration, measurement, data hub, market/payment, medical |
| Ops | docker compose smoke, health/readiness, metrics existence |

### 10.2 골든 E2E 시나리오

1. `E2E-MEASURE-001`: 로그인 -> 장비 등록 -> 카트리지 인식 -> 측정 -> 결과 -> 히스토리.
2. `E2E-OFFLINE-001`: 오프라인 측정 -> 로컬 저장 -> 네트워크 복구 -> 서버 동기화.
3. `E2E-RISK-001`: 위험 수치 측정 -> AI risk -> 가족/응급 알림.
4. `E2E-MARKET-001`: 카트리지 구매 -> 결제 -> entitlement -> 측정 권한.
5. `E2E-MEDICAL-001`: 진료 예약 -> 데이터 공유 동의 -> 영상 상담 -> 처방.
6. `E2E-ADMIN-001`: 운영자 로그인 -> 감사 로그 조회 -> 설정 변경 승인.

---

## 11. 운영 구축 계획

### 11.1 Health/Readiness

| Endpoint | 의미 |
|---|---|
| `/health/live` | 프로세스 생존 |
| `/health/ready` | 요청 처리 가능. 필수 의존성 반영 |
| `/health/dependencies` | DB/Redis/Kafka/Milvus/ES/MinIO/gRPC 상세 |
| `/metrics` | Prometheus metrics |

### 11.2 필수 지표

1. `measurement_session_started_total`
2. `measurement_session_failed_total{reason}`
3. `measurement_processing_duration_ms`
4. `ble_disconnect_total`
5. `nfc_read_failed_total`
6. `vector_store_failed_total`
7. `sync_conflict_total`
8. `audit_event_write_failed_total`
9. `gateway_request_duration_ms{route,status}`
10. `dependency_health{dependency}`

---

## 12. 규제/보안 구축 계획

1. 모든 PHI 접근에 audit event 기록.
2. 측정 결과 생성 시 raw packet hash, cartridge lot, calibration params, model version 저장.
3. consent scope가 없는 데이터 공유는 service layer에서 차단.
4. RBAC는 UI guard와 Gateway middleware, service check 3중 적용.
5. release마다 traceability matrix를 생성한다.
6. 위험 제어 테스트를 CI 품질 게이트에 포함한다.

---

## 13. Sprint 실행 계획

### Sprint H0: 하네스 지도화

기간: 1주

작업:
1. route-api-service-data-test 매트릭스 생성.
2. mock/stub/demo/fallback register 작성.
3. 핵심 메뉴별 completion score 산정.
4. 운영 차단 항목과 개발 허용 항목 분리.

산출물:
- `docs/audit/harness-capillary-map.md`
- `docs/audit/mock-retirement-register.md`
- `docs/audit/menu-completion-scorecard.md`

### Sprint H1: 측정 골든 패스

기간: 2주

작업:
1. MeasurementOrchestrator 구현.
2. Rust packet/NFC fixtures 구현.
3. Go measurement session state machine 구현.
4. measurement raw/process/audit schema 추가.
5. Flutter Measure E2E 작성.

품질 게이트:
- Rust tests pass
- Go measurement tests pass
- Flutter measure integration pass
- Docker compose smoke pass

### Sprint H2: 실데이터 메뉴 결선

기간: 4주

작업:
1. Home/Data Hub/Devices real provider 전환.
2. Market/Subscription/Payment entitlement 결선.
3. Medical case lifecycle 결선.
4. Family/Emergency event escalation 결선.
5. Admin observability/audit 결선.

품질 게이트:
- mock register P0 항목 80% 은퇴
- 핵심 메뉴 API contract 90% 통과

### Sprint H3: 오프라인/동기화

기간: 3주

작업:
1. Flutter local outbox/inbox 구현.
2. 서버 sync endpoint와 conflict payload 구현.
3. conflict resolver 실제 payload 연결.
4. offline measurement E2E 구현.

품질 게이트:
- offline measurement pass
- sync replay idempotency pass
- conflict resolution pass

### Sprint H4: 운영/규제 하네스

기간: 3주

작업:
1. dependency-aware readiness 구현.
2. Prometheus/Grafana dashboard 보강.
3. audit event 표준화.
4. traceability matrix generator 구현.
5. release evidence bundle 생성.

품질 게이트:
- readiness dependency failure test pass
- audit write coverage pass
- traceability link coverage 90% 이상

### Sprint H5: 릴리스 안정화

기간: 2주

작업:
1. smoke/regression/performance suite 고정.
2. chaos test: DB down, Kafka down, BLE disconnect, FFI unavailable.
3. security scan, dependency review.
4. release runbook 작성.

품질 게이트:
- P0/P1 결함 0
- 핵심 E2E 100% pass
- 운영 runbook review complete

---

## 14. 즉시 착수 우선순위

1. `docs/audit/mock-retirement-register.md` 작성 및 P0 mock 분류.
2. `MeasurementOrchestrator` 설계 및 skeleton 구현.
3. measurement-service session state machine 테스트 먼저 작성.
4. Gateway readiness를 실제 의존성 기반으로 보강.
5. Flutter Provider의 silent catch 패턴을 `Failed/Stale` 상태로 전환.

---

## 15. Definition of Done

완성된 플랫폼으로 판정하려면 다음을 모두 만족해야 한다.

1. 핵심 메뉴 9개(Home, Measure, Data Hub, Devices, AI Coach, Market, Medical, Family, Admin)의 completion score 85점 이상.
2. Measure, Auth, Payment, Medical data sharing, Emergency는 92점 이상.
3. 운영 빌드에서 silent mock/stub fallback 0건.
4. 모든 핵심 API가 contract test를 통과.
5. Docker compose smoke가 clean 환경에서 1회 명령으로 통과.
6. 오프라인 측정과 복구 동기화 E2E 통과.
7. PHI 접근 audit write coverage 100%.
8. REQ/RISK/TEST/BUILD traceability coverage 90% 이상.
9. 장애 주입 테스트에서 사용자에게 복구 가능 상태가 노출.
10. 릴리스 evidence bundle 생성.

---

## 16. 운영 원칙

만파식은 의료기기와 헬스케어 SaaS가 결합된 플랫폼이다. 따라서 "화면이 그럴듯하게 보이는 것"은 완료가 아니다. 완료는 실제 장비, 실제 데이터, 실제 장애, 실제 감사, 실제 복구가 하나의 흐름으로 동작하는 상태다.

이 계획은 그 상태를 만들기 위한 실행 기준이다.
