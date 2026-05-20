# Measure gRPC Service Binary Smoke 게이트

**작성일**: 2026-05-01  
**작성자**: Codex  
**범위**: 실제 `measurement-service` `cmd/main.go` 바이너리 cross-process 검증

## 상세 구현계획

1. `go build`로 `backend/services/measurement-service/cmd`의 production `main` 바이너리를 임시 디렉터리에 생성한다.
2. 테스트가 `127.0.0.1:0`으로 빈 gRPC/HTTP 포트를 예약한 뒤, 실제 바이너리를 별도 OS process로 기동한다.
3. runtime 환경은 `GRPC_PORT`, `HTTP_PORT`, `VERSION`, `SHUTDOWN_TIMEOUT_SECONDS`, `TENANCY_ENFORCED=false`만 명시 주입한다.
4. `DB_HOST`, `MILVUS_HOST`, `KAFKA_BROKERS`, `ELASTICSEARCH_URL`은 주입하지 않아 main wiring이 memory repository fallback으로 기동되는지 확인한다.
5. HTTP `/health`가 `status=healthy`, `service=measurement-service`, `version=smoke-test`를 반환할 때까지 poll 한다.
6. gRPC health service가 `measurement-service`를 `SERVING`으로 반환할 때까지 poll 한다.
7. generated Go `MeasurementServiceClient`로 `StartSession -> StreamMeasurement -> GetMeasurementHistory -> EndSession` lifecycle을 실행한다.
8. 측정 응답의 primary value, unit, confidence, fingerprint vector, `ProcessedAt`, history 저장 조회, session 종료 응답을 검증한다.
9. 테스트 종료 시 실제 service process에 `SIGTERM`을 보내고 main의 graceful shutdown path를 통해 종료한다.

## 구현 결과

- `cmd/main.go`를 직접 호출한 것은 아니고, `go build`로 생성한 실제 service binary를 별도 process로 실행했다.
- helper-process smoke에서 남아 있던 "test binary 안에서 handler를 조립한다"는 한계를 제거했다.
- main wiring의 다음 표면을 검증했다.
  - env 기반 `GRPC_PORT`, `HTTP_PORT`, `VERSION`, shutdown timeout
  - 외부 DB/Milvus/Kafka/Elasticsearch 미설정 시 memory fallback
  - HTTP health endpoint
  - gRPC health service
  - generated client/server transport
  - measurement session, stream processing, history read, session end lifecycle

## 변경 파일

- `backend/services/measurement-service/cmd/main_smoke_test.go`
  - `TestMeasurementServiceBinarySmoke` 추가.
  - 실제 `measurement-service` binary build/start/stop helper 추가.
  - HTTP health와 gRPC health readiness poll helper 추가.
  - `StartSession`, `StreamMeasurement`, `GetMeasurementHistory`, `EndSession` lifecycle 검증.
- `docs/audit/measure-grpc-service-binary-smoke-gate.md`
  - 상세 구현계획, 실행 범위, 품질 게이트 문서화.
- `docs/audit/measure-grpc-helper-process-smoke-gate.md`
  - 실제 service binary 후속 gate 완료 상태 반영.
- `docs/audit/mock-retirement-register.md`
  - MR-009 증거에 service binary smoke 완료 추가.

## 품질 게이트

```bash
cd /home/kangjh3kang/Manpasik
/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -v \
  ./backend/services/measurement-service/cmd \
  -run TestMeasurementServiceBinarySmoke
/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 \
  ./backend/services/measurement-service/cmd \
  ./backend/services/measurement-service/internal/handler \
  ./backend/services/measurement-service/internal/service

cd /home/kangjh3kang/Manpasik/frontend/flutter-app
bash scripts/check_proto_generation_preflight.sh
export PATH=/mnt/d/우리집/flutter_cache/flutter/bin:/usr/local/bin:/usr/bin:/bin
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

- containerized measurement-service smoke 하네스는 `docs/audit/measure-grpc-compose-container-smoke-gate.md`에서 추가 완료했다.
- 현재 WSL 환경에서는 Docker Compose 연결이 없어 실제 container up은 blocked 상태다.
- Dart proto toolchain 버전 정렬 후 전체 generated Dart compile gate를 실행한다.
- service binary smoke를 Gateway `/measurements/process` e2e와 연결해 Flutter native direct stream, Gateway bridge, measurement-service 저장 경로의 선택 기준을 더 명확히 한다.
