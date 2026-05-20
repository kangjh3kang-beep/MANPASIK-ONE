# Measure gRPC Helper Process Smoke 게이트

**작성일**: 2026-05-01  
**작성자**: Codex  
**범위**: measurement-service `StreamMeasurement` helper-process cross-process 검증

## 상세 구현계획

1. 현재 Go test binary를 별도 child process로 재실행해 parent process와 gRPC server process를 분리한다.
2. child process는 memory repository 기반 `MeasurementService`와 `MeasurementHandler`를 구성하고 `127.0.0.1:0` TCP listener에 `grpc.Server`를 기동한다.
3. child process가 실제 bind 주소를 임시 address file에 기록하면 parent process가 generated Go `MeasurementServiceClient`로 접속한다.
4. parent process는 `StartSession -> StreamMeasurement -> GetMeasurementHistory -> EndSession` 순서로 실제 gRPC API lifecycle을 실행한다.
5. `StreamMeasurement` 응답의 primary value, unit, confidence, fingerprint vector, `ProcessedAt`을 검증한다.
6. `GetMeasurementHistory`가 같은 child process repository에 저장된 측정 결과를 반환하는지 확인한다.
7. 테스트 종료 시 parent process가 child process에 interrupt signal을 보내고 child process는 `GracefulStop()`으로 gRPC server를 종료한다.

## 구현 결과

- `grpc_stream_transport_smoke_test.go`에 helper-process cross-process smoke를 추가했다.
- 기존 TCP loopback smoke는 같은 process 안에서 OS network stack을 통과했지만, 이번 gate는 client와 server를 다른 process로 분리한다.
- 실제 production `cmd/main.go` binary를 기동한 것은 아니며, test binary helper process 안에서 동일 handler/service/repository 조립을 사용한다.
- 따라서 이번 gate는 process boundary, TCP transport, generated Go client/server, session lifecycle, 저장 조회 연동을 검증한다.
- 다음 gate에서는 실제 `measurement-service` binary 또는 dev compose service를 대상으로 환경 변수, 포트, repository readiness까지 포함해 검증한다.

## 변경 파일

- `backend/services/measurement-service/internal/handler/grpc_stream_transport_smoke_test.go`
  - `TestStreamMeasurementAgainstHelperProcess` 추가.
  - `TestMeasurementServiceProcessHelper` 추가.
  - `startMeasurementServiceHelperProcess` helper 추가.
- `docs/audit/measure-grpc-helper-process-smoke-gate.md`
  - 상세 구현계획, 실행 범위, 한계, 품질 게이트 문서화.
- `docs/audit/measure-grpc-tcp-loopback-smoke-gate.md`
  - helper-process 후속 gate 완료 상태 반영.
- `docs/audit/measure-grpc-transport-smoke-gate.md`
  - transport smoke 증거 사슬에 helper-process gate 완료 반영.
- `docs/audit/mock-retirement-register.md`
  - MR-009 증거에 helper-process cross-process smoke 완료 추가.

## 품질 게이트

```bash
cd /home/kangjh3kang/Manpasik
/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -v \
  ./backend/services/measurement-service/internal/handler \
  -run TestStreamMeasurementAgainstHelperProcess
/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 \
  ./backend/services/measurement-service/internal/handler \
  ./backend/services/measurement-service/internal/service

cd /home/kangjh3kang/Manpasik/frontend/flutter-app
bash scripts/check_proto_generation_preflight.sh
export PATH="/mnt/d/우리집/flutter_cache/flutter/bin:/usr/local/bin:/usr/bin:/bin"
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
```

결과: PASS

## 다음 단계

- 실제 `backend/services/measurement-service/cmd/main.go` binary를 별도 process로 띄우는 service-process smoke는 `docs/audit/measure-grpc-service-binary-smoke-gate.md`에서 추가 완료했다.
- dev compose 환경에서 measurement-service health/readiness와 gRPC port를 포함한 smoke gate를 추가한다.
- Dart proto toolchain 버전 정렬 후 전체 generated Dart compile gate를 실행한다.
