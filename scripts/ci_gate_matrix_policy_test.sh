#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MATRIX="$ROOT/docs/ci/ci-gate-matrix.md"
WORKFLOW="$ROOT/.github/workflows/ci.yml"

fail() {
  echo "CI_GATE_MATRIX_POLICY_FAIL $1" >&2
  exit 1
}

if [[ ! -f "$MATRIX" ]]; then
  fail "missing matrix document"
fi

required_markers=(
  "matrix_version: 1"
  "ssot-governance"
  "Flutter evidence UI gate: contract"
  "Flutter evidence UI gate: ui"
  "Flutter evidence UI gate: data"
  "Flutter web release gate"
  "Flutter web timing collection policy"
  "Wasm readiness policy"
  "Security release gate"
  "Assay evidence gate"
  "blocking: true"
  "blocking: false"
  "execution_mode: category_shards"
  "execution_mode: pull_request_and_release_branch"
  "release_target: js_web"
  "wasm_release_target: false"
  "timing_capture: true"
  "timing_marker: FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS"
  "minimum_samples_before_relaxing: 5"
  "threshold_policy_status: advisory"
  "Flutter web timing artifact"
  "artifact_name: flutter-web-timing-report"
  "artifact_path: /tmp/manpasik_flutter_web_timing.env"
  "actions/upload-artifact@v4"
  "Flutter web timing threshold promotion policy"
  "promotion_policy_version: 1"
  "automated_gate_relaxation: false"
  "Flutter web timing review fixture policy"
  "review_fixture_policy_version: 1"
  "fixture_root: docs/ci/fixtures/flutter-web-timing"
  "Flutter web release review checklist policy"
  "release_review_checklist_version: 1"
  "review_scope: flutter_web_timing_release_review"
  "Flutter web timing governance index policy"
  "governance_index_version: 1"
  "governance_scope: flutter_web_timing_ci"
)

for marker in "${required_markers[@]}"; do
  if ! grep -Fq "$marker" "$MATRIX"; then
    fail "missing matrix marker: $marker"
  fi
done

required_files=(
  "docs/ci/flutter-evidence-ui-gate-policy.md"
  "docs/ci/flutter-web-release-gate-policy.md"
  "docs/ci/flutter-web-gate-timing-collection.md"
  "docs/ci/flutter-web-timing-threshold-promotion.md"
  "docs/ci/flutter-web-timing-review-fixtures.md"
  "docs/ci/flutter-web-release-review-checklist.md"
  "docs/ci/flutter-web-timing-governance-index.md"
  "docs/audit/templates/flutter-web-timing-review-template.md"
  "docs/ci/flutter-wasm-readiness.md"
  "scripts/flutter_evidence_ui_gate_policy_test.sh"
  "scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh"
  "scripts/ci_web_gate_policy_test.sh"
  "scripts/ci_web_gate_timing_collection_policy_test.sh"
  "scripts/ci_flutter_web_timing_artifact_workflow_test.sh"
  "scripts/ci_web_timing_threshold_promotion_policy_test.sh"
  "scripts/ci_web_timing_review_fixture_policy_test.sh"
  "scripts/ci_web_release_review_checklist_policy_test.sh"
  "scripts/ci_web_timing_governance_index_policy_test.sh"
  "scripts/ci_wasm_readiness_policy_test.sh"
  "scripts/security_release_gate.sh"
  "scripts/assay_evidence_gate.sh"
)

for path in "${required_files[@]}"; do
  if [[ ! -f "$ROOT/$path" ]]; then
    fail "missing referenced file: $path"
  fi
  if ! grep -Fq "$path" "$MATRIX"; then
    fail "matrix does not reference file: $path"
  fi
done

workflow_markers=(
  "SSOT Governance"
  "Flutter evidence UI gate: contract"
  "Flutter evidence UI gate: ui"
  "Flutter evidence UI gate: data"
  "Flutter web release gate"
  "Flutter web timing collection policy"
  "Flutter web timing threshold promotion policy"
  "Flutter web timing review fixture policy"
  "Flutter web release review checklist policy"
  "Flutter web timing governance index policy"
)

for marker in "${workflow_markers[@]}"; do
  if ! grep -Fq "$marker" "$WORKFLOW"; then
    fail "workflow missing marker: $marker"
  fi
done

echo "CI_GATE_MATRIX_POLICY_PASS"
