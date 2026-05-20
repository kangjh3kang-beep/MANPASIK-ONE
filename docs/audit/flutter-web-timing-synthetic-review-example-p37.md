# Flutter Web Timing Synthetic Review Example P37 Audit

## Scope

P37은 P36 release-review checklist를 비운영 synthetic fixture와 audit example로 재현 검증하는 단계다.
Synthetic fixture는 real timing review directory convention을 따르지만 production evidence로 사용할 수 없으며, checklist path와 wrapper command가 실제 파일 구조에서 동작하는지 검증한다.

## Changes

- `scripts/ci_web_timing_synthetic_review_example_test.sh`
  - Synthetic fixture directory와 audit file marker를 확인하고, fixture validator 및 threshold review wrapper를 실행한다.
- `docs/ci/fixtures/flutter-web-timing/2026-05-19-synthetic-review/`
  - `manifest.env`, five synthetic sample files, `aggregate.env`를 추가했다.
- `docs/audit/flutter-web-timing-synthetic-review-example.md`
  - Wrapper가 검증할 수 있는 synthetic audit example을 추가했다.
- `scripts/ci_web_release_review_checklist_policy_test.sh`
  - Checklist policy guard가 synthetic example script, fixture path, audit path를 요구하고 synthetic test를 실행하도록 확장했다.
- `docs/ci/flutter-web-release-review-checklist.md`
  - Synthetic dry-run command와 fixture/audit example path를 문서화했다.
- `docs/superpowers/plans/2026-05-19-flutter-web-timing-synthetic-review-example-p37.md`
  - P37 RED-GREEN 구현 순서와 검증 명령을 기록했다.

## TDD Record

- RED: `bash scripts/ci_web_timing_synthetic_review_example_test.sh`
  - 실패: `missing synthetic review dir`
- GREEN: `bash scripts/ci_web_timing_synthetic_review_example_test.sh`
  - 통과: `CI_WEB_TIMING_SYNTHETIC_REVIEW_EXAMPLE_PASS`
- RED: `bash scripts/ci_web_release_review_checklist_policy_test.sh`
  - 실패: `checklist missing marker: synthetic_review_example: docs/ci/fixtures/flutter-web-timing/2026-05-19-synthetic-review`
- GREEN: `bash scripts/ci_web_release_review_checklist_policy_test.sh`
  - 통과: `CI_WEB_RELEASE_REVIEW_CHECKLIST_POLICY_PASS`

## Verification

- `bash scripts/ci_web_timing_synthetic_review_example_test.sh`: PASS
- `bash scripts/ci_web_release_review_checklist_policy_test.sh`: PASS
- `bash scripts/ci_web_timing_review_fixture_policy_test.sh`: PASS
- `bash scripts/ci_gate_matrix_policy_test.sh`: PASS
- `bash -n scripts/ci_web_timing_synthetic_review_example_test.sh scripts/ci_web_release_review_checklist_policy_test.sh scripts/flutter_web_timing_review_fixture_validate.sh scripts/flutter_web_timing_threshold_change_review_check.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS

## Self Review

- Synthetic example is explicitly marked with `synthetic_example=true` and `synthetic_review: true`.
- The synthetic audit decision is `rejected`, so the example cannot be interpreted as approval for a real threshold change.
- Checklist policy guard now executes the synthetic example test, exercising validator and wrapper paths in CI.

## Residual Notes

- P38 can add a release-governance index page that links all Flutter web timing CI policies from one entrypoint.
- No production timing artifacts were added in this step.
