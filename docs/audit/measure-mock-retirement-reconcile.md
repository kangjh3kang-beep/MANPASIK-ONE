# Measure Mock Retirement 재대조

**작성일**: 2026-05-01  
**작성자**: Codex  
**범위**: Measure native repository process path, MR-002/MR-009 재대조

## 목적

이전 단계까지 Measure golden path는 Flutter orchestrator, Rust readiness gate, Gateway process route, measurement-service storage/vector persistence, remote trace audit persistence까지 연결됐다. 이번 단계는 `docs/audit/mock-retirement-register.md`의 Measure 관련 항목을 실제 코드와 대조하고, 네이티브 Flutter 경로에 남아 있던 로컬 echo fallback을 제거한다.

## 확인 결과

- `MR-002` Rust FFI stub fallback은 readiness gate로 운영 차단 조건을 충족한다.
  - `RustBridgeDiagnostics.canRunMeasurement == false`이면 `MeasurementGoldenPathOrchestrator`가 session start 전에 중단한다.
  - 실패 trace event가 남고, server/process/end 호출은 일어나지 않는다.
- 신규 `MR-009`로 네이티브 `processMeasurement` local echo fallback을 등록하고 바로 퇴역 처리했다.
  - 기존 `MeasurementRepositoryImpl.processMeasurement()`는 checked-in Dart gRPC binding에 `StreamMeasurement`가 없다는 이유로 Rust 결과를 그대로 반환했다.
  - 이제 동일 payload mapper를 통해 Gateway `POST /api/v1/measurements/process`로 전송한다.
  - 따라서 네이티브 Measure도 backend StreamMeasurement, measurement storage, fingerprint vector storage 경로를 탄다.

## 변경 파일

- `frontend/flutter-app/lib/features/measurement/data/measurement_process_gateway_mapper.dart`
  - Gateway process payload encode/decode 공통 mapper 추가.
- `frontend/flutter-app/lib/features/measurement/data/measurement_repository_rest.dart`
  - Web/REST repository가 공통 mapper를 사용하도록 정리.
- `frontend/flutter-app/lib/features/measurement/data/measurement_repository_impl.dart`
  - Native repository의 process path를 local echo에서 Gateway REST process 호출로 교체.
  - `accessTokenProvider` 값을 REST Authorization header에 반영.
- `frontend/flutter-app/lib/core/providers/grpc_provider.dart`
  - Native `MeasurementRepositoryImpl`에 shared `ManPaSikRestClient`를 주입.
- `frontend/flutter-app/lib/core/network/tenant_interceptor.dart`
  - SharedPreferences 접근 실패가 REST 요청 전체를 막지 않도록 fail-open 처리.
- `frontend/flutter-app/test/features/measurement/data/measurement_repository_impl_test.dart`
  - Native process path가 Gateway endpoint로 payload와 Authorization을 전송하는지 검증.
- `docs/audit/mock-retirement-register.md`
  - MR-009 등록 및 Measure 재대조 결과 추가.

## 품질 게이트

```bash
cd /home/kangjh3kang/Manpasik/frontend/flutter-app
export PATH="/mnt/d/우리집/flutter_cache/flutter/bin:/usr/local/bin:/usr/bin:/bin"
dart format lib/core/network/tenant_interceptor.dart \
  lib/core/providers/grpc_provider.dart \
  lib/features/measurement/data/measurement_process_gateway_mapper.dart \
  lib/features/measurement/data/measurement_repository_impl.dart \
  lib/features/measurement/data/measurement_repository_rest.dart \
  test/features/measurement/data/measurement_repository_impl_test.dart
flutter test --no-pub \
  test/features/measurement/data/measurement_repository_impl_test.dart \
  test/features/measurement/data/measurement_repository_rest_test.dart \
  test/features/measurement/data/measurement_trace_sink_rest_test.dart \
  test/features/measurement/application/measurement_golden_path_orchestrator_test.dart \
  test/core/network/tenant_interceptor_test.dart
flutter analyze --no-pub \
  lib/core/network/tenant_interceptor.dart \
  lib/core/providers/grpc_provider.dart \
  lib/features/measurement/data/measurement_process_gateway_mapper.dart \
  lib/features/measurement/data/measurement_repository_impl.dart \
  lib/features/measurement/data/measurement_repository_rest.dart \
  test/features/measurement/data/measurement_repository_impl_test.dart
```

결과: PASS

## 후속 진행

- MR-001 RealMode history 실패는 `docs/audit/measure-history-stale-state.md`에서 stale/error 상태 노출로 보강 완료했다.
- MR-009 native process REST bridge는 `docs/audit/measure-native-grpc-stream-binding.md`에서 `StreamMeasurement` gRPC 직결로 치환 완료했다.
- 다음 단계는 Dart protobuf/protoc plugin 버전을 정렬해 전체 generated Dart를 공식 산출물로 전환할 수 있는지 검증하고, 실제 measurement-service integration smoke test를 붙이는 것이다.
