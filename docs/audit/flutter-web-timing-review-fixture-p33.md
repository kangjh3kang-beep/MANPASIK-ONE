# Flutter Web Timing Review Fixture P33 Audit

## Scope

P33은 실제 `flutter-web-timing-report` 5-sample measured cost review를 재현 가능한 fixture와 audit template으로 남기기 위한 규약 단계다.
P32 aggregator가 만든 `aggregate.env`를 정책 변경 전 review evidence로 기록할 수 있게 폴더 구조, manifest, template, CI guard를 고정했다.

## Changes

- `scripts/ci_web_timing_review_fixture_policy_test.sh`
  - fixture convention 문서, fixture root README, audit template, promotion policy cross-reference, workflow, matrix 연결을 검증한다.
- `docs/ci/flutter-web-timing-review-fixtures.md`
  - fixture root, manifest, sample file naming, aggregate output, artifact source, no-PHI 규칙을 문서화했다.
- `docs/ci/fixtures/flutter-web-timing/README.md`
  - 실제 review directory 구조와 데이터 반입 제한을 기록했다.
- `docs/audit/templates/flutter-web-timing-review-template.md`
  - aggregate results, decision, sign-off, rollback notes를 기록하는 audit template을 추가했다.
- `docs/ci/flutter-web-timing-threshold-promotion.md`
  - review fixture policy, fixture root, audit template을 promotion policy에 교차 참조했다.
- `.github/workflows/ci.yml`
  - `ssot-governance` job에 `Flutter web timing review fixture policy` step을 추가했다.
- `scripts/ci_gate_matrix_policy_test.sh`
  - matrix가 review fixture policy doc/template/guard/workflow를 참조하도록 검증한다.
- `docs/ci/ci-gate-matrix.md`
  - `Flutter web timing review fixture policy` row와 fixture marker를 추가했다.
- `docs/superpowers/plans/2026-05-14-flutter-web-timing-review-fixture-p33.md`
  - P33 RED-GREEN 구현 순서와 검증 명령을 기록했다.

## TDD Record

- RED: `bash scripts/ci_web_timing_review_fixture_policy_test.sh`
  - 실패: `missing file: docs/ci/flutter-web-timing-review-fixtures.md`
- RED after docs/workflow: `bash scripts/ci_web_timing_review_fixture_policy_test.sh`
  - 실패: `matrix missing marker: Flutter web timing review fixture policy`
- RED: `bash scripts/ci_gate_matrix_policy_test.sh`
  - 실패: `missing matrix marker: Flutter web timing review fixture policy`
- GREEN: `bash scripts/ci_web_timing_review_fixture_policy_test.sh`
  - 통과: `CI_WEB_TIMING_REVIEW_FIXTURE_POLICY_PASS`
- GREEN: `bash scripts/ci_gate_matrix_policy_test.sh`
  - 통과: `CI_GATE_MATRIX_POLICY_PASS`

## Verification

- `bash scripts/ci_web_timing_review_fixture_policy_test.sh`: PASS
- `bash scripts/ci_web_timing_threshold_promotion_policy_test.sh`: PASS
- `bash scripts/flutter_web_timing_sample_aggregate_test.sh`: PASS
- `bash scripts/ci_gate_matrix_policy_test.sh`: PASS
- `bash -n scripts/ci_web_timing_review_fixture_policy_test.sh scripts/ci_web_timing_threshold_promotion_policy_test.sh scripts/flutter_web_timing_sample_aggregate.sh scripts/flutter_web_timing_sample_aggregate_test.sh scripts/ci_gate_matrix_policy_test.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS

## Self Review

- Fixture convention requires GitHub Actions artifacts only and explicitly forbids PHI, identifiers, tokens, and raw medical measurements.
- Audit template records decision, aggregate statistics, sign-off, and rollback notes so threshold changes remain reviewable.
- Matrix and workflow now run the fixture policy guard alongside the threshold promotion guard.

## Residual Notes

- P34 can add an optional fixture manifest validator that checks a concrete `<review-id>` directory before a threshold change PR.
- No real CI timing artifacts were added in this step.
