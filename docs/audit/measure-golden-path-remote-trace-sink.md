# Measure Golden Path Remote Trace Sink Audit

**Date**: 2026-05-01  
**Owner**: Codex  
**Scope**: H1 Measure golden path observability remote intake

## 목적

Measure golden path는 로컬 `AppLogger` phase trace까지 갖췄지만, 운영 장애 분석에는 클라이언트에서 Gateway까지 도달하는 원격 trace intake가 필요하다. 이번 단계는 Flutter trace event를 Gateway로 전송하고, Gateway가 PHI 최소화 계약을 강제하는 경로를 추가했다.

## 구현 계약

### Flutter

- `MeasurementGoldenPathTraceEvent.toRemoteObservabilityJson()`
  - `schema_version=measure_trace.v1`
  - `source`, `route`, `phase`, `elapsed_ms`, `occurred_at`, `session_id`, `cartridge_id`, `engine_mode`, `unit`, `confidence`, `failure_reason`, `diagnostic_message`
  - `primary_value`는 전송하지 않고 `has_primary_value` boolean만 남긴다.
- `MeasurementGoldenPathRestTraceSink`
  - `ManPaSikRestClient.recordMeasurementTraceEvent()`로 `POST /measurements/trace-events`를 호출한다.
  - 전송 실패는 측정 플로우를 막지 않고 `AppLogger.warning`으로만 남긴다.
- `measurementGoldenPathTraceSinkProvider`
  - 로컬 로그 sink와 REST remote sink를 합성한다.
  - `/measure` 화면의 `MeasurementGoldenPathOrchestrator`가 이 provider를 사용한다.

### Gateway

- `POST /api/v1/measurements/trace-events`
  - `phase` 필수
  - `elapsed_ms >= 0`
  - `primary_value`가 포함되면 `400 Bad Request`
  - 정상 수신 시 `202 Accepted`
  - 현재 audit-service에는 write RPC가 없어 `audit_status`로 후속 연결 상태를 명시한다.

## 추가 파일

- `frontend/flutter-app/lib/features/measurement/data/measurement_trace_sink_rest.dart`
- `frontend/flutter-app/test/features/measurement/data/measurement_trace_sink_rest_test.dart`

## 변경 파일

- `frontend/flutter-app/lib/features/measurement/application/measurement_golden_path_orchestrator.dart`
- `frontend/flutter-app/lib/core/services/rest_client.dart`
- `frontend/flutter-app/lib/core/providers/grpc_provider.dart`
- `frontend/flutter-app/lib/features/measurement/presentation/measure_screen.dart`
- `frontend/flutter-app/test/features/measurement/application/measurement_golden_path_orchestrator_test.dart`
- `backend/services/gateway/internal/handler/measurement_routes.go`
- `backend/services/gateway/internal/handler/e2e_test.go`

## 품질 게이트

```bash
cd /home/kangjh3kang/Manpasik
export PATH=/home/kangjh3kang/sdk/go-go1.26.2/bin:/usr/local/bin:/usr/bin:/bin
gofmt -w backend/services/gateway/internal/handler/measurement_routes.go backend/services/gateway/internal/handler/e2e_test.go
go test ./backend/services/measurement-service/... ./backend/services/gateway/...
```

결과: PASS

```bash
cd /home/kangjh3kang/Manpasik/frontend/flutter-app
export PATH="/mnt/d/우리집/flutter_cache/flutter/bin:/usr/local/bin:/usr/bin:/bin"
dart format <changed Dart files>
flutter test --no-pub test/features/measurement/application/measurement_golden_path_orchestrator_test.dart test/features/measurement/data/measurement_trace_sink_rest_test.dart test/features/measurement/data/measurement_repository_rest_test.dart
flutter analyze --no-pub <changed Dart files>
```

결과: PASS

## 환경 메모

`flutter pub get`은 현재 WSL Flutter SDK cache의 `bin/cache/flutter.version.json` 누락 때문에 실패한다. 기존 lock/package config 기반의 `--no-pub` 테스트와 분석은 통과했다.

## 다음 단계

- audit-service에 write RPC 또는 dedicated intake를 추가하고 `audit_status`를 실제 persisted 상태로 승격한다.
- `docs/audit/mock-retirement-register.md`의 Measure 관련 mock/stub 항목을 native Rust/DB 연결 상태와 대조한다.
