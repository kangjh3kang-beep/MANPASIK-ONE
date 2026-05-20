#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SCRIPT="$ROOT/scripts/flutter_web_timing_sample_aggregate.sh"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

fail() {
  echo "FLUTTER_WEB_TIMING_SAMPLE_AGGREGATE_TEST_FAIL $1" >&2
  exit 1
}

write_sample() {
  local file="$1"
  local duration="$2"
  local branch_type="$3"
  cat >"$TMP_DIR/$file" <<EOF_SAMPLE
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

write_sample pr-1.env 12 pull_request
write_sample pr-2.env 14 pull_request
write_sample pr-3.env 15 pull_request
write_sample release-1.env 20 release_branch
write_sample release-2.env 30 release_branch

OUTPUT="$TMP_DIR/aggregate.env"
bash "$SCRIPT" --input-dir "$TMP_DIR" --output "$OUTPUT"

grep -Fq "aggregate_version=1" "$OUTPUT" || fail "missing aggregate_version"
grep -Fq "sample_count=5" "$OUTPUT" || fail "missing sample_count"
grep -Fq "branch_types=pull_request,release_branch" "$OUTPUT" || fail "missing branch types"
grep -Fq "median_seconds=15" "$OUTPUT" || fail "missing median"
grep -Fq "p95_seconds=30" "$OUTPUT" || fail "missing p95"
grep -Fq "worst_case_seconds=30" "$OUTPUT" || fail "missing worst case"

if bash "$SCRIPT" --input-dir "$TMP_DIR" --min-samples 6 --output "$TMP_DIR/too-few.env" >/tmp/manpasik_aggregate_too_few.out 2>&1; then
  fail "expected min sample failure"
fi
grep -Fq "minimum samples not met" /tmp/manpasik_aggregate_too_few.out || fail "wrong min sample failure"

echo "FLUTTER_WEB_TIMING_SAMPLE_AGGREGATE_TEST_PASS"
