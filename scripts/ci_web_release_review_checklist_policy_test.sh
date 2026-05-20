#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CHECKLIST="$ROOT/docs/ci/flutter-web-release-review-checklist.md"
MATRIX="$ROOT/docs/ci/ci-gate-matrix.md"
WORKFLOW="$ROOT/.github/workflows/ci.yml"

fail() {
  echo "CI_WEB_RELEASE_REVIEW_CHECKLIST_POLICY_FAIL $1" >&2
  exit 1
}

for file in "$CHECKLIST" "$MATRIX" "$WORKFLOW"; do
  if [[ ! -f "$file" ]]; then
    fail "missing file: ${file#$ROOT/}"
  fi
done

required_markers=(
  "release_review_checklist_version: 1"
  "review_scope: flutter_web_timing_release_review"
  "review_order: policy_gates,artifact_collection,offline_aggregation,fixture_validation,threshold_review_check,decision_record"
  "required_before_threshold_change: true"
  "no_phi_review_required: true"
  "checklist_guard: scripts/ci_web_release_review_checklist_policy_test.sh"
  "synthetic_review_example: docs/ci/fixtures/flutter-web-timing/2026-05-19-synthetic-review"
  "synthetic_review_audit: docs/audit/flutter-web-timing-synthetic-review-example.md"
)

for marker in "${required_markers[@]}"; do
  if ! grep -Fq "$marker" "$CHECKLIST"; then
    fail "checklist missing marker: $marker"
  fi
done

required_commands=(
  "bash scripts/ci_web_gate_timing_collection_policy_test.sh"
  "bash scripts/ci_web_timing_threshold_promotion_policy_test.sh"
  "bash scripts/ci_web_timing_review_fixture_policy_test.sh"
  "bash scripts/flutter_web_timing_sample_aggregate.sh --input-dir docs/ci/fixtures/flutter-web-timing/<review-id>/samples --output docs/ci/fixtures/flutter-web-timing/<review-id>/aggregate.env"
  "bash scripts/flutter_web_timing_review_fixture_validate.sh --review-dir docs/ci/fixtures/flutter-web-timing/<review-id>"
  "bash scripts/flutter_web_timing_threshold_change_review_check.sh --review-dir docs/ci/fixtures/flutter-web-timing/<review-id> --audit-file docs/audit/<review-doc>.md"
  "bash scripts/ci_web_timing_synthetic_review_example_test.sh"
)

for command in "${required_commands[@]}"; do
  if ! grep -Fq "$command" "$CHECKLIST"; then
    fail "checklist missing command: $command"
  fi
done

required_files=(
  "docs/ci/flutter-web-gate-timing-collection.md"
  "docs/ci/flutter-web-timing-threshold-promotion.md"
  "docs/ci/flutter-web-timing-review-fixtures.md"
  "docs/audit/templates/flutter-web-timing-review-template.md"
  "scripts/flutter_web_timing_sample_aggregate.sh"
  "scripts/flutter_web_timing_review_fixture_validate.sh"
  "scripts/flutter_web_timing_threshold_change_review_check.sh"
  "scripts/ci_web_timing_synthetic_review_example_test.sh"
  "docs/audit/flutter-web-timing-synthetic-review-example.md"
)

for path in "${required_files[@]}"; do
  if [[ ! -f "$ROOT/$path" ]]; then
    fail "missing referenced file: $path"
  fi
  if ! grep -Fq "$path" "$CHECKLIST"; then
    fail "checklist does not reference file: $path"
  fi
done

bash "$ROOT/scripts/ci_web_timing_synthetic_review_example_test.sh" >/dev/null

workflow_markers=(
  "Flutter web release review checklist policy"
  "scripts/ci_web_release_review_checklist_policy_test.sh"
)

for marker in "${workflow_markers[@]}"; do
  if ! grep -Fq "$marker" "$WORKFLOW"; then
    fail "workflow missing marker: $marker"
  fi
done

matrix_markers=(
  "Flutter web release review checklist policy"
  "docs/ci/flutter-web-release-review-checklist.md"
  "scripts/ci_web_release_review_checklist_policy_test.sh"
)

for marker in "${matrix_markers[@]}"; do
  if ! grep -Fq "$marker" "$MATRIX"; then
    fail "matrix missing marker: $marker"
  fi
done

echo "CI_WEB_RELEASE_REVIEW_CHECKLIST_POLICY_PASS"
