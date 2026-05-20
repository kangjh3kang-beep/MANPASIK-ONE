# Dart Proto Official Generated Replacement 게이트

**작성일**: 2026-05-01  
**작성자**: Codex  
**범위**: Flutter 앱 `grpc/protobuf` lock 정렬 및 checked-in official generated Dart 교체

## 상세 구현계획

1. Flutter 앱 `pubspec.yaml`의 gRPC/protobuf 계열을 공식 생성기와 맞춘다.
2. `grpc ^5.1.0`, `protobuf ^6.0.0`, `fixnum ^1.1.1`을 앱 direct dependency로 고정한다.
3. `flutter pub get`으로 실제 `pubspec.lock`를 갱신한다.
4. `backend/shared/proto/manpasik.proto`를 `protoc --dart_out=grpc`로 `frontend/flutter-app/lib/generated`에 공식 생성한다.
5. checked-in generated 파일은 `manpasik.pb.dart`, `manpasik.pbgrpc.dart`, `manpasik.pbenum.dart`, `manpasik.pbjson.dart` 4개로 확장한다.
6. 공식 generated API 차이로 드러난 소비자 코드를 좁게 정렬한다.
7. preflight가 `ready_for_compile_gate`로 바뀌었는지 확인한다.
8. wire golden, native repository stream, REST/history, orchestrator, Go measurement service 회귀를 실행한다.

## 구현 결과

- Flutter 앱 lock이 `grpc 5.1.0`, `protobuf 6.0.0` 계열로 정렬됐다.
- `fixnum`은 generated 파일이 직접 import하므로 direct dependency로 승격했다.
- 수동 generated stub을 공식 `protoc_plugin 25.0.0` 산출물로 교체했다.
- 2026-05-01 후속 Auth `user_id` 계약 보강 후 현재 official generated Dart 산출물은 64,750라인이다.
- 기존 `check_proto_generation_preflight.sh` 판정이 `PROTO_PREFLIGHT_FULL_REPLACEMENT=ready_for_compile_gate`로 전환됐다.
- `MeasurementData` wire golden은 공식 generated output에서도 동일하게 유지된다.

## 소비자 정렬

- `MeasurementRepositoryImpl`
  - 중복 import 제거.
  - official `MeasurementServiceClient.streamMeasurement` surface 유지.
- `AuthRepositoryImpl`
  - 2026-05-01 후속 Auth 계약 보강으로 `LoginResponse.userId`를 우선 사용한다.
  - legacy 빈 응답에 한해 fallback 성격의 `unknown`을 유지한다.
- `DeviceRepositoryImpl`
  - `DeviceInfo.status`가 `int`가 아니라 `DeviceStatus` enum으로 노출되므로 enum switch로 정렬.
- `UserRepositoryImpl`
  - `SubscriptionTier` enum을 도메인 `int` DTO로 넘길 때 `.value` 사용.
- `ChatNotifier`
  - official client 이름인 `AiInferenceServiceClient`로 정렬.
- `AdminSettingsProvider`
  - official pbgrpc export에 맞춰 중복 pb import 제거.
- `measurement_stream_wire_contract_test`
  - 테스트 이름을 official Dart proto 상태에 맞게 갱신.

## 변경 파일

- `frontend/flutter-app/pubspec.yaml`
- `frontend/flutter-app/pubspec.lock`
- `frontend/flutter-app/lib/generated/manpasik.pb.dart`
- `frontend/flutter-app/lib/generated/manpasik.pbgrpc.dart`
- `frontend/flutter-app/lib/generated/manpasik.pbenum.dart`
- `frontend/flutter-app/lib/generated/manpasik.pbjson.dart`
- `backend/shared/proto/manpasik.proto`
- `backend/shared/gen/go/v1/manpasik.pb.go`
- `backend/shared/gen/go/v1/manpasik_grpc.pb.go`
- `frontend/flutter-app/lib/features/measurement/data/measurement_repository_impl.dart`
- `frontend/flutter-app/lib/features/auth/data/auth_repository_impl.dart`
- `frontend/flutter-app/lib/features/devices/data/device_repository_impl.dart`
- `frontend/flutter-app/lib/features/user/data/user_repository_impl.dart`
- `frontend/flutter-app/lib/shared/providers/admin_settings_provider.dart`
- `frontend/flutter-app/lib/shared/providers/chat_provider.dart`
- `frontend/flutter-app/test/features/measurement/data/measurement_stream_wire_contract_test.dart`

## 품질 게이트

```bash
cd /home/kangjh3kang/Manpasik/frontend/flutter-app
export PATH=/mnt/d/우리집/flutter_cache/flutter/bin:/usr/local/bin:/usr/bin:/bin
flutter pub get
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
  lib/generated/manpasik.pbenum.dart \
  lib/generated/manpasik.pbjson.dart \
  lib/features/measurement/data/measurement_repository_impl.dart \
  lib/features/auth/data/auth_repository_impl.dart \
  lib/features/user/data/user_repository_impl.dart \
  lib/shared/providers/chat_provider.dart \
  test/features/measurement/data/measurement_stream_wire_contract_test.dart \
  test/features/measurement/data/measurement_repository_impl_test.dart

cd /home/kangjh3kang/Manpasik
/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 \
  ./backend/services/measurement-service/cmd \
  ./backend/services/measurement-service/internal/handler \
  ./backend/services/measurement-service/internal/service
```

결과:

- `flutter pub get`: PASS
- `PROTO_PREFLIGHT_STATUS=generated_streammeasurement_available`
- `PROTO_PREFLIGHT_FULL_REPLACEMENT=ready_for_compile_gate`
- `PROTO_COMPILE_GATE_STATUS=passed`
- `PROTO_COMPILE_GATE_GENERATED_DART_LINES=64750`
- `PROTO_COMPILE_GATE_GRPC_VERSION=5.1.0`
- `PROTO_COMPILE_GATE_PROTOBUF_VERSION=6.0.0`
- Flutter stream/repository/orchestrator tests: PASS
- Flutter auth proto contract test: PASS
- Flutter generated/consumer analyze target: PASS
- Go auth service/handler regression: PASS
- Go measurement cmd/handler/service regression: PASS

## 잔여 리스크

- 2026-05-01 후속 보강으로 `LoginResponse.user_id = 5`가 추가되어 Auth 도메인의 정상 로그인 `unknown` fallback 리스크는 해소됐다.
- Gateway REST 로그인 응답이 별도 DTO를 쓰는 경로는 추가 점검이 필요하다.
- Docker Compose container smoke는 현재 WSL Docker integration 부재로 runtime blocked 상태다.

## 다음 단계

- Docker Desktop WSL integration을 활성화한 뒤 compose container smoke를 실제 PASS로 확정한다.
- Gateway auth REST/E2E 응답의 사용자 식별자 전달 상태를 점검한다.
- 이후 Gateway `/measurements/process`와 native direct stream 중 운영 모드별 라우팅 정책을 문서화한다.
