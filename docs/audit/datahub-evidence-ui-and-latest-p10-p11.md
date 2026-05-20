# DataHub Evidence UI And Latest Selection Audit (P10-P11)

## 목적

P9에서 DataHub domain/REST 계층이 measurement evidence metadata를 보존하도록 확장했지만, 화면에는 안전한 evidence 상태가 표시되지 않았고 summary latest evidence가 응답 순서에 의존했다. 이번 단계는 P4/P7 배지 표시 원칙을 재사용해 DataHub UI까지 연결하고, `measured_at` 기준 최신 evidence 선택을 테스트로 고정한다.

## 변경 사항

- `frontend/flutter-app/lib/features/data_hub/presentation/data_hub_screen.dart`
  - `_HeroChartCard`와 `_DetailPanel`에 `MeasurementEvidenceBadge`를 연결했다.
  - 새 의료 판정 문구를 만들지 않고 `MeasurementEvidencePresentation` 경로를 재사용한다.
- `frontend/flutter-app/lib/features/data_hub/data/data_hub_repository_rest.dart`
  - REST history item을 `_trendPointFromMeasurement`로 표준화한다.
  - `getTrendData` 결과를 timestamp 오름차순으로 정렬한다.
  - `getBiomarkerSummary`와 `getAllBiomarkerSummaries`가 timestamp 최신 point의 evidence metadata와 latest value를 사용한다.
- `frontend/flutter-app/test/features/data_hub/presentation/data_hub_evidence_badge_test.dart`
  - DataHub 화면이 `research_only` summary를 `연구용` 배지로 표시하고 금지 문구를 노출하지 않는지 검증한다.
- `frontend/flutter-app/test/features/data_hub/data/data_hub_repository_rest_test.dart`
  - out-of-order history response에서 timestamp 정렬과 최신 evidence 선택을 검증한다.

## TDD 기록

- RED: `flutter test --no-pub test/features/data_hub/presentation/data_hub_evidence_badge_test.dart`
  - 실패 이유: `DataHubScreen`에 `연구용` 배지가 렌더링되지 않음.
- GREEN: 같은 테스트 재실행 PASS.
- RED: `flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart`
  - 실패 이유: trend data가 응답 순서를 유지하고 summary latest evidence/value가 timestamp 최신값을 선택하지 않음.
- GREEN: 같은 테스트 재실행 PASS.

## 자체 코드리뷰

- 의료/규정 문구: DataHub 화면에서 "정상", "위험", "진단", "확정" 같은 판정 표현을 새로 만들지 않았다.
- 하위 호환성: evidence가 없는 legacy response는 기존 기본값 `unknown`, `false`, `[]` 경로를 유지한다.
- 결정성: REST history 응답 순서와 무관하게 parsed `measured_at` 기준으로 latest evidence를 선택한다.
- 범위 제한: DataHub UI와 REST summary mapping만 수정했고 차트 알고리즘/점수 산식은 별도 단계로 남겼다.

## 품질 게이트

- `flutter test --no-pub test/features/data_hub/presentation/data_hub_evidence_badge_test.dart test/features/data_hub/domain/data_hub_domain_test.dart test/features/data_hub/data/data_hub_repository_rest_test.dart test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart`: PASS
- `flutter analyze --no-pub --no-fatal-infos` targeted P10/P11 files: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted P10/P11 tracked files: PASS
- P10/P11 targeted trailing whitespace check: PASS

## 다음 단계 지침

- P12에서는 DataHub trend 산식이 value-sort 기반 `_computeTrend`에 의존하는 문제를 timestamp 기반 추세 계산으로 분리해 검증한다.
- 이후 Flutter 화면 smoke 또는 golden-like screenshot 검증으로 DataHub 배지가 모바일/데스크톱에서 overflow 없이 표시되는지 확인한다.
