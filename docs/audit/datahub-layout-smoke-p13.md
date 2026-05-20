# DataHub Layout Smoke Audit (P13)

## 목적

P10-P12에서 DataHub 화면에 evidence badge와 timestamp 기반 summary/trend를 연결했다. 이번 단계는 좁은 모바일 화면에서 긴 metric label, evidence badge, period chip이 함께 표시될 때 레이아웃 예외가 발생하지 않는지 회귀 smoke로 고정한다.

## 변경 사항

- `docs/superpowers/plans/2026-05-13-datahub-layout-smoke-p13.md`
  - P13 layout smoke 상세 계획과 검증 명령을 기록했다.
- `frontend/flutter-app/test/features/data_hub/presentation/data_hub_layout_smoke_test.dart`
  - 320px 모바일 폭, 긴 metric label, `research_only` evidence summary를 주입한다.
  - `연구용` 배지 표시와 `tester.takeException() == null`을 검증한다.

## TDD/검증 기록

- 초기 smoke: `flutter test --no-pub test/features/data_hub/presentation/data_hub_layout_smoke_test.dart`: PASS
- 계획상 overflow RED를 예상했지만 현 P10-P12 레이아웃은 좁은 폭에서도 예외 없이 동작했다.
- 생산 코드 수정은 하지 않고, 회귀 방지 테스트만 추가했다.

## 자체 코드리뷰

- 테스트는 실제 `DataHubScreen`과 Riverpod provider override를 사용한다.
- 새 의료 판정 문구를 만들지 않고 기존 `연구용` 배지 표면만 확인한다.
- 좁은 모바일 폭과 긴 metric label을 함께 사용해 layout smoke의 압력을 높였다.

## 품질 게이트

- `flutter test --no-pub test/features/data_hub/presentation/data_hub_layout_smoke_test.dart test/features/data_hub/presentation/data_hub_evidence_badge_test.dart test/features/data_hub/data/data_hub_repository_rest_test.dart test/features/data_hub/domain/data_hub_domain_test.dart test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart`: PASS
- `flutter analyze --no-pub --no-fatal-infos` targeted P13 Flutter files: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `flutter build web --no-pub`: PASS, WebAssembly dry-run compatibility warnings only

## 다음 단계 지침

- P14에서는 DataHub REST compatibility의 남은 작은 구멍인 `getTotalMeasurementCount`의 `totalCount` camelCase fallback을 TDD로 보강한다.
- 이후 Wasm dry-run 경고는 별도 web platform compatibility 계획으로 분리한다.
