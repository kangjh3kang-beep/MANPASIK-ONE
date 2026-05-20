#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
POLICY="$ROOT/docs/ci/flutter-web-timing-threshold-promotion.md"
COLLECTION_POLICY="$ROOT/docs/ci/flutter-web-gate-timing-collection.md"
MATRIX="$ROOT/docs/ci/ci-gate-matrix.md"
WORKFLOW="$ROOT/.github/workflows/ci.yml"

fail() {
  echo "CI_WEB_TIMING_THRESHOLD_PROMOTION_POLICY_FAIL $1" >&2
  exit 1
}

for file in "$POLICY" "$COLLECTION_POLICY" "$MATRIX" "$WORKFLOW"; do
  if [[ ! -f "$file" ]]; then
    fail "missing file: ${file#$ROOT/}"
  fi
done

required_policy_markers=(
  "promotion_policy_version: 1"
  "sample_source_artifact: flutter-web-timing-report"
  "minimum_successful_artifact_samples: 5"
  "required_branch_types: pull_request,release_branch"
  "required_runner_context_field: runner_context"
  "required_duration_field: latest_duration_seconds"
  "required_distribution_fields: min_duration_seconds,max_duration_seconds,duration_seconds_values"
  "required_statistics: median_seconds,p95_seconds,worst_case_seconds"
  "threshold_change_requires: measured_cost_review"
  "automated_gate_relaxation: false"
  "relaxation_allowed_before_minimum_samples: false"
  "nightly_split_change_requires: docs/ci/flutter-web-release-gate-policy.md"
  "review_record_required: docs/audit"
  "sample_aggregator_script: scripts/flutter_web_timing_sample_aggregate.sh"
  "sample_aggregator_test: scripts/flutter_web_timing_sample_aggregate_test.sh"
  "aggregation_output_fields: aggregate_version,source_artifact,sample_count,duration_seconds_values,median_seconds,p95_seconds,worst_case_seconds,branch_types,runner_contexts"
)

for marker in "${required_policy_markers[@]}"; do
  if ! grep -Fq "$marker" "$POLICY"; then
    fail "missing promotion marker: $marker"
  fi
done

for path in scripts/flutter_web_timing_sample_aggregate.sh scripts/flutter_web_timing_sample_aggregate_test.sh; do
  if [[ ! -f "$ROOT/$path" ]]; then
    fail "missing aggregator file: $path"
  fi
done

cross_doc_markers=(
  "threshold_promotion_policy: docs/ci/flutter-web-timing-threshold-promotion.md"
  "sample_source_artifact: flutter-web-timing-report"
)

for marker in "${cross_doc_markers[@]}"; do
  if ! grep -Fq "$marker" "$COLLECTION_POLICY"; then
    fail "collection policy missing marker: $marker"
  fi
done

workflow_markers=(
  "Flutter web timing threshold promotion policy"
  "scripts/ci_web_timing_threshold_promotion_policy_test.sh"
)

for marker in "${workflow_markers[@]}"; do
  if ! grep -Fq "$marker" "$WORKFLOW"; then
    fail "workflow missing marker: $marker"
  fi
done

matrix_markers=(
  "Flutter web timing threshold promotion policy"
  "docs/ci/flutter-web-timing-threshold-promotion.md"
  "scripts/ci_web_timing_threshold_promotion_policy_test.sh"
)

for marker in "${matrix_markers[@]}"; do
  if ! grep -Fq "$marker" "$MATRIX"; then
    fail "matrix missing marker: $marker"
  fi
done

echo "CI_WEB_TIMING_THRESHOLD_PROMOTION_POLICY_PASS"
