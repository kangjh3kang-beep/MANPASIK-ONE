# Measure gRPC Transport Smoke 게이트

**작성일**: 2026-05-01  
**작성자**: Codex  
**범위**: measurement-service `StreamMeasurement` in-process gRPC transport 검증

## 상세 구현계획

1. `bufconn`으로 외부 포트 없이 실제 `grpc.Server`를 메모리 안에서 기동한다.
2. 기존 `MeasurementHandler`와 memory test repository를 서버에 등록한다.
3. generated Go `MeasurementServiceClient`로 `StreamMeasurement` bidi stream을 연다.
4. client `Send`, `CloseSend`, `Recv`, EOF 확인 순서로 실제 transport stream lifecycle을 검증한다.
5. 응답값뿐 아니라 measurement 저장소, fingerprint vector 저장소 side effect까지 확인한다.

## 구현 결과

- `grpc_stream_transport_smoke_test.go`를 추가했다.
- fake stream이 아닌 generated client/server transport를 통과하는 `StreamMeasurement` smoke를 확보했다.
- 측정 frame 전송 후 `MeasurementResult` 응답, `ProcessedAt`, measurement 저장, fingerprint vector 저장을 함께 검증한다.
- 이전 wire golden smoke와 조합되어 다음 범위를 커버한다.
  - Dart 수동 protobuf binary contract
  - Go generated proto unmarshal
  - measurement-service handler/service/repository path
  - in-process gRPC transport lifecycle

## 변경 파일

- `backend/services/measurement-service/internal/handler/grpc_stream_transport_smoke_test.go`
  - `bufconn` 기반 gRPC server/client smoke test 추가.
- `docs/audit/measure-stream-wire-smoke-gate.md`
  - network-level 후속 gate 완료 상태 반영.
- `docs/audit/mock-retirement-register.md`
  - MR-009 증거에 transport smoke 완료를 추가.

## 품질 게이트

```bash
cd /home/kangjh3kang/Manpasik
export PATH="/home/kangjh3kang/sdk/go-go1.26.2/bin:$PATH"
go test -count=1 ./backend/services/measurement-service/internal/handler ./backend/services/measurement-service/internal/service

cd /home/kangjh3kang/Manpasik/frontend/flutter-app
bash scripts/check_proto_generation_preflight.sh
export PATH="/mnt/d/우리집/flutter_cache/flutter/bin:/usr/local/bin:/usr/bin:/bin"
flutter test --no-pub \
  test/features/measurement/data/measurement_stream_wire_contract_test.dart \
  test/features/measurement/data/measurement_repository_impl_test.dart \
  test/features/measurement/data/measurement_repository_rest_test.dart \
  test/features/measurement/application/measurement_golden_path_orchestrator_test.dart
```

결과: PASS

## 다음 단계

- TCP loopback multi-frame smoke는 `docs/audit/measure-grpc-tcp-loopback-smoke-gate.md`에서 추가 완료했다.
- helper-process cross-process smoke는 `docs/audit/measure-grpc-helper-process-smoke-gate.md`에서 추가 완료했다.
- 실제 `measurement-service` binary service-process smoke는 `docs/audit/measure-grpc-service-binary-smoke-gate.md`에서 추가 완료했다.
- 다음 단계는 dev compose container smoke와 Dart generated proto compile gate다.
