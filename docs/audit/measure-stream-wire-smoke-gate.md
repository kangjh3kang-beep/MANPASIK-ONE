# Measure Stream Wire Smoke 게이트

**작성일**: 2026-05-01  
**작성자**: Codex  
**범위**: Dart native stream surface와 Go measurement-service wire contract 검증

## 상세 구현계획

1. Dart 수동 protobuf 스텁이 생성하는 `MeasurementData` binary payload를 golden hex로 고정한다.
2. golden payload에 `session_id`, `raw_channels`, `DifferentialCorrection`, `EnvironmentMeta`를 포함해 Measure 핵심 입력면을 커버한다.
3. Go measurement-service handler test에서 같은 golden hex를 `v1.MeasurementData`로 unmarshal한다.
4. unmarshal된 frame을 실제 `MeasurementHandler.StreamMeasurement()` fake stream에 넣어 저장소 저장, fingerprint vector 저장, stream response까지 검증한다.
5. 전체 Dart generated proto 치환은 별도 preflight script로 분리해, 현재 lock된 protobuf 환경에서 덮어쓰기 리스크를 명확히 보고한다.

## 구현 결과

- Dart와 Go가 같은 `MeasurementData` wire payload를 공유하는 smoke gate를 추가했다.
- Dart test는 수동 스텁의 `writeToBuffer()` 결과가 Go proto golden hex와 일치하는지 검증한다.
- Go test는 Dart golden hex를 unmarshal하고 실제 measurement-service stream handler를 통과시킨다.
- proto generation preflight script를 추가해 전체 generated Dart 전환 가능성과 차단 사유를 자동 보고한다.

## 변경 파일

- `frontend/flutter-app/test/features/measurement/data/measurement_stream_wire_contract_test.dart`
  - Dart 수동 `MeasurementData` 스텁의 binary wire contract를 golden hex로 고정한다.
- `backend/services/measurement-service/internal/handler/grpc_stream_wire_contract_test.go`
  - 같은 golden hex를 Go generated proto로 읽고 실제 `StreamMeasurement` handler에 흘려보낸다.
- `frontend/flutter-app/scripts/check_proto_generation_preflight.sh`
  - 임시 디렉터리에 Dart gRPC 산출물을 생성하고 `streamMeasurement` surface 및 Timestamp 호환성을 점검한다.
- `docs/audit/dart-proto-generation-preflight.md`
  - 전체 generated Dart 전환 전 preflight 판정과 보류 사유를 문서화한다.

## 품질 게이트

```bash
cd /home/kangjh3kang/Manpasik
export PATH="/home/kangjh3kang/sdk/go-go1.26.2/bin:$PATH"
go test ./backend/services/measurement-service/internal/handler ./backend/services/measurement-service/internal/service

cd /home/kangjh3kang/Manpasik/frontend/flutter-app
export PATH="/mnt/d/우리집/flutter_cache/flutter/bin:/usr/local/bin:/usr/bin:/bin"
bash scripts/check_proto_generation_preflight.sh
flutter test --no-pub \
  test/features/measurement/data/measurement_stream_wire_contract_test.dart \
  test/features/measurement/data/measurement_repository_impl_test.dart \
  test/features/measurement/data/measurement_repository_rest_test.dart \
  test/features/measurement/data/measurement_trace_sink_rest_test.dart \
  test/features/measurement/application/measurement_golden_path_orchestrator_test.dart \
  test/core/providers/measurement_history_provider_test.dart \
  test/features/measurement/domain/measurement_domain_test.dart \
  test/features/domain_models_test.dart
flutter analyze --no-pub \
  lib/generated/manpasik.pb.dart \
  lib/generated/manpasik.pbgrpc.dart \
  lib/features/measurement/data/measurement_repository_impl.dart \
  test/features/measurement/data/measurement_stream_wire_contract_test.dart \
  test/features/measurement/data/measurement_repository_impl_test.dart
```

결과: PASS

## 다음 단계

- in-process gRPC transport smoke는 `docs/audit/measure-grpc-transport-smoke-gate.md`에서 추가 완료했다.
- 실제 service process/container smoke 하네스는 후속 gate들에서 추가 완료했다.
- 전체 Dart generated proto compile gate는 `docs/audit/dart-proto-full-generated-compile-gate.md`에서 완료했다.
- Flutter 앱 의존성 lock 정렬과 checked-in official generated file 교체는 `docs/audit/dart-proto-official-generated-replacement.md`에서 완료했다.
- 다음 단계는 Auth login 사용자 식별자 보강 또는 Docker Compose container smoke runtime PASS 확보다.
