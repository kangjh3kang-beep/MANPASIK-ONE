# Measurement Evidence Persistence P5

## 목적

P4까지 UI 표시용으로 전달되던 measurement evidence fields를 저장 스키마와 history response contract까지 확장했다. 이번 단계는 `research_only`, `diagnostic_ready=false`, `evidence_gaps`를 손실 없이 보존하는 것이며, 임상 판정 문구나 추가 UI 설명은 포함하지 않는다.

## Task 1 완료 범위: History Contract RED

- `backend/services/measurement-service/internal/handler/grpc_stream_test.go`
  - `GetMeasurementHistory` 응답의 `MeasurementSummary`가 evidence fields를 노출하는지 검증했다.
- `backend/services/measurement-service/internal/service/measurement_test.go`
  - mock repository history summary가 evidence fields를 보존하는지 검증했다.
- `backend/services/measurement-service/internal/repository/postgres/measurement_schema_test.go`
  - `04-measurement.sql`에 evidence columns와 summary view projection이 있는지 정적 guard를 추가했다.
- `frontend/flutter-app/test/generated/measurement_summary_evidence_contract_test.dart`
  - Dart generated `MeasurementSummary`가 evidence fields를 직렬화/역직렬화하는지 검증했다.

## Task 1 TDD 기록

- RED: Go handler/service tests가 `MeasurementSummary` evidence fields 부재로 build fail.
- RED: repository schema guard가 `evidence_status` 컬럼 부재로 fail.
- RED: Dart generated contract test가 `evidenceStatus`, `diagnosticReady`, `evidenceGaps` 부재로 compile fail.

## Task 2 완료 범위: Proto and Service Mapping

- `backend/shared/proto/manpasik.proto`
  - `MeasurementSummary`에 `evidence_status = 6`, `diagnostic_ready = 7`, `evidence_gaps = 8`을 append-only로 추가했다.
- Go/Dart generated output을 재생성했다.
- `backend/services/measurement-service/internal/service/measurement.go`
  - `MeasurementSummary`에 evidence fields를 추가했다.
- `backend/services/measurement-service/internal/handler/grpc.go`
  - history response summary에 evidence fields를 매핑했다.

## Task 3 완료 범위: Persistence Columns and Repository Mapping

- `infrastructure/database/init/04-measurement.sql`
  - `measurement_data`에 `evidence_status`, `diagnostic_ready`, `evidence_gaps`를 추가했다.
  - `measurement_summary` view가 새 evidence fields를 노출하게 했다.
- `backend/services/measurement-service/internal/repository/postgres/measurement.go`
  - Store INSERT에 evidence fields를 추가했다.
  - GetHistory SELECT/SCAN이 저장된 evidence fields를 `MeasurementSummary`로 복원한다.
  - 누락된 status는 `unknown`으로 보수적으로 처리한다.

## 자체 코드리뷰

- proto는 기존 `MeasurementSummary` field 1~5를 유지하고 6~8만 추가했다.
- 저장 기본값은 `unknown`, `false`, empty array라 legacy 데이터가 진단 가능 상태로 오인되지 않는다.
- repository는 empty status를 `unknown`으로 복원해 하위 호환 경로도 안전하게 실패한다.
- history response는 stream response와 같은 evidence semantics를 노출한다.

## 품질 게이트

- RED: `go test -count=1 ./services/measurement-service/internal/handler ./services/measurement-service/internal/service ./services/measurement-service/internal/repository/postgres`: FAIL, summary fields/schema missing
- RED: `flutter test --no-pub test/generated/measurement_summary_evidence_contract_test.dart`: FAIL, generated Dart summary evidence fields missing
- GREEN: `go test -count=1 ./services/measurement-service/internal/handler ./services/measurement-service/internal/service ./services/measurement-service/internal/repository/postgres`: PASS
- GREEN: `flutter test --no-pub test/generated/measurement_summary_evidence_contract_test.dart`: PASS
- `flutter analyze --no-pub --no-fatal-infos` targeted generated/test files: PASS
- `bash scripts/check_proto_generation_compile_gate.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted P5 tracked files: PASS
- P5 targeted trailing whitespace check: PASS

## 다음 단계 지침

1. P6에서는 Flutter history model/REST mapper가 `MeasurementSummary` evidence fields를 읽도록 별도 TDD로 확장한다.
2. 실제 migration 환경에서는 기존 measurement_data 테이블에 대한 ALTER migration 파일을 별도 운영 절차로 추가해야 한다.
3. 전체 smoke는 measurement-service gRPC history와 gateway REST history를 동시에 검증하는 경로로 확장한다.
