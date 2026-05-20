#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
INDEX="$ROOT/docs/ci/flutter-web-timing-governance-index.md"
MATRIX="$ROOT/docs/ci/ci-gate-matrix.md"
WORKFLOW="$ROOT/.github/workflows/ci.yml"

fail() {
  echo "CI_WEB_TIMING_GOVERNANCE_INDEX_POLICY_FAIL $1" >&2
  exit 1
}

for file in "$INDEX" "$MATRIX" "$WORKFLOW"; do
  if [[ ! -f "$file" ]]; then
    fail "missing file: ${file#$ROOT/}"
  fi
done

required_markers=(
  "governance_index_version: 1"
  "governance_scope: flutter_web_timing_ci"
  "entrypoint_status: canonical"
  "policy_chain: release_gate,timing_collection,artifact_upload,promotion_policy,review_fixture,release_review_checklist,synthetic_example"
  "guard_chain: release_gate_policy,release_gate_workflow,timing_collection_policy,timing_artifact_workflow,threshold_promotion_policy,review_fixture_policy,release_review_checklist_policy,synthetic_review_example,governance_index_policy,gate_matrix_policy"
  "index_guard: scripts/ci_web_timing_governance_index_policy_test.sh"
)

for marker in "${required_markers[@]}"; do
  if ! grep -Fq "$marker" "$INDEX"; then
    fail "index missing marker: $marker"
  fi
done

required_refs=(
  "docs/ci/flutter-web-release-gate-policy.md"
  "docs/ci/flutter-web-gate-timing-collection.md"
  "docs/ci/flutter-web-timing-threshold-promotion.md"
  "docs/ci/flutter-web-timing-review-fixtures.md"
  "docs/ci/flutter-web-release-review-checklist.md"
  "docs/ci/ci-gate-matrix.md"
  "docs/audit/templates/flutter-web-timing-review-template.md"
  "docs/audit/flutter-web-release-gate-timing-p27.md"
  "docs/audit/flutter-web-timing-collection-policy-p28.md"
  "docs/audit/flutter-web-timing-report-p29.md"
  "docs/audit/flutter-web-timing-artifact-p30.md"
  "docs/audit/flutter-web-timing-threshold-promotion-p31.md"
  "docs/audit/flutter-web-timing-sample-aggregator-p32.md"
  "docs/audit/flutter-web-timing-review-fixture-p33.md"
  "docs/audit/flutter-web-timing-review-validator-p34.md"
  "docs/audit/flutter-web-timing-threshold-change-review-check-p35.md"
  "docs/audit/flutter-web-release-review-checklist-p36.md"
  "docs/audit/flutter-web-timing-synthetic-review-example-p37.md"
  "docs/audit/flutter-web-timing-governance-index-p38.md"
  "scripts/ci_web_gate_policy_test.sh"
  "scripts/ci_flutter_web_gate_workflow_test.sh"
  "scripts/ci_web_gate_timing_collection_policy_test.sh"
  "scripts/ci_flutter_web_timing_artifact_workflow_test.sh"
  "scripts/ci_web_timing_threshold_promotion_policy_test.sh"
  "scripts/ci_web_timing_review_fixture_policy_test.sh"
  "scripts/ci_web_release_review_checklist_policy_test.sh"
  "scripts/ci_gate_matrix_policy_test.sh"
  "scripts/flutter_web_release_gate.sh"
  "scripts/flutter_web_release_gate_test.sh"
  "scripts/flutter_web_timing_report.sh"
  "scripts/flutter_web_timing_report_test.sh"
  "scripts/flutter_web_timing_sample_aggregate.sh"
  "scripts/flutter_web_timing_sample_aggregate_test.sh"
  "scripts/flutter_web_timing_review_fixture_validate.sh"
  "scripts/flutter_web_timing_review_fixture_validate_test.sh"
  "scripts/flutter_web_timing_threshold_change_review_check.sh"
  "scripts/flutter_web_timing_threshold_change_review_check_test.sh"
  "scripts/ci_web_timing_synthetic_review_example_test.sh"
)

for ref in "${required_refs[@]}"; do
  if [[ ! -f "$ROOT/$ref" ]]; then
    fail "missing referenced file: $ref"
  fi
  if ! grep -Fq "$ref" "$INDEX"; then
    fail "index does not reference file: $ref"
  fi
done

workflow_markers=(
  "Flutter web timing governance index policy"
  "scripts/ci_web_timing_governance_index_policy_test.sh"
)

for marker in "${workflow_markers[@]}"; do
  if ! grep -Fq "$marker" "$WORKFLOW"; then
    fail "workflow missing marker: $marker"
  fi
done

matrix_markers=(
  "Flutter web timing governance index policy"
  "docs/ci/flutter-web-timing-governance-index.md"
  "scripts/ci_web_timing_governance_index_policy_test.sh"
)

for marker in "${matrix_markers[@]}"; do
  if ! grep -Fq "$marker" "$MATRIX"; then
    fail "matrix missing marker: $marker"
  fi
done

echo "CI_WEB_TIMING_GOVERNANCE_INDEX_POLICY_PASS"
