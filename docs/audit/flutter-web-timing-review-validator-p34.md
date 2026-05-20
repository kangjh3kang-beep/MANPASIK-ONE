# Flutter Web Timing Review Validator P34 Audit

## Scope

P34는 concrete `docs/ci/fixtures/flutter-web-timing/<review-id>` directory가 threshold change PR 전에 검증 가능하도록 validator를 추가한 단계다.
Validator는 manifest, sample files, aggregate output, no-PHI attestation, forbidden sensitive fields, aggregate reproducibility를 검사한다.

## Changes

- `scripts/flutter_web_timing_review_fixture_validate_test.sh`
  - 임시 review fixture를 만들고 validator 성공 경로와 `no_phi_attestation=false` 실패 경로를 검증한다.
- `scripts/flutter_web_timing_review_fixture_validate.sh`
  - `manifest.env`, five sample files, `aggregate.env`, allowed review status, GitHub Actions artifact source, no-PHI attestation을 검증한다.
  - sample/aggregate/manifest에서 sensitive field keys를 차단한다.
  - P32 aggregator를 재실행해 `aggregate.env` 재현성을 검증한다.
- `scripts/ci_web_timing_review_fixture_policy_test.sh`
  - fixture policy가 validator script/test, manifest required fields, forbidden fixture fields를 문서화했는지 검증한다.
- `docs/ci/flutter-web-timing-review-fixtures.md`
  - validator 사용법과 검증 범위를 문서화했다.
- `docs/superpowers/plans/2026-05-14-flutter-web-timing-review-validator-p34.md`
  - P34 RED-GREEN 구현 순서와 검증 명령을 기록했다.

## TDD Record

- RED: `bash scripts/flutter_web_timing_review_fixture_validate_test.sh`
  - 실패: `scripts/flutter_web_timing_review_fixture_validate.sh` missing
- GREEN: `bash scripts/flutter_web_timing_review_fixture_validate_test.sh`
  - 통과: `FLUTTER_WEB_TIMING_REVIEW_FIXTURE_VALIDATE_TEST_PASS`
- RED: `bash scripts/ci_web_timing_review_fixture_policy_test.sh`
  - 실패: `policy missing marker: fixture_validator_script: scripts/flutter_web_timing_review_fixture_validate.sh`
- GREEN: `bash scripts/ci_web_timing_review_fixture_policy_test.sh`
  - 통과: `CI_WEB_TIMING_REVIEW_FIXTURE_POLICY_PASS`

## Verification

- `bash scripts/flutter_web_timing_review_fixture_validate_test.sh`: PASS
- `bash scripts/ci_web_timing_review_fixture_policy_test.sh`: PASS
- `bash scripts/ci_web_timing_threshold_promotion_policy_test.sh`: PASS
- `bash scripts/ci_gate_matrix_policy_test.sh`: PASS
- `bash -n scripts/flutter_web_timing_review_fixture_validate.sh scripts/flutter_web_timing_review_fixture_validate_test.sh scripts/ci_web_timing_review_fixture_policy_test.sh scripts/ci_web_timing_threshold_promotion_policy_test.sh scripts/ci_gate_matrix_policy_test.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS

## Self Review

- Validator keeps real review fixture validation offline and deterministic.
- Aggregate reproducibility is checked with the same aggregator used by the promotion policy.
- Forbidden field scan blocks common PHI, identifier, token, and raw measurement keys before fixture review.

## Residual Notes

- P35 can add a small CLI wrapper or documented PR checklist that points reviewers to the validator command before threshold policy edits.
- No real timing fixture directory was added in this step.
