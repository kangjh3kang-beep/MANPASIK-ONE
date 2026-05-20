# Flutter Web Timing Threshold Promotion P31 Audit

## Scope

P31은 Flutter web release gate timing artifact가 쌓인 뒤 PR 필수 web gate를 완화할 수 있는 조건을 정책으로 고정한 단계다.
목표는 최소 sample 수, review 입력, 통계 항목, 자동 완화 금지, 관련 문서 변경 조건을 CI에서 검증 가능하게 만드는 것이다.

## Changes

- `scripts/ci_web_timing_threshold_promotion_policy_test.sh`
  - promotion policy 문서, timing collection 문서, workflow, matrix 연결을 검증한다.
- `docs/ci/flutter-web-timing-threshold-promotion.md`
  - 최소 5개 성공 artifact sample, branch type, runner context, median/p95/worst-case review, 자동 완화 금지를 문서화했다.
- `.github/workflows/ci.yml`
  - `ssot-governance` job에 `Flutter web timing threshold promotion policy` step을 추가했다.
- `docs/ci/flutter-web-gate-timing-collection.md`
  - threshold promotion policy 문서와 `flutter-web-timing-report` sample source를 교차 참조한다.
- `scripts/ci_gate_matrix_policy_test.sh`
  - promotion policy marker, policy 문서, guard script, workflow step 연결을 검증한다.
- `docs/ci/ci-gate-matrix.md`
  - `Flutter web timing threshold promotion policy` row와 promotion marker를 추가했다.
- `docs/superpowers/plans/2026-05-14-flutter-web-timing-threshold-promotion-p31.md`
  - P31 RED-GREEN 구현 순서와 검증 명령을 기록했다.

## TDD Record

- RED: `bash scripts/ci_web_timing_threshold_promotion_policy_test.sh`
  - 실패: `missing file: docs/ci/flutter-web-timing-threshold-promotion.md`
- RED after policy/workflow: `bash scripts/ci_web_timing_threshold_promotion_policy_test.sh`
  - 실패: `matrix missing marker: Flutter web timing threshold promotion policy`
- RED: `bash scripts/ci_gate_matrix_policy_test.sh`
  - 실패: `missing matrix marker: Flutter web timing threshold promotion policy`
- GREEN: `bash scripts/ci_web_timing_threshold_promotion_policy_test.sh`
  - 통과: `CI_WEB_TIMING_THRESHOLD_PROMOTION_POLICY_PASS`
- GREEN: `bash scripts/ci_gate_matrix_policy_test.sh`
  - 통과: `CI_GATE_MATRIX_POLICY_PASS`

## Verification

- `bash scripts/ci_web_timing_threshold_promotion_policy_test.sh`: PASS
- `bash scripts/ci_web_gate_timing_collection_policy_test.sh`: PASS
- `bash scripts/ci_gate_matrix_policy_test.sh`: PASS
- `bash scripts/ci_flutter_web_timing_artifact_workflow_test.sh`: PASS
- `bash -n scripts/ci_web_timing_threshold_promotion_policy_test.sh scripts/ci_flutter_web_timing_artifact_workflow_test.sh scripts/ci_web_gate_timing_collection_policy_test.sh scripts/ci_gate_matrix_policy_test.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS

## Self Review

- Promotion is policy-only and deliberately does not parse artifact files yet, because actual CI artifacts are not present in the repository.
- The policy forbids automatic relaxation, which protects the PR web release gate until measured cost review is recorded.
- The collection policy now points to the promotion policy so sample collection and threshold decisions do not drift.

## Residual Notes

- P32 can add an offline sample aggregator once at least 5 CI artifact files are exported into a review folder.
- The full Flutter web build was not re-run in this step because the change is governance wiring and policy validation.
