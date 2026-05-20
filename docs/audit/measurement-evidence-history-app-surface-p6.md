# Measurement Evidence History App Surface P6

## 목적

P5에서 `MeasurementSummary` proto와 저장 계층까지 확장된 evidence fields를 gateway REST history와 Flutter history domain/native/REST repository까지 보존했다. 이번 범위는 history item 데이터 표면 확장이며, UI 상세 문구나 의료 판정 로직은 포함하지 않는다.

## Task 1 완료 범위: Gateway REST History Contract

- `backend/services/gateway/internal/handler/e2e_test.go`
  - `/api/v1/measurements/history` 응답이 `measurements`와 함께 `evidence_status`, `diagnostic_ready`, `evidence_gaps`, `research_only`를 포함하는지 검증했다.
  - mock `GetMeasurementHistory` response가 evidence fields를 가진 `MeasurementSummary` 1건을 반환하게 했다.

## Task 1 TDD 기록

- RED: gateway history response가 `{"measurements":[],"total_count":0}`만 반환해 `evidence_status` assertion 실패.
- GREEN: mock history response에 `research_only`, `DiagnosticReady=false`, `clinical_lock_required` gap을 추가해 route test PASS.

## Task 2 완료 범위: Flutter History Domain and Mappers

- `frontend/flutter-app/lib/features/measurement/domain/measurement_repository.dart`
  - `MeasurementHistoryItem`에 `evidenceStatus`, `diagnosticReady`, `evidenceGaps`를 추가했다.
  - legacy 기본값은 `unknown`, `false`, `[]`다.
- `frontend/flutter-app/lib/features/measurement/data/measurement_repository_impl.dart`
  - native gRPC `MeasurementSummary` evidence fields를 history item으로 매핑한다.
  - 테스트 가능하도록 optional `MeasurementHistoryCall` override를 추가했다.
- `frontend/flutter-app/lib/features/measurement/data/measurement_repository_rest.dart`
  - REST history item의 snake_case와 legacy camelCase evidence fields를 모두 decode한다.
- `frontend/flutter-app/test/features/measurement/data/measurement_repository_impl_test.dart`
  - native history evidence mapping contract를 추가했다.
- `frontend/flutter-app/test/features/measurement/data/measurement_repository_rest_test.dart`
  - REST history snake_case/camelCase evidence mapping contract를 추가했다.

## Task 2 TDD 기록

- RED: Flutter tests가 `MeasurementHistoryItem` evidence getters 부재와 `getMeasurementHistory` override 부재로 compile fail.
- GREEN: domain fields, native history override/mapping, REST decode helpers를 추가해 targeted Flutter tests PASS.

## 자체 코드리뷰

- `research_only`는 history item에서도 그대로 보존하며 `diagnosticReady=false`를 유지한다.
- REST mapper는 snake_case를 기본으로 읽고 camelCase도 허용해 gateway/legacy 응답을 함께 수용한다.
- 누락된 evidence status는 `unknown`으로 처리해 오래된 응답이 진단 가능 상태로 보이지 않는다.
- native repository의 optional override는 테스트 주입용이며 기존 생성자 호출은 깨지지 않는다.

## 품질 게이트

- RED: `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./services/gateway/internal/handler -run TestE2E_GetMeasurementHistory`: FAIL, evidence fields missing
- GREEN: `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./services/gateway/internal/handler -run TestE2E_GetMeasurementHistory`: PASS
- RED: `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/measurement/data/measurement_repository_rest_test.dart test/features/measurement/data/measurement_repository_impl_test.dart`: FAIL, history domain fields and native override missing
- GREEN: `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/measurement/data/measurement_repository_rest_test.dart test/features/measurement/data/measurement_repository_impl_test.dart`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `/mnt/d/우리집/flutter_cache/flutter/bin/flutter analyze --no-pub --no-fatal-infos` targeted P6 Flutter files: PASS
- `git diff --check` targeted P6 tracked files: PASS
- P6 targeted trailing whitespace check: PASS

## 다음 단계 지침

1. P7에서는 history 화면이나 data hub가 evidence label을 표시할 경우 P4의 `MeasurementEvidencePresentation`을 재사용한다.
2. 실제 gateway REST history smoke는 measurement-service history와 gateway route를 함께 띄워 gRPC→REST→Flutter 경로를 확인한다.
3. 운영 DB에는 init SQL 외에 기존 테이블용 ALTER migration을 별도 절차로 추가한다.
