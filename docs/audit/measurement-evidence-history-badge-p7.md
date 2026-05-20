# Measurement Evidence History Badge P7

## 목적

P6에서 Flutter history item까지 전달된 evidence fields를 결과 화면 최신 측정 카드에 compact badge로 표시했다. 이번 범위는 `MeasurementEvidencePresentation`을 재사용한 안전한 배지 표시이며, 의료 판정 문구나 data hub 전체 표시 확장은 포함하지 않는다.

## Task 1 완료 범위: Evidence Badge Widget

- `frontend/flutter-app/lib/features/measurement/presentation/widgets/measurement_evidence_badge.dart`
  - `MeasurementEvidencePresentation.from(...)`을 사용해 evidence status를 짧은 badge label로 렌더링한다.
- `frontend/flutter-app/test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart`
  - `research_only`가 `연구용`으로 표시되는지 검증했다.
  - "정상", "위험", "진단", "확정" 표현이 표시되지 않는지 검증했다.

## Task 1 TDD 기록

- RED: badge widget 파일/클래스가 없어 widget test compile fail.
- GREEN: badge widget을 추가해 widget test PASS.

## Task 2 완료 범위: Result Screen Integration

- `frontend/flutter-app/lib/features/measurement/presentation/measurement_result_screen.dart`
  - 최신 측정 카드에서 `MeasurementHistoryItem`의 `evidenceStatus`, `diagnosticReady`, `evidenceGaps`를 `MeasurementEvidenceBadge`로 표시한다.

## 자체 코드리뷰

- 배지 문구는 P4 helper를 통해서만 생성되며 화면에서 별도 의료 문구를 만들지 않는다.
- `research_only`는 `연구용`으로 표시되고 진단/정상/위험/확정 표현을 사용하지 않는다.
- 결과 화면의 기존 값, 카트리지 타입, AI 분석 카드 흐름은 유지했다.
- DataHub 확장은 별도 단계로 남겼다.

## 품질 게이트

- RED: `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart`: FAIL, widget missing
- GREEN: `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `/mnt/d/우리집/flutter_cache/flutter/bin/flutter analyze --no-pub --no-fatal-infos` targeted P7 Flutter files: PASS
- `git diff --check` targeted P7 tracked files: PASS
- P7 targeted trailing whitespace check: PASS

## 다음 단계 지침

1. P8에서 DataHub나 home dashboard에 evidence badge를 추가할 때도 `MeasurementEvidenceBadge`를 재사용한다.
2. 화면 수준 widget test가 필요하면 `measurementHistoryProvider` override 기반으로 `MeasurementResultScreen` 전체 렌더링을 검증한다.
3. 실제 서비스 smoke는 history API 응답이 화면 배지로 표시되는 경로까지 확인한다.
