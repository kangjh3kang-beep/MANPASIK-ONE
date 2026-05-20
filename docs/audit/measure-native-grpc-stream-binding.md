# Measure Native gRPC Stream Binding 게이트

**작성일**: 2026-05-01  
**작성자**: Codex  
**범위**: MR-009 native process REST bridge 축소 및 `StreamMeasurement` 직결

## 목적

이전 단계에서 네이티브 Flutter `processMeasurement()`는 local echo fallback을 제거하고 Gateway `POST /measurements/process`를 통해 backend storage/vector 경로에 도달했다. 이번 단계는 그 임시 REST bridge를 줄이고, native repository가 `MeasurementService.StreamMeasurement` bidi stream을 직접 호출하도록 보강한다.

## 구현 결과

- 수동 Dart gRPC 스텁에 `MeasurementData`, `DifferentialCorrection`, `EnvironmentMeta`, `MeasurementResult` 메시지를 추가했다.
- `MeasurementServiceClient.streamMeasurement()` client method를 추가했다.
- `MeasurementRepositoryImpl.processMeasurement()`가 REST Gateway process endpoint 대신 gRPC stream으로 단일 measurement frame을 전송한다.
- `AuthInterceptor`가 붙은 기존 `MeasurementServiceClient` 경로를 사용하므로 native gRPC 인증 경로가 유지된다.
- Web/REST repository는 기존 Gateway process mapper를 그대로 사용한다.

## 변경 파일

- `frontend/flutter-app/lib/generated/manpasik.pb.dart`
  - `StreamMeasurement`에 필요한 수동 protobuf 메시지 클래스를 추가했다.
- `frontend/flutter-app/lib/generated/manpasik.pbgrpc.dart`
  - `streamMeasurement(Stream<MeasurementData>)` client method를 추가했다.
- `frontend/flutter-app/lib/features/measurement/data/measurement_repository_impl.dart`
  - native `processMeasurement()`를 gRPC stream 직결로 전환했다.
  - 테스트 주입용 `MeasurementStreamCall`을 추가해 네트워크 없이 stream frame 계약을 검증한다.
- `frontend/flutter-app/lib/core/providers/grpc_provider.dart`
  - native measurement repository에서 REST process client 주입을 제거했다.
- `frontend/flutter-app/test/features/measurement/data/measurement_repository_impl_test.dart`
  - native process path가 `MeasurementData` frame을 gRPC stream으로 전송하고 `MeasurementResult`를 도메인 result로 매핑하는지 검증한다.

## 제약 및 결정

- 전체 `protoc --dart_out=grpc` 재생성은 `/tmp` 검증 기준 약 6.4만 줄 산출물과 `pbenum/pbjson` 추가 파일을 만들었다. 기존 Flutter 앱은 수동 축약 스텁에 맞춰져 있어 이번 단계에서는 필요한 stream surface만 좁게 보강했다.
- 현재 lock된 `protobuf 3.1.0` 환경에서는 최신 protoc plugin이 생성한 well-known Timestamp import 경로가 바로 컴파일되지 않았다. 따라서 native direct stream result의 `processedAt`은 기존 수동 스텁 한계에 맞춰 `null`로 유지한다.

## 품질 게이트

```bash
cd /home/kangjh3kang/Manpasik/frontend/flutter-app
export PATH="/mnt/d/우리집/flutter_cache/flutter/bin:/usr/local/bin:/usr/bin:/bin"
dart format \
  lib/generated/manpasik.pb.dart \
  lib/generated/manpasik.pbgrpc.dart \
  lib/features/measurement/data/measurement_repository_impl.dart \
  lib/core/providers/grpc_provider.dart \
  test/features/measurement/data/measurement_repository_impl_test.dart
flutter test --no-pub \
  test/features/measurement/data/measurement_repository_impl_test.dart \
  test/features/measurement/data/measurement_repository_rest_test.dart \
  test/features/measurement/application/measurement_golden_path_orchestrator_test.dart
flutter analyze --no-pub \
  lib/generated/manpasik.pb.dart \
  lib/generated/manpasik.pbgrpc.dart \
  lib/features/measurement/data/measurement_repository_impl.dart \
  lib/core/providers/grpc_provider.dart \
  test/features/measurement/data/measurement_repository_impl_test.dart
```

결과: PASS

추가 회귀 게이트:

- `flutter test --no-pub test/core/providers/measurement_history_provider_test.dart test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_repository_rest_test.dart test/features/measurement/data/measurement_trace_sink_rest_test.dart test/features/measurement/application/measurement_golden_path_orchestrator_test.dart test/features/measurement/domain/measurement_domain_test.dart test/features/domain_models_test.dart`: PASS
- `flutter analyze --no-pub lib/generated/manpasik.pb.dart lib/generated/manpasik.pbgrpc.dart lib/features/measurement/data/measurement_repository_impl.dart lib/core/providers/grpc_provider.dart lib/features/measurement/domain/measurement_repository.dart lib/features/measurement/data/measurement_repository_rest.dart lib/shared/providers/ecosystem_providers.dart lib/features/measurement/presentation/measurement_result_screen.dart lib/features/home/presentation/home_screen.dart test/core/providers/measurement_history_provider_test.dart test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_repository_rest_test.dart`: PASS

## 다음 단계

- Dart/Go wire smoke gate는 `docs/audit/measure-stream-wire-smoke-gate.md`에서 추가 완료했다.
- Dart proto generation preflight는 `docs/audit/dart-proto-generation-preflight.md`에서 추가 완료했다.
- 다음 단계는 `protobuf`, `grpc`, `protoc_plugin` 버전 조합을 정렬한 뒤 전체 Dart generated proto compile gate를 실행하는 것이다.
