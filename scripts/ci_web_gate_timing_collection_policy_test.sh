#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
POLICY="$ROOT/docs/ci/flutter-web-gate-timing-collection.md"
WEB_POLICY="$ROOT/docs/ci/flutter-web-release-gate-policy.md"
MATRIX="$ROOT/docs/ci/ci-gate-matrix.md"
WORKFLOW="$ROOT/.github/workflows/ci.yml"

fail() {
  echo "CI_WEB_GATE_TIMING_COLLECTION_POLICY_FAIL $1" >&2
  exit 1
}

for file in "$POLICY" "$WEB_POLICY" "$MATRIX" "$WORKFLOW"; do
  if [[ ! -f "$file" ]]; then
    fail "missing file: ${file#$ROOT/}"
  fi
done

required_policy_markers=(
  "timing_collection_version: 1"
  "source_marker: FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS"
  "log_query_pattern: FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS="
  "minimum_samples_before_relaxing: 5"
  "threshold_policy_status: advisory"
  "blocking_threshold_seconds: unset"
  "nightly_split_requires: measured_cost_review"
  "current_nightly_split: false"
  "review_dataset: ci_logs"
  "report_artifact_format: key_value_v1"
  "report_script: scripts/flutter_web_timing_report.sh"
  "report_test: scripts/flutter_web_timing_report_test.sh"
  "report_required_fields: report_version,source_marker,sample_count,duration_seconds_values,latest_duration_seconds,min_duration_seconds,max_duration_seconds,branch_type,runner_context"
  "ci_artifact_upload: true"
  "artifact_name: flutter-web-timing-report"
  "artifact_path: /tmp/manpasik_flutter_web_timing.env"
  "artifact_upload_action: actions/upload-artifact@v4"
  "report_generation_step: Flutter web timing report"
  "report_upload_step: Upload Flutter web timing report"
)

for marker in "${required_policy_markers[@]}"; do
  if ! grep -Fq "$marker" "$POLICY"; then
    fail "missing timing collection marker: $marker"
  fi
done

web_policy_markers=(
  "timing_capture: true"
  "timing_marker: FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS"
  "nightly_split: false"
)

for marker in "${web_policy_markers[@]}"; do
  if ! grep -Fq "$marker" "$WEB_POLICY"; then
    fail "web release policy missing marker: $marker"
  fi
done

workflow_markers=(
  "Flutter web timing collection policy"
  "scripts/ci_web_gate_timing_collection_policy_test.sh"
  "Flutter web timing report"
  "Upload Flutter web timing report"
  "actions/upload-artifact@v4"
  "flutter-web-timing-report"
  "/tmp/manpasik_flutter_web_timing.env"
)

for marker in "${workflow_markers[@]}"; do
  if ! grep -Fq "$marker" "$WORKFLOW"; then
    fail "workflow missing marker: $marker"
  fi
done

for path in scripts/flutter_web_timing_report.sh scripts/flutter_web_timing_report_test.sh; do
  if [[ ! -f "$ROOT/$path" ]]; then
    fail "missing report file: $path"
  fi
done

echo "CI_WEB_GATE_TIMING_COLLECTION_POLICY_PASS"
