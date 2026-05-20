# DataHub Export Compatibility Audit (P15)

## 목적

P14에서 DataHub total count의 snake/camel compatibility를 맞췄다. 이번 단계는 export response에서도 `filePath`, `fhirJson`, `recordCount` 같은 camelCase 응답을 잃지 않게 해 DataHub export 결과의 파일 경로와 record count를 보존한다.

## 변경 사항

- `frontend/flutter-app/test/features/data_hub/data/data_hub_repository_rest_test.dart`
  - `/api/v1/health-records/export/fhir` local server helper를 추가했다.
  - `{"filePath": "/tmp/export.fhir.json", "recordCount": 3}` response를 `ExportResult`로 보존하는지 검증했다.
- `frontend/flutter-app/lib/features/data_hub/data/data_hub_repository_rest.dart`
  - `exportData`가 `file_path`, `filePath`, `fhir_json`, `fhirJson` 순서로 파일 경로를 해석한다.
  - `record_count`, `recordCount`를 모두 record count로 처리한다.

## TDD 기록

- RED: `flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart`
  - 실패 이유: camelCase `filePath` response에서 `ExportResult.filePath`가 empty string으로 반환됨.
- GREEN: 같은 테스트 재실행 PASS.

## 자체 코드리뷰

- public `ExportResult` API는 변경하지 않았다.
- 기존 snake_case 응답과 `fhir_json` fallback은 유지했다.
- P14의 `_intField`를 재사용해 count fallback 정책을 통일했다.

## 품질 게이트

- `flutter test --no-pub test/features/data_hub/presentation/data_hub_layout_smoke_test.dart test/features/data_hub/presentation/data_hub_evidence_badge_test.dart test/features/data_hub/data/data_hub_repository_rest_test.dart test/features/data_hub/domain/data_hub_domain_test.dart test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart`: PASS
- `flutter analyze --no-pub --no-fatal-infos` targeted P15 Flutter files: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `flutter build web --no-pub`: PASS, WebAssembly dry-run compatibility warnings only
- `git diff --check` targeted P15 tracked files: PASS
- P15 targeted trailing whitespace check: PASS

## 다음 단계 지침

- P16에서는 Wasm dry-run compatibility warning을 별도 platform plan으로 분리해, 실제 지원 대상이 JS web인지 Wasm web인지 정책을 확정한다.
- DataHub REST compatibility는 현재 history item, total count, export result에 대해 snake/camel fallback을 갖췄다.
