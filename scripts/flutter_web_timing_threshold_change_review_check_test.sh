#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CHECKER="$ROOT/scripts/flutter_web_timing_threshold_change_review_check.sh"
AGGREGATOR="$ROOT/scripts/flutter_web_timing_sample_aggregate.sh"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

fail() {
  echo "FLUTTER_WEB_TIMING_THRESHOLD_CHANGE_REVIEW_CHECK_TEST_FAIL $1" >&2
  exit 1
}

write_sample() {
  local file="$1"
  local duration="$2"
  local branch_type="$3"
  cat >"$TMP_DIR/review/samples/$file" <<EOF_SAMPLE
report_version=1
source_marker=FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS
sample_count=1
duration_seconds_values=$duration
latest_duration_seconds=$duration
min_duration_seconds=$duration
max_duration_seconds=$duration
branch_type=$branch_type
runner_context=Linux-X64
EOF_SAMPLE
}

mkdir -p "$TMP_DIR/review/samples"
write_sample sample-01.env 12 pull_request
write_sample sample-02.env 14 pull_request
write_sample sample-03.env 15 pull_request
write_sample sample-04.env 20 release_branch
write_sample sample-05.env 30 release_branch

cat >"$TMP_DIR/review/manifest.env" <<'EOF_MANIFEST'
review_id=2026-05-14-pr-gate-cost-review
artifact_source=github_actions
artifact_name=flutter-web-timing-report
sample_count=5
sample_files=sample-01.env,sample-02.env,sample-03.env,sample-04.env,sample-05.env
aggregate_output=aggregate.env
review_status=proposed
no_phi_attestation=true
EOF_MANIFEST

bash "$AGGREGATOR" --input-dir "$TMP_DIR/review/samples" --output "$TMP_DIR/review/aggregate.env" >/dev/null

cat >"$TMP_DIR/review-audit.md" <<EOF_AUDIT
# Flutter Web Timing Review

review_template_version: 1
review_fixture_policy: docs/ci/flutter-web-timing-review-fixtures.md
aggregate_source: aggregate.env
review_id: 2026-05-14-pr-gate-cost-review
decision: proposed
sample_count: 5
median_seconds: 15
p95_seconds: 30
worst_case_seconds: 30
branch_types: pull_request,release_branch
runner_contexts: Linux-X64
validator_command: bash scripts/flutter_web_timing_review_fixture_validate.sh --review-dir $TMP_DIR/review
EOF_AUDIT

bash "$CHECKER" --review-dir "$TMP_DIR/review" --audit-file "$TMP_DIR/review-audit.md" | grep -Fq "FLUTTER_WEB_TIMING_THRESHOLD_CHANGE_REVIEW_CHECK_PASS" || fail "expected wrapper pass"

cp "$TMP_DIR/review-audit.md" "$TMP_DIR/bad-audit.md"
sed -i 's/p95_seconds: 30/p95_seconds: 29/' "$TMP_DIR/bad-audit.md"
if bash "$CHECKER" --review-dir "$TMP_DIR/review" --audit-file "$TMP_DIR/bad-audit.md" >/tmp/manpasik_threshold_review_bad.out 2>&1; then
  fail "expected p95 mismatch failure"
fi
grep -Fq "audit field mismatch: p95_seconds" /tmp/manpasik_threshold_review_bad.out || fail "wrong p95 mismatch failure"

echo "FLUTTER_WEB_TIMING_THRESHOLD_CHANGE_REVIEW_CHECK_TEST_PASS"
