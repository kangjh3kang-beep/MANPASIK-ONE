# DataHub Total Count Compatibility Audit (P14)

## 목적

P9-P12에서 DataHub history item은 snake_case/camelCase를 모두 처리하도록 강화했지만, `getTotalMeasurementCount`는 `total_count`만 읽고 있었다. gateway 또는 legacy mock이 `totalCount`를 반환하면 총 측정 수가 0으로 떨어질 수 있어 count field도 동일한 호환성 정책으로 맞춘다.

## 변경 사항

- `frontend/flutter-app/test/features/data_hub/data/data_hub_repository_rest_test.dart`
  - history response가 `{"measurements": [], "totalCount": 7}`을 반환할 때 `getTotalMeasurementCount()`가 7을 반환하는지 검증했다.
- `frontend/flutter-app/lib/features/data_hub/data/data_hub_repository_rest.dart`
  - `_intField` helper를 추가했다.
  - `getTotalMeasurementCount`가 `total_count`와 `totalCount`를 모두 읽게 했다.

## TDD 기록

- RED: `flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart`
  - 실패 이유: `totalCount` response에서 count가 0으로 반환됨.
- GREEN: 같은 테스트 재실행 PASS.

## 자체 코드리뷰

- 기존 snake_case 경로는 유지하고 camelCase fallback만 추가했다.
- helper는 `num`을 `toInt()`로 처리해 JSON decoder가 int/double을 주더라도 안정적으로 동작한다.
- history item field fallback 정책과 count field fallback 정책이 일관된다.

## 품질 게이트

- `flutter test --no-pub test/features/data_hub/presentation/data_hub_layout_smoke_test.dart test/features/data_hub/presentation/data_hub_evidence_badge_test.dart test/features/data_hub/data/data_hub_repository_rest_test.dart test/features/data_hub/domain/data_hub_domain_test.dart test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart`: PASS
- `flutter analyze --no-pub --no-fatal-infos` targeted P14 Flutter files: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS

## 다음 단계 지침

- P15에서는 DataHub export path의 response field fallback(`file_path`/`fhir_json`)과 record count camelCase compatibility를 검토한다.
- WebAssembly dry-run compatibility warning은 별도 platform compatibility plan으로 분리한다.
