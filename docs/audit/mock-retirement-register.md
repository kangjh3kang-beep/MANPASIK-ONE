# ManPaSik Mock Retirement Register

**작성일**: 2026-05-01  
**작성자**: Codex  
**목적**: demo/mock/stub/fallback 의존 지점을 운영 차단 항목과 개발 허용 항목으로 분류한다.

| ID | 위치 | 유형 | 위험도 | 운영 허용 | 은퇴 전략 |
|---|---|---|---|---|---|
| MR-001 | `frontend/flutter-app/lib/core/providers/grpc_provider.dart` | Demo Mode mock data | P1 | 조건부 | 명시적 DemoMode에서만 사용, RealMode 실패는 Failed/Stale로 노출. 2026-05-01 history stale/error 상태 노출 보강 완료 |
| MR-002 | `frontend/flutter-app/lib/core/services/rust_ffi_stub.dart` | Rust FFI stub fallback | P0 | 불가 | 운영 빌드에서 native 미초기화 시 측정 시작 차단 및 진단 노출 |
| MR-003 | `frontend/flutter-app/lib/features/devices/data/device_repository_impl.dart` | Native demo device | P1 | 불가 | BLE scan 결과와 서버 등록 장비를 분리 |
| MR-004 | `backend/services/*/cmd/main.go` | memory repository fallback | P1 | 불가 | 운영 환경에서는 readiness fail, 개발 환경에서만 허용 |
| MR-005 | `backend/services/telemedicine-service` tests | mock WebRTC provider | P2 | 테스트만 | provider contract test로 보완 |
| MR-006 | `backend/services/ai-inference-service` tests | mock LLM client | P2 | 테스트만 | LLM adapter contract/fault test 유지 |
| MR-007 | `frontend/flutter-app/lib/features/family/data/family_repository_rest.dart` | placeholder methods | P1 | 불가 | Family API 결선 완료 후 제거 |
| MR-008 | Market mock images/data | sample asset/data | P2 | 데모만 | 상품 API 응답과 asset fallback 분리 |
| MR-009 | `frontend/flutter-app/lib/features/measurement/data/measurement_repository_impl.dart` | Native process local echo fallback | P0 | 퇴역 완료 | 2026-05-01 local echo 제거 후 Gateway `/measurements/process` 경유. 2026-05-01 native `StreamMeasurement` gRPC 직결, helper-process cross-process smoke, 실제 service binary smoke 완료. Container smoke 하네스 구축 완료, 현재 WSL Docker Compose 연결 부재로 실행 blocked. 공식 Dart generated compile gate 및 checked-in 교체 완료 |

## P0 은퇴 조건

1. 측정 시작 경로에서 silent stub fallback 제거.
2. `MeasurementGoldenPathOrchestrator`가 engine availability를 판정.
3. 운영 빌드에서 `RustBridge.isNativeEnabled == false`이면 측정 대신 진단 상태를 노출.

## 2026-05-01 Measure 재대조 결과

| 항목 | 판정 | 증거 | 잔여 조치 |
|---|---|---|---|
| MR-002 | 운영 차단 조건 충족 | `MeasurementGoldenPathOrchestrator` readiness gate가 `RustBridgeDiagnostics.canRunMeasurement == false`이면 session start 전에 실패 처리. 테스트가 session/process/end 미호출을 검증 | 모바일 native Rust 라이브러리 패키징 후 release smoke test |
| MR-009 | 퇴역 완료, native gRPC 직결 완료 | `MeasurementRepositoryImpl.processMeasurement()`가 더 이상 Rust 처리 결과를 로컬 echo로 반환하지 않고 `MeasurementService.StreamMeasurement` gRPC stream으로 `MeasurementData` frame을 전송. Dart/Go golden wire smoke, in-process transport smoke, TCP loopback multi-frame smoke, helper-process cross-process smoke, 실제 service binary smoke가 서버 handler/main wiring 저장/응답 계약을 검증. Compose container smoke 하네스는 추가됐고 Docker 연결 환경에서 실행 가능. 공식 Dart generated 64,750라인 compile gate가 `grpc 5.1.0`, `protobuf 6.0.0`에서 통과했고 checked-in generated file 교체 완료 | Docker WSL integration 활성화 후 compose container smoke PASS 확보 |
| MR-001 | 조건부 유지, RealMode 보강 완료 | DemoMode mock history/device data는 명시적 `authState.isDemo`에서만 사용. RealMode history 실패는 `MeasurementHistoryResult.isStale/errorMessage`로 노출하고 UI가 동기화 문제를 표시 | DemoMode 표시/권한 정책을 release flavor에서 재검증 |

