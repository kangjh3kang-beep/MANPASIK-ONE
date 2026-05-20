# Flutter Web Timing Governance Index

```yaml
governance_index_version: 1
governance_scope: flutter_web_timing_ci
entrypoint_status: canonical
policy_chain: release_gate,timing_collection,artifact_upload,promotion_policy,review_fixture,release_review_checklist,synthetic_example
guard_chain: release_gate_policy,release_gate_workflow,timing_collection_policy,timing_artifact_workflow,threshold_promotion_policy,review_fixture_policy,release_review_checklist_policy,synthetic_review_example,governance_index_policy,gate_matrix_policy
index_guard: scripts/ci_web_timing_governance_index_policy_test.sh
```

## Purpose

This is the canonical entrypoint for Flutter web timing CI governance.
Start here before changing Flutter web timing thresholds, PR web release gate policy, or release review evidence requirements.

## Policy Chain

1. Release gate:
   - `docs/ci/flutter-web-release-gate-policy.md`
   - `scripts/flutter_web_release_gate.sh`
2. Timing collection and report artifact:
   - `docs/ci/flutter-web-gate-timing-collection.md`
   - `scripts/flutter_web_timing_report.sh`
3. Threshold promotion:
   - `docs/ci/flutter-web-timing-threshold-promotion.md`
   - `scripts/flutter_web_timing_sample_aggregate.sh`
4. Review fixtures:
   - `docs/ci/flutter-web-timing-review-fixtures.md`
   - `docs/audit/templates/flutter-web-timing-review-template.md`
   - `scripts/flutter_web_timing_review_fixture_validate.sh`
   - `scripts/flutter_web_timing_threshold_change_review_check.sh`
5. Release review checklist:
   - `docs/ci/flutter-web-release-review-checklist.md`
6. Synthetic example:
   - `scripts/ci_web_timing_synthetic_review_example_test.sh`

## Guard Chain

Run these guards before changing timing thresholds, release-gate timing evidence, or review fixture requirements:

1. Release gate policy and workflow:
   - `scripts/ci_web_gate_policy_test.sh`
   - `scripts/ci_flutter_web_gate_workflow_test.sh`
   - `scripts/flutter_web_release_gate_test.sh`
2. Timing collection and artifact export:
   - `scripts/ci_web_gate_timing_collection_policy_test.sh`
   - `scripts/ci_flutter_web_timing_artifact_workflow_test.sh`
   - `scripts/flutter_web_timing_report_test.sh`
3. Promotion, fixture, and reviewer checks:
   - `scripts/ci_web_timing_threshold_promotion_policy_test.sh`
   - `scripts/flutter_web_timing_sample_aggregate_test.sh`
   - `scripts/ci_web_timing_review_fixture_policy_test.sh`
   - `scripts/flutter_web_timing_review_fixture_validate_test.sh`
   - `scripts/flutter_web_timing_threshold_change_review_check_test.sh`
4. Release review and governance:
   - `scripts/ci_web_release_review_checklist_policy_test.sh`
   - `scripts/ci_web_timing_synthetic_review_example_test.sh`
   - `scripts/ci_web_timing_governance_index_policy_test.sh`
   - `scripts/ci_gate_matrix_policy_test.sh`

## Governance Matrix

Use `docs/ci/ci-gate-matrix.md` to confirm workflow step, policy doc, guard script, blocking status, and execution mode.

## Audit Trail

- `docs/audit/flutter-web-release-gate-timing-p27.md`
- `docs/audit/flutter-web-timing-collection-policy-p28.md`
- `docs/audit/flutter-web-timing-report-p29.md`
- `docs/audit/flutter-web-timing-artifact-p30.md`
- `docs/audit/flutter-web-timing-threshold-promotion-p31.md`
- `docs/audit/flutter-web-timing-sample-aggregator-p32.md`
- `docs/audit/flutter-web-timing-review-fixture-p33.md`
- `docs/audit/flutter-web-timing-review-validator-p34.md`
- `docs/audit/flutter-web-timing-threshold-change-review-check-p35.md`
- `docs/audit/flutter-web-release-review-checklist-p36.md`
- `docs/audit/flutter-web-timing-synthetic-review-example-p37.md`
- `docs/audit/flutter-web-timing-governance-index-p38.md`

## Change Rules

- Do not relax the PR web release gate from this index alone; use the release review checklist and audit template.
- Do not treat synthetic examples as measured production evidence.
- Do not remove any referenced policy, script, or audit template without updating this index and its guard in the same change.
