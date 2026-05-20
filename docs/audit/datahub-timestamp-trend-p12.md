# DataHub Timestamp Trend Audit (P12)

## 목적

P11에서 DataHub summary latest evidence는 `measured_at` 기준으로 안정화했지만, trend 계산은 값 정렬 리스트를 입력으로 받아 시간 순서의 상승/하강을 왜곡할 수 있었다. 이번 단계는 추세 판정을 timestamp 순서의 값 흐름으로 계산하게 고정한다.

## 변경 사항

- `frontend/flutter-app/test/features/data_hub/data/data_hub_repository_rest_test.dart`
  - older 130, middle 100, newer 70으로 감소하는 history를 out-of-order response로 주입한다.
  - `getBiomarkerSummary`와 `getAllBiomarkerSummaries`가 모두 `falling`을 반환해야 한다.
- `frontend/flutter-app/lib/features/data_hub/data/data_hub_repository_rest.dart`
  - summary trend에는 timestamp 정렬된 `orderedValues`를 사용한다.
  - min/max에는 별도 `sortedValues`를 사용해 통계 집계와 추세 판정을 분리한다.
  - `getAllBiomarkerSummaries`는 biomarker별 `TrendDataPoint` 목록을 보관하고 timestamp 정렬 후 latest/trend를 계산한다.

## TDD 기록

- RED: `flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart`
  - 실패 이유: 감소하는 timestamp sequence가 `rising`으로 계산됨.
- GREEN: 같은 테스트 재실행 PASS.

## 자체 코드리뷰

- 기존 P11의 timestamp 최신 evidence 계약과 같은 `TrendDataPoint` 정렬 경로를 사용한다.
- 평균/min/max는 기존 통계 의미를 유지하고, trend만 시간 순서 기반으로 분리했다.
- response order가 뒤섞여도 `getTrendData`, `getBiomarkerSummary`, `getAllBiomarkerSummaries`가 같은 시간 기준을 공유한다.

## 품질 게이트

- `flutter test --no-pub test/features/data_hub/presentation/data_hub_evidence_badge_test.dart test/features/data_hub/domain/data_hub_domain_test.dart test/features/data_hub/data/data_hub_repository_rest_test.dart test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart`: PASS
- `flutter analyze --no-pub --no-fatal-infos` targeted P10-P12 Flutter files: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `flutter build web --no-pub`: PASS, WebAssembly dry-run compatibility warnings only
- `git diff --check` targeted P10-P12 tracked files: PASS
- P10-P12 targeted trailing whitespace check: PASS

## 다음 단계 지침

- P13에서는 DataHub 배지와 trend 카드가 모바일/데스크톱 레이아웃에서 overflow 없이 보이는지 Flutter widget constraint test 또는 Playwright/browser smoke로 검증한다.
- 이후 DataHub trend 계산 임계값 5%가 임상/연구용 문맥에 적절한지 assay evidence 정책과 연결해 검토한다.