### 품질 게이트

- `flutter test --no-pub test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_repository_rest_test.dart test/features/measurement/data/measurement_trace_sink_rest_test.dart test/features/measurement/application/measurement_golden_path_orchestrator_test.dart test/core/network/tenant_interceptor_test.dart`: PASS
- `flutter analyze --no-pub lib/core/network/tenant_interceptor.dart lib/core/providers/grpc_provider.dart lib/features/measurement/data/measurement_process_gateway_mapper.dart lib/features/measurement/data/measurement_repository_impl.dart lib/features/measurement/data/measurement_repository_rest.dart test/features/measurement/data/measurement_repository_impl_test.dart`: PASS
- `flutter test --no-pub test/core/providers/measurement_history_provider_test.dart test/features/measurement/data/measurement_repository_rest_test.dart test/features/measurement/domain/measurement_domain_test.dart test/features/domain_models_test.dart test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_trace_sink_rest_test.dart test/features/measurement/application/measurement_golden_path_orchestrator_test.dart`: PASS
- `flutter analyze --no-pub lib/features/measurement/domain/measurement_repository.dart lib/features/measurement/data/measurement_repository_rest.dart lib/core/providers/grpc_provider.dart lib/shared/providers/ecosystem_providers.dart lib/features/measurement/presentation/measurement_result_screen.dart lib/features/home/presentation/home_screen.dart test/core/providers/measurement_history_provider_test.dart test/features/measurement/data/measurement_repository_rest_test.dart`: PASS
- `flutter test --no-pub test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_repository_rest_test.dart test/features/measurement/application/measurement_golden_path_orchestrator_test.dart`: PASS
- `flutter analyze --no-pub lib/generated/manpasik.pb.dart lib/generated/manpasik.pbgrpc.dart lib/features/measurement/data/measurement_repository_impl.dart lib/core/providers/grpc_provider.dart test/features/measurement/data/measurement_repository_impl_test.dart`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -v ./backend/services/measurement-service/internal/handler -run TestStreamMeasurementAgainstHelperProcess`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./backend/services/measurement-service/internal/handler ./backend/services/measurement-service/internal/service`: PASS
- `bash scripts/check_proto_generation_preflight.sh`: PASS, `StreamMeasurement` 생성 가능, 전체 치환은 locked `protobuf 3.1.0` Timestamp import 호환성으로 보류
- `flutter test --no-pub test/features/measurement/data/measurement_stream_wire_contract_test.dart test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_repository_rest_test.dart test/features/measurement/application/measurement_golden_path_orchestrator_test.dart`: PASS
- `flutter analyze --no-pub lib/generated/manpasik.pb.dart lib/generated/manpasik.pbgrpc.dart lib/features/measurement/data/measurement_repository_impl.dart test/features/measurement/data/measurement_stream_wire_contract_test.dart test/features/measurement/data/measurement_repository_impl_test.dart`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -v ./backend/services/measurement-service/cmd -run TestMeasurementServiceBinarySmoke`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./backend/services/measurement-service/cmd ./backend/services/measurement-service/internal/handler ./backend/services/measurement-service/internal/service`: PASS
- `bash -n scripts/measure_service_compose_smoke.sh`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -v ./backend/services/measurement-service/cmd -run 'TestMeasurementService(BinarySmoke|ExternalEndpointSmoke)'`: PASS, external endpoint smoke는 env 미지정 시 SKIP
- `scripts/measure_service_compose_smoke.sh`: BLOCKED, `MEASUREMENT_COMPOSE_SMOKE_STATUS=blocked_docker_compose_unavailable`
- `scripts/check_proto_generation_compile_gate.sh`: PASS, `PROTO_COMPILE_GATE_STATUS=passed`, `grpc 5.1.0`, `protobuf 6.0.0`, generated Dart 64,736 lines
- `flutter pub get`: PASS, app lock aligned to `grpc 5.1.0`, `protobuf 6.0.0`, `fixnum 1.1.1`
- `protoc --dart_out=grpc:lib/generated`: PASS, official checked-in generated files replaced/expanded to `pb`, `pbgrpc`, `pbenum`, `pbjson`
- `scripts/check_proto_generation_preflight.sh`: PASS, `PROTO_PREFLIGHT_FULL_REPLACEMENT=ready_for_compile_gate`
