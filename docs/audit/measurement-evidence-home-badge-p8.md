# Measurement Evidence Home Badge P8

## 목적

P7에서 결과 화면에 표시한 evidence badge를 Home dashboard의 최근 측정 요약에도 확장했다. 사용자가 첫 화면에서 최근 측정을 볼 때 `research_only` 상태가 `연구용` 배지로 일관되게 표시되도록 하는 범위다.

## Task 1 완료 범위: Home Badge Screen Test

- `frontend/flutter-app/test/features/home/presentation/home_measurement_evidence_badge_test.dart`
  - `homeDashboardProvider`를 override해 `research_only` latest measurement를 주입했다.
  - HomeScreen에서 `연구용` 배지가 렌더링되는지 검증했다.
  - "정상", "위험", "진단", "확정" 표현이 표시되지 않는지 검증했다.

## Task 1 TDD 기록

- RED: HomeScreen이 latest measurement evidence badge를 렌더링하지 않아 `연구용` finder가 0건으로 FAIL.
- GREEN: Home hero card에 `MeasurementEvidenceBadge`를 연결해 screen widget test PASS.

## 완료 범위

- `frontend/flutter-app/lib/features/home/presentation/home_screen.dart`
  - `MeasurementEvidenceBadge`를 import했다.
  - `_HeroBentoCard`의 최근 측정 정보 아래에 latest measurement evidence badge를 표시한다.

## 자체 코드리뷰

- Home 화면도 P7 badge widget을 재사용하므로 문구 생성 경로가 P4 helper로 유지된다.
- `research_only`는 `연구용`으로만 표시하며 진단/정상/위험/확정 표현을 추가하지 않았다.
- DataHub와 다른 카드 확장은 별도 단계로 남겼다.

## 품질 게이트

- RED: `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/home/presentation/home_measurement_evidence_badge_test.dart`: FAIL, badge not rendered
- GREEN: `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/home/presentation/home_measurement_evidence_badge_test.dart`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/home/presentation/home_measurement_evidence_badge_test.dart test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart`: PASS
- `/mnt/d/우리집/flutter_cache/flutter/bin/flutter analyze --no-pub --no-fatal-infos` targeted P8 Flutter files: PASS

## 다음 단계 지침

1. P9에서는 DataHub trend/summary에도 evidence metadata를 확장할지 결정한다.
2. DataHub 확장은 `TrendDataPoint`에 evidence fields를 넣을지, summary card에만 표시할지 먼저 계획한다.
3. 실제 smoke에서는 REST history response에서 HomeScreen badge까지 이어지는 화면 경로를 확인한다.
