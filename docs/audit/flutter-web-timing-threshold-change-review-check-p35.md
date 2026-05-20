# Flutter Web Timing Threshold Change Review Check P35 Audit

## Scope

P35는 threshold 변경 PR에서 reviewer가 concrete fixture validator와 audit template 값을 빠짐없이 실행/작성했는지 검증하는 wrapper command를 추가한 단계다.
Wrapper는 P34 validator를 실행한 뒤 `manifest.env`, `aggregate.env`, audit Markdown의 machine-check fields가 서로 일치하는지 확인한다.

## Changes

- `scripts/flutter_web_timing_threshold_change_review_check_test.sh`
  - 임시 review fixture와 audit file을 만들고 wrapper 성공 경로와 `p95_seconds` mismatch 실패 경로를 검증한다.
- `scripts/flutter_web_timing_threshold_change_review_check.sh`
  - `--review-dir`와 `--audit-file`을 받아 fixture validator 실행, aggregate field read, audit field consistency check를 수행한다.
- `scripts/ci_web_timing_review_fixture_policy_test.sh`
  - fixture policy가 wrapper script/test, required audit fields, template `validator_command` field를 문서화했는지 검증한다.
- `docs/ci/flutter-web-timing-review-fixtures.md`
  - threshold change review check command와 required audit fields를 문서화했다.
- `docs/audit/templates/flutter-web-timing-review-template.md`
  - wrapper가 읽을 수 있는 machine-check field block을 추가했다.
- `docs/superpowers/plans/2026-05-14-flutter-web-timing-threshold-change-review-check-p35.md`
  - P35 RED-GREEN 구현 순서와 검증 명령을 기록했다.

## TDD Record

- RED: `bash scripts/flutter_web_timing_threshold_change_review_check_test.sh`
  - 실패: `scripts/flutter_web_timing_threshold_change_review_check.sh` missing
- GREEN: `bash scripts/flutter_web_timing_threshold_change_review_check_test.sh`
  - 통과: `FLUTTER_WEB_TIMING_THRESHOLD_CHANGE_REVIEW_CHECK_TEST_PASS`
- RED: `bash scripts/ci_web_timing_review_fixture_policy_test.sh`
  - 실패: `policy missing marker: threshold_change_review_check_script: scripts/flutter_web_timing_threshold_change_review_check.sh`
- GREEN: `bash scripts/ci_web_timing_review_fixture_policy_test.sh`
  - 통과: `CI_WEB_TIMING_REVIEW_FIXTURE_POLICY_PASS`

## Verification

- `bash scripts/flutter_web_timing_threshold_change_review_check_test.sh`: PASS
- `bash scripts/flutter_web_timing_review_fixture_validate_test.sh`: PASS
- `bash scripts/ci_web_timing_review_fixture_policy_test.sh`: PASS
- `bash scripts/ci_web_timing_threshold_promotion_policy_test.sh`: PASS
- `bash scripts/ci_gate_matrix_policy_test.sh`: PASS
- `bash -n scripts/flutter_web_timing_threshold_change_review_check.sh scripts/flutter_web_timing_threshold_change_review_check_test.sh scripts/flutter_web_timing_review_fixture_validate.sh scripts/flutter_web_timing_review_fixture_validate_test.sh scripts/ci_web_timing_review_fixture_policy_test.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS

## Self Review

- Wrapper keeps threshold-change review local and deterministic.
- Audit file must record exact aggregate values, so a stale or partially copied review is rejected.
- Decision remains constrained to `proposed`, `approved`, or `rejected`.

## Residual Notes

- P36 can add a small `docs/ci` reviewer checklist page that links all timing gate commands in release-review order.
- No real timing review fixture was added in this step.
