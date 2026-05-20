# DataHub Evidence Metadata P9

## 목적

P8까지 Home/Result 화면에 표시되던 evidence status를 DataHub trend/summary 도메인 모델에서도 보존하도록 확장했다. 이번 단계는 UI 표시가 아니라 DataHub 분석 모델이 `research_only`, `diagnosticReady=false`, `evidenceGaps`를 잃지 않게 만드는 데이터 표면 보강이다.

## Task 1 완료 범위: Domain Evidence Metadata

- `frontend/flutter-app/lib/features/data_hub/domain/data_hub_repository.dart`
  - `TrendDataPoint`에 `evidenceStatus`, `diagnosticReady`, `evidenceGaps`를 추가했다.
  - `BiomarkerSummary`에 `latestEvidenceStatus`, `latestDiagnosticReady`, `latestEvidenceGaps`를 추가했다.
  - 기본값은 `unknown`, `false`, `[]`로 보수적으로 설정했다.
- `frontend/flutter-app/test/features/data_hub/domain/data_hub_domain_test.dart`
  - 기본값과 explicit evidence metadata 생성 계약을 검증했다.

## Task 1 TDD 기록

- RED: DataHub domain tests가 evidence field/getter 부재로 compile fail.
- GREEN: domain fields와 보수적 기본값을 추가해 domain test PASS.

## Task 2 완료 범위: REST Repository Evidence Mapping

- `frontend/flutter-app/lib/features/data_hub/data/data_hub_repository_rest.dart`
  - REST history item의 snake_case/camelCase fields를 읽는 helper를 추가했다.
  - `getTrendData()`가 evidence metadata를 `TrendDataPoint`로 전달한다.
  - `getBiomarkerSummary()`가 최신 point evidence metadata를 summary로 전달한다.
  - `getAllBiomarkerSummaries()`도 type별 첫 history item의 evidence metadata를 summary에 보존한다.
- `frontend/flutter-app/test/features/data_hub/data/data_hub_repository_rest_test.dart`
  - snake_case history evidence mapping을 검증했다.
  - legacy camelCase history evidence mapping을 검증했다.
  - `getBiomarkerSummary()` 최신 evidence metadata 보존을 검증했다.

## Task 2 TDD 기록

- RED: snake_case evidence status가 `unknown`으로 떨어지고, camelCase response는 측정 시각/필드명을 못 읽어 trend point가 0건으로 FAIL.
- GREEN: field helper와 evidence mapping을 추가해 DataHub REST tests PASS.

## 자체 코드리뷰

- UI 표시를 추가하지 않고 DataHub 모델 표면만 확장했다.
- legacy 응답의 기본값은 `unknown`, `false`, `[]`라 진단 가능 상태로 오인되지 않는다.
- `research_only`는 `diagnosticReady=false`와 함께 그대로 보존한다.
- summary latest evidence는 현재 history 응답 순서를 기준으로 첫 항목을 사용한다. 정렬 정책을 바꾸려면 별도 단계에서 measured_at 기반 정렬을 테스트로 고정한다.

## 품질 게이트

- RED: `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/data_hub/domain/data_hub_domain_test.dart`: FAIL, domain evidence fields missing
- GREEN: `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/data_hub/domain/data_hub_domain_test.dart`: PASS
- RED: `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart`: FAIL, REST evidence mapping missing
- GREEN: `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/data_hub/domain/data_hub_domain_test.dart test/features/data_hub/data/data_hub_repository_rest_test.dart`: PASS
- `/mnt/d/우리집/flutter_cache/flutter/bin/flutter analyze --no-pub --no-fatal-infos` targeted P9 Flutter files: PASS
- `git diff --check` targeted P9 tracked files: PASS
- P9 targeted trailing whitespace check: PASS

## 다음 단계 지침

1. P10에서 DataHub UI에 evidence badge를 표시할 경우 `MeasurementEvidenceBadge`를 재사용한다.
2. Summary latest evidence의 시간 기준 정렬이 필요하면 measured_at 기반 정렬 테스트를 먼저 작성한다.
3. REST history 실제 smoke는 gateway response에서 DataHub trend model까지 이어지는 경로를 검증한다.
