#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
POLICY="$ROOT/docs/ci/flutter-web-timing-review-fixtures.md"
FIXTURE_README="$ROOT/docs/ci/fixtures/flutter-web-timing/README.md"
TEMPLATE="$ROOT/docs/audit/templates/flutter-web-timing-review-template.md"
PROMOTION_POLICY="$ROOT/docs/ci/flutter-web-timing-threshold-promotion.md"
MATRIX="$ROOT/docs/ci/ci-gate-matrix.md"
WORKFLOW="$ROOT/.github/workflows/ci.yml"

fail() {
  echo "CI_WEB_TIMING_REVIEW_FIXTURE_POLICY_FAIL $1" >&2
  exit 1
}

for file in "$POLICY" "$FIXTURE_README" "$TEMPLATE" "$PROMOTION_POLICY" "$MATRIX" "$WORKFLOW"; do
  if [[ ! -f "$file" ]]; then
    fail "missing file: ${file#$ROOT/}"
  fi
done

required_policy_markers=(
  "review_fixture_policy_version: 1"
  "fixture_root: docs/ci/fixtures/flutter-web-timing"
  "fixture_manifest: manifest.env"
  "required_exported_artifacts: 5"
  "required_sample_files: sample-01.env,sample-02.env,sample-03.env,sample-04.env,sample-05.env"
  "aggregate_output_file: aggregate.env"
  "aggregator_script: scripts/flutter_web_timing_sample_aggregate.sh"
  "audit_template: docs/audit/templates/flutter-web-timing-review-template.md"
  "promotion_policy: docs/ci/flutter-web-timing-threshold-promotion.md"
  "review_status_allowed: proposed,approved,rejected"
  "artifact_source_required: github_actions"
  "no_phi_allowed: true"
  "fixture_validator_script: scripts/flutter_web_timing_review_fixture_validate.sh"
  "fixture_validator_test: scripts/flutter_web_timing_review_fixture_validate_test.sh"
  "manifest_required_fields: review_id,artifact_source,artifact_name,sample_count,sample_files,aggregate_output,review_status,no_phi_attestation"
  "forbidden_fixture_fields: user_id,patient_id,device_id,access_token,refresh_token,raw_channels,s_det,s_ref,primary_value"
  "threshold_change_review_check_script: scripts/flutter_web_timing_threshold_change_review_check.sh"
  "threshold_change_review_check_test: scripts/flutter_web_timing_threshold_change_review_check_test.sh"
  "threshold_change_review_check_required: true"
  "threshold_change_review_required_fields: review_id,decision,sample_count,median_seconds,p95_seconds,worst_case_seconds,branch_types,runner_contexts,validator_command"
)

for marker in "${required_policy_markers[@]}"; do
  if ! grep -Fq "$marker" "$POLICY"; then
    fail "policy missing marker: $marker"
  fi
done

for path in scripts/flutter_web_timing_review_fixture_validate.sh scripts/flutter_web_timing_review_fixture_validate_test.sh scripts/flutter_web_timing_threshold_change_review_check.sh scripts/flutter_web_timing_threshold_change_review_check_test.sh; do
  if [[ ! -f "$ROOT/$path" ]]; then
    fail "missing validator file: $path"
  fi
done

template_markers=(
  "review_template_version: 1"
  "review_fixture_policy: docs/ci/flutter-web-timing-review-fixtures.md"
  "aggregate_source: aggregate.env"
  "required_decision_fields: decision,median_seconds,p95_seconds,worst_case_seconds,sample_count,branch_types,runner_contexts"
  "validator_command:"
)

for marker in "${template_markers[@]}"; do
  if ! grep -Fq "$marker" "$TEMPLATE"; then
    fail "template missing marker: $marker"
  fi
done

cross_doc_markers=(
  "review_fixture_policy: docs/ci/flutter-web-timing-review-fixtures.md"
  "review_fixture_root: docs/ci/fixtures/flutter-web-timing"
  "review_audit_template: docs/audit/templates/flutter-web-timing-review-template.md"
)

for marker in "${cross_doc_markers[@]}"; do
  if ! grep -Fq "$marker" "$PROMOTION_POLICY"; then
    fail "promotion policy missing marker: $marker"
  fi
done

workflow_markers=(
  "Flutter web timing review fixture policy"
  "scripts/ci_web_timing_review_fixture_policy_test.sh"
)

for marker in "${workflow_markers[@]}"; do
  if ! grep -Fq "$marker" "$WORKFLOW"; then
    fail "workflow missing marker: $marker"
  fi
done

matrix_markers=(
  "Flutter web timing review fixture policy"
  "docs/ci/flutter-web-timing-review-fixtures.md"
  "scripts/ci_web_timing_review_fixture_policy_test.sh"
)

for marker in "${matrix_markers[@]}"; do
  if ! grep -Fq "$marker" "$MATRIX"; then
    fail "matrix missing marker: $marker"
  fi
done

echo "CI_WEB_TIMING_REVIEW_FIXTURE_POLICY_PASS"
