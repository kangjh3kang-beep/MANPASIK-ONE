# Flutter Web Release Review Checklist P36 Audit

## Scope

P36은 Flutter web timing threshold 또는 PR gate 정책 변경 전에 reviewer가 따라야 할 release-review command order를 한 문서로 고정한 단계다.
Checklist policy guard가 checklist marker, 필수 command, 참조 문서/스크립트, workflow, matrix 연결을 검증한다.

## Changes

- `scripts/ci_web_release_review_checklist_policy_test.sh`
  - checklist 문서, 필수 command, 참조 파일, workflow step, matrix row를 검증한다.
- `docs/ci/flutter-web-release-review-checklist.md`
  - policy gates, artifact export, offline aggregation, fixture validation, threshold review check, audit decision 순서를 문서화했다.
- `.github/workflows/ci.yml`
  - `ssot-governance` job에 `Flutter web release review checklist policy` step을 추가했다.
- `scripts/ci_gate_matrix_policy_test.sh`
  - matrix가 checklist policy doc/guard/workflow를 참조하도록 검증한다.
- `docs/ci/ci-gate-matrix.md`
  - `Flutter web release review checklist policy` row와 checklist marker를 추가했다.
- `docs/superpowers/plans/2026-05-19-flutter-web-release-review-checklist-p36.md`
  - P36 RED-GREEN 구현 순서와 검증 명령을 기록했다.

## TDD Record

- RED: `bash scripts/ci_web_release_review_checklist_policy_test.sh`
  - 실패: `missing file: docs/ci/flutter-web-release-review-checklist.md`
- RED after checklist/workflow: `bash scripts/ci_web_release_review_checklist_policy_test.sh`
  - 실패: `matrix missing marker: Flutter web release review checklist policy`
- RED: `bash scripts/ci_gate_matrix_policy_test.sh`
  - 실패: `missing matrix marker: Flutter web release review checklist policy`
- GREEN: `bash scripts/ci_web_release_review_checklist_policy_test.sh`
  - 통과: `CI_WEB_RELEASE_REVIEW_CHECKLIST_POLICY_PASS`
- GREEN: `bash scripts/ci_gate_matrix_policy_test.sh`
  - 통과: `CI_GATE_MATRIX_POLICY_PASS`

## Verification

- `bash scripts/ci_web_release_review_checklist_policy_test.sh`: PASS
- `bash scripts/ci_web_timing_review_fixture_policy_test.sh`: PASS
- `bash scripts/ci_web_timing_threshold_promotion_policy_test.sh`: PASS
- `bash scripts/ci_gate_matrix_policy_test.sh`: PASS
- `bash -n scripts/ci_web_release_review_checklist_policy_test.sh scripts/ci_web_timing_review_fixture_policy_test.sh scripts/ci_web_timing_threshold_promotion_policy_test.sh scripts/ci_gate_matrix_policy_test.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS

## Self Review

- Checklist references the current command names from P30-P35 and is guarded by CI.
- The checklist keeps PHI/no-identifier review requirements visible before timing policy changes.
- Matrix and workflow now include the checklist policy alongside timing collection, promotion, and fixture gates.

## Residual Notes

- P37 can add an optional concrete review example with synthetic timing data if a non-production example is needed.
- No real timing artifact fixture was added in this step.
