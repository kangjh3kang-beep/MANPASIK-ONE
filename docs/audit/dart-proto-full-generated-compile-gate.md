# Dart Proto Full Generated Compile 게이트

**작성일**: 2026-05-01  
**작성자**: Codex  
**범위**: 공식 `protoc --dart_out=grpc` 전체 산출물의 격리 compile 검증

## 상세 구현계획

1. checked-in `frontend/flutter-app/lib/generated` 수동 스텁은 건드리지 않는다.
2. `backend/shared/proto/manpasik.proto`를 임시 Dart package의 `lib/generated`로 전체 생성한다.
3. 생성 산출물에 `MeasurementData`와 `MeasurementServiceClient.streamMeasurement`가 있는지 확인한다.
4. 임시 package의 toolchain은 최신 `protoc_plugin 25.0.0` 산출물과 맞는 `grpc ^5.1.0`, `protobuf ^6.0.0`, `fixnum ^1.1.1`로 고정한다.
5. 임시 package에서 `dart pub get`, `dart analyze bin lib`, `dart run bin/compile_smoke.dart`를 순서대로 실행한다.
6. smoke entrypoint는 `MeasurementData`, `DifferentialCorrection`, `EnvironmentMeta`, `MeasurementResult`, `MeasurementServiceClient`를 실제로 참조한다.
7. 성공 시 생성 라인 수와 resolved `grpc`/`protobuf` 버전을 출력한다.

## 구현 결과

- `frontend/flutter-app/scripts/check_proto_generation_compile_gate.sh`를 추가했다.
- 공식 generated Dart 전체 산출물 64,736라인이 격리 package에서 컴파일 가능함을 확인했다.
- resolved toolchain은 `grpc 5.1.0`, `protobuf 6.0.0`이다.
- 기존 `check_proto_generation_preflight.sh`는 현재 앱 lock 기준의 위험을 계속 보고한다.
- 따라서 판정은 다음과 같다.
  - 공식 generated output compile 가능성: 확인 완료
  - 현재 앱 checked-in stub 즉시 치환: 아직 보류
  - 보류 사유: 앱 `pubspec.lock`은 `grpc 4.1.0`, `protobuf 3.1.0`이고, 현재 `flutter pub` resolution은 이 WSL Flutter SDK 환경에서 정상 수행되지 않았다.

## 변경 파일

- `frontend/flutter-app/scripts/check_proto_generation_compile_gate.sh`
  - 임시 Dart package 생성.
  - `manpasik.proto` 전체 Dart gRPC 생성.
  - `grpc ^5.1.0`, `protobuf ^6.0.0` compile/analyze/run smoke.
- `docs/audit/dart-proto-full-generated-compile-gate.md`
  - 상세 구현계획, 실행 결과, 품질 게이트 문서화.
- `docs/audit/dart-proto-generation-preflight.md`
  - preflight 이후 compile gate 완료 상태 반영.
- `docs/audit/measure-stream-wire-smoke-gate.md`
  - 전체 generated compile gate 완료 상태 반영.
- `docs/audit/mock-retirement-register.md`
  - MR-009 잔여 조치를 앱 의존성 정렬과 checked-in generated 교체 단계로 좁힘.

## 품질 게이트

```bash
cd /home/kangjh3kang/Manpasik/frontend/flutter-app
export PATH=/mnt/d/우리집/flutter_cache/flutter/bin:/usr/local/bin:/usr/bin:/bin
bash -n scripts/check_proto_generation_compile_gate.sh
scripts/check_proto_generation_preflight.sh
scripts/check_proto_generation_compile_gate.sh
flutter test --no-pub \
  test/features/measurement/data/measurement_stream_wire_contract_test.dart \
  test/features/measurement/data/measurement_repository_impl_test.dart \
  test/features/measurement/data/measurement_repository_rest_test.dart \
  test/features/measurement/application/measurement_golden_path_orchestrator_test.dart
flutter analyze --no-pub \
  lib/generated/manpasik.pb.dart \
  lib/generated/manpasik.pbgrpc.dart \
  lib/features/measurement/data/measurement_repository_impl.dart \
  test/features/measurement/data/measurement_stream_wire_contract_test.dart \
  test/features/measurement/data/measurement_repository_impl_test.dart

cd /home/kangjh3kang/Manpasik
/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 \
  ./backend/services/measurement-service/cmd \
  ./backend/services/measurement-service/internal/handler \
  ./backend/services/measurement-service/internal/service
```

결과:

- `PROTO_PREFLIGHT_STATUS=generated_streammeasurement_available`
- `PROTO_PREFLIGHT_FULL_REPLACEMENT=blocked_timestamp_import_incompatible_with_current_pub_cache`
- `PROTO_COMPILE_GATE_STATUS=passed`
- `PROTO_COMPILE_GATE_GENERATED_DART_LINES=64736`
- `PROTO_COMPILE_GATE_GRPC_VERSION=5.1.0`
- `PROTO_COMPILE_GATE_PROTOBUF_VERSION=6.0.0`
- Flutter stream subset test/analyze: PASS
- Go measurement cmd/handler/service regression: PASS

## 다음 단계

- Flutter 앱 `pubspec.yaml`/`pubspec.lock` 정렬과 checked-in official generated 교체는 `docs/audit/dart-proto-official-generated-replacement.md`에서 완료했다.
- 다음 단계는 Auth login 사용자 식별자 보강 또는 Docker Compose container smoke runtime PASS 확보다.
