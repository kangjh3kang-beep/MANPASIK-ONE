# Measurement Evidence App Surface P3

## 목적

P2에서 gRPC/Dart proto 계약에 추가한 measurement evidence fields를 REST gateway와 Flutter repository result까지 노출했다. 이번 단계는 UI 표시 문구나 저장소 영속화가 아니라, 화면/저장 계층이 사용할 수 있는 안전한 데이터 표면을 만드는 범위다.

## Task 1 완료 범위: Gateway REST Contract

- `backend/services/gateway/internal/handler/e2e_test.go`
  - `/api/v1/measurements/process` REST 응답에 `evidence_status`, `diagnostic_ready`, `evidence_gaps`가 포함되는지 검증했다.
  - mock `StreamMeasurement` 응답을 `research_only`, `DiagnosticReady=false`, `clinical_lock_required` gap으로 맞췄다.

## Task 1 TDD 기록

- RED: route test가 `research_only` evidence status를 찾지 못해 FAIL. 응답에는 빈 `evidence_status`와 빈 gaps만 있었다.
- GREEN: mock stream response에 P2 evidence contract를 추가해 gateway route test PASS.

## Task 2 완료 범위: Flutter Native Repository

- `frontend/flutter-app/lib/features/measurement/domain/measurement_repository.dart`
  - `ProcessMeasurementResult`에 `evidenceStatus`, `diagnosticReady`, `evidenceGaps`를 추가했다.
  - legacy server 호환 기본값은 `unknown`, `false`, `[]`로 설정했다.
- `frontend/flutter-app/lib/features/measurement/data/measurement_repository_impl.dart`
  - native gRPC `MeasurementResult`의 evidence fields를 도메인 결과로 복사한다.
- `frontend/flutter-app/test/features/measurement/data/measurement_repository_impl_test.dart`
  - native process result가 evidence fields를 보존하는지 검증한다.

## Task 2 TDD 기록

- RED: `ProcessMeasurementResult`에 evidence getters가 없어 native repository test가 compile fail.
- GREEN: 도메인 필드와 native mapping을 추가해 target test PASS.

## Task 3 완료 범위: Flutter REST Mapper

- `frontend/flutter-app/lib/features/measurement/data/measurement_process_gateway_mapper.dart`
  - gateway proto JSON의 snake_case evidence keys를 읽는다.
  - 하위 호환을 위해 camelCase keys도 함께 읽는다.
- `frontend/flutter-app/test/features/measurement/data/measurement_process_gateway_mapper_test.dart`
  - snake_case와 camelCase response 모두 evidence fields를 도메인 결과로 매핑하는지 검증한다.

## Task 3 TDD 기록

- RED: REST mapper가 evidence fields를 읽지 못해 `evidenceStatus='unknown'`으로 떨어지고, camelCase 기본 필드는 fallback으로 빠졌다.
- GREEN: `_field()`와 `_numberField()` helper로 snake_case/camelCase를 모두 읽게 해 mapper tests PASS.

## 자체 코드리뷰

- `research_only`는 그대로 전달하며 진단 가능/확정 표현으로 바꾸지 않았다.
- `diagnosticReady` 기본값은 false라 legacy backend나 누락 응답이 안전한 쪽으로 실패한다.
- `evidenceGaps`는 새 리스트로 변환해 외부 응답 list와 도메인 객체의 결합을 줄였다.
- Gateway는 `writeProtoJSON`의 `UseProtoNames: true` 정책을 유지하므로 REST 표준 응답은 snake_case다.
- UI 표시와 저장소 영속화는 아직 하지 않았다. 문구/규정 리스크와 DB migration 리스크가 있어 다음 단계에서 별도 TDD로 진행한다.

## 품질 게이트

- RED: `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./services/gateway/internal/handler -run TestE2E_ProcessMeasurementGoldenPath`: FAIL, `research_only` missing
- GREEN: `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./services/gateway/internal/handler -run TestE2E_ProcessMeasurementGoldenPath`: PASS
- RED: `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/measurement/data/measurement_repository_impl_test.dart`: FAIL, `ProcessMeasurementResult` evidence fields missing
- GREEN: `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/measurement/data/measurement_repository_impl_test.dart`: PASS
- RED: `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/measurement/data/measurement_process_gateway_mapper_test.dart`: FAIL, evidence mapper missing
- GREEN: `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/measurement/data/measurement_process_gateway_mapper_test.dart`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./shared/assay ./services/measurement-service/internal/handler ./services/measurement-service/internal/service`: PASS
- `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_process_gateway_mapper_test.dart`: PASS
- `/mnt/d/우리집/flutter_cache/flutter/bin/flutter analyze --no-pub --no-fatal-infos lib/features/measurement/domain/measurement_repository.dart lib/features/measurement/data/measurement_repository_impl.dart lib/features/measurement/data/measurement_process_gateway_mapper.dart test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_process_gateway_mapper_test.dart`: PASS
- `git diff --check` targeted tracked files: PASS
- P3 untracked file trailing whitespace check: PASS

## 다음 단계 지침

1. P4에서 UI 표면을 추가할 때는 `research_only`를 "연구용/검증 준비 중"처럼 보수적으로 표시하고, 진단/정상/위험 확정 문구와 섞지 않는다.
2. 저장 계층 영속화는 Timescale/Postgres schema migration, repository test, history response contract를 먼저 작성한다.
3. Gateway REST smoke를 실제 measurement-service binary/compose smoke에 연결하면 gRPC→REST→Flutter 전체 경로 검증이 가능하다.
