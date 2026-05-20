# Measure gRPC TCP Loopback Smoke 게이트

**작성일**: 2026-05-01  
**작성자**: Codex  
**범위**: measurement-service `StreamMeasurement` TCP loopback multi-frame 검증

## 상세 구현계획

1. 외부 Docker/Compose 의존 없이 `127.0.0.1:0` TCP listener를 열어 실제 OS network stack을 사용한다.
2. 기존 `MeasurementHandler`와 memory test repository를 `grpc.Server`에 등록한다.
3. generated Go `MeasurementServiceClient`로 TCP 주소에 접속해 `StreamMeasurement` bidi stream을 연다.
4. 한 stream 안에서 measurement frame 2개를 연속 전송하고 `CloseSend` 후 응답 2개와 EOF를 검증한다.
5. 각 응답의 primary value, unit, confidence, fingerprint vector, `ProcessedAt`을 확인한다.
6. measurement 저장소에 2건이 저장됐는지와 마지막 fingerprint vector가 vector 저장소에 반영됐는지 확인한다.

## 구현 결과

- `grpc_stream_transport_smoke_test.go`에 TCP loopback multi-frame smoke를 추가했다.
- 이전 `bufconn` smoke보다 한 단계 실제 네트워크 계층에 가까운 transport를 검증한다.
- multi-frame stream을 검증해 native direct stream이 향후 연속 측정 frame에도 대응할 수 있는 서버 계약을 고정했다.

## 변경 파일

- `backend/services/measurement-service/internal/handler/grpc_stream_transport_smoke_test.go`
  - `TestStreamMeasurementOverTCPLoopbackHandlesMultipleFrames` 추가.
- `docs/audit/measure-grpc-transport-smoke-gate.md`
  - TCP loopback 후속 gate 완료 상태 반영.
- `docs/audit/mock-retirement-register.md`
  - MR-009 증거에 TCP loopback multi-frame smoke 완료 추가.

## 품질 게이트

```bash
cd /home/kangjh3kang/Manpasik
/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -v \
  ./backend/services/measurement-service/internal/handler \
  -run TestStreamMeasurementOverTCP
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

- helper-process cross-process smoke는 `docs/audit/measure-grpc-helper-process-smoke-gate.md`에서 추가 완료했다.
- 실제 `measurement-service` binary service-process smoke는 `docs/audit/measure-grpc-service-binary-smoke-gate.md`에서 추가 완료했다.
- 다음 단계는 dev compose service를 별도 process/container로 띄우고 generated client가 실제 포트에 접속하는 container smoke다.
- Dart proto toolchain 버전 정렬 후 전체 generated Dart compile gate를 실행한다.
